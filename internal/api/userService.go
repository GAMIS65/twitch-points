package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gamis65/twitch-points/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nicklaw5/helix/v2"
)

type ChannelCustomRewardsParams struct {
	BroadcasterID                     string `query:"broadcaster_id"`
	Title                             string `json:"title"`
	Cost                              int    `json:"cost"`
	Prompt                            string `json:"prompt"`
	IsEnabled                         bool   `json:"is_enabled"`
	BackgroundColor                   string `json:"background_color,omitempty"`
	IsUserInputRequired               bool   `json:"is_user_input_required"`
	IsMaxPerStreamEnabled             bool   `json:"is_max_per_stream_enabled"`
	MaxPerStream                      int    `json:"max_per_stream"`
	IsMaxPerUserPerStreamEnabled      bool   `json:"is_max_per_user_per_stream_enabled"`
	MaxPerUserPerStream               int    `json:"max_per_user_per_stream"`
	IsGlobalCooldownEnabled           bool   `json:"is_global_cooldown_enabled"`
	GlobalCooldownSeconds             int    `json:"global_cooldown_seconds"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"should_redemptions_skip_request_queue"`
}

type TwitchAPIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type CreateCustomRewardResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (s *Server) getUserData(accessToken string) (*helix.User, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     s.oauthConfig.ClientID,
		ClientSecret: s.oauthConfig.ClientSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Twitch client: %w", err)
	}

	client.SetUserAccessToken(accessToken)

	resp, err := client.GetUsers(&helix.UsersParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("twitch API error: %s", resp.Error)
	}

	if len(resp.Data.Users) == 0 {
		return nil, fmt.Errorf("no user information found")
	}

	return &resp.Data.Users[0], nil
}

func (s *Server) addRewardHandler(w http.ResponseWriter, r *http.Request) {
	session, err := s.sessionStore.Get(r, "twitch-oauth-session")
	if err != nil {
		slog.Error("Error getting session", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userID, ok := session.Values["user_id"].(string)
	if !ok || userID == "" {
		slog.Error("User ID not found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	accessToken, ok := session.Values["access_token"].(string)
	if !ok || accessToken == "" {
		slog.Error("Access token not found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	existingReward, err := s.db.GetRewardsByStreamer(r.Context(), pgtype.Text{String: userID, Valid: true})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		slog.Error("Error getting streamers rewards from the database", "error", err, "id", userID)
		http.Error(w, "Error getting streamers rewards from the database", http.StatusInternalServerError)
		return
	}

	if len(existingReward) > 0 {
		err := s.db.DeleteRewardsByStreamerID(r.Context(), pgtype.Text{String: userID, Valid: true})
		if err != nil {
			slog.Error("Error while trying to delete existing rewards", "error", err, "id", userID)
			http.Error(w, "Error while trying to delete existing rewards", http.StatusInternalServerError)
			return
		}
		slog.Info("Existing rewards deleted for user", "id", userID)
	}

	rewardID, err := createCustomChannelPointReward(s.oauthConfig.ClientID, accessToken, &ChannelCustomRewardsParams{
		BroadcasterID:                userID,
		Title:                        "1 Giveaway Entry",
		Cost:                         100,
		IsEnabled:                    true,
		IsMaxPerUserPerStreamEnabled: true,
		MaxPerUserPerStream:          1,
	})

	if err != nil {
		slog.Error("Failed to create a channel point reward", "error", err, "id", userID)
		http.Error(w, "Failed to create a channel point reward", http.StatusInternalServerError)
		return
	}

	_, err = s.db.CreateReward(r.Context(), db.CreateRewardParams{
		RewardID:   rewardID,
		StreamerID: pgtype.Text{String: userID, Valid: true},
	})

	if err != nil {
		slog.Error("Error adding new reward to database", "error", err, "id", userID, "reward_id", rewardID)
		http.Error(w, "Error adding new reward to database", http.StatusInternalServerError)
		return
	}

	slog.Info("Channel point reward created successfully", "id", userID, "reward_id", rewardID)
	w.WriteHeader(http.StatusOK)
}

// TODO: Move this somewhere
func createCustomChannelPointReward(clientID string, accessToken string, reward *ChannelCustomRewardsParams) (string, error) {
	url := fmt.Sprintf("https://api.twitch.tv/helix/channel_points/custom_rewards?broadcaster_id=%s", reward.BroadcasterID)

	payloadBytes, err := json.Marshal(reward)
	if err != nil {
		return "", fmt.Errorf("error marshaling JSON payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("error creating HTTP request: %w", err)
	}

	req.Header.Set("Client-Id", clientID)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending HTTP request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		var twitchError TwitchAPIError
		if err := json.Unmarshal(bodyBytes, &twitchError); err == nil {
			return "", fmt.Errorf("twitch api error: status %d, message: %s", twitchError.Status, twitchError.Message)
		}
		return "", fmt.Errorf("twitch api error: status %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the successful response to extract the reward ID
	var successResponse CreateCustomRewardResponse
	if err := json.Unmarshal(bodyBytes, &successResponse); err != nil {
		return "", fmt.Errorf("error decoding successful response: %w, body: %s", err, string(bodyBytes))
	}

	if len(successResponse.Data) > 0 {
		return successResponse.Data[0].ID, nil
	}

	return "", fmt.Errorf("successful response did not contain reward data: %s", string(bodyBytes))
}

package api

import (
	"errors"
	"fmt"
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
	client, err := helix.NewClient(&helix.Options{
		ClientID:     s.oauthConfig.ClientID,
		ClientSecret: s.oauthConfig.ClientSecret,
	})

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

	logger := s.logger.With(
		slog.String("user_id", userID),
	)

	accessToken, ok := session.Values["access_token"].(string)
	if !ok || accessToken == "" {
		logger.Error("Access token not found in session")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	existingReward, err := s.db.GetRewardsByStreamer(r.Context(), pgtype.Text{String: userID, Valid: true})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.Error("Error getting streamers rewards from the database", "error", err)
		http.Error(w, "Error getting streamers rewards from the database", http.StatusInternalServerError)
		return
	}

	if len(existingReward) > 0 {
		err := s.db.DeleteRewardsByStreamerID(r.Context(), pgtype.Text{String: userID, Valid: true})
		if err != nil {
			logger.Error("Error while trying to delete existing rewards", "error", err)
			http.Error(w, "Error while trying to delete existing rewards", http.StatusInternalServerError)
			return
		}
		logger.Info("Existing rewards deleted for user")
	}

	response, err := client.CreateCustomReward(&helix.ChannelCustomRewardsParams{
		BroadcasterID:                userID,
		Title:                        "1 Giveaway Entry",
		Cost:                         100,
		IsEnabled:                    true,
		IsMaxPerUserPerStreamEnabled: true,
		MaxPerUserPerStream:          1,
	})

	rewardID := response.Data.ChannelCustomRewards[0].ID

	if err != nil {
		logger.Error("Failed to create a channel point reward", "error", err)
		http.Error(w, "Failed to create a channel point reward", http.StatusInternalServerError)
		return
	}

	_, err = s.db.CreateReward(r.Context(), db.CreateRewardParams{
		RewardID:   rewardID,
		StreamerID: pgtype.Text{String: userID, Valid: true},
	})

	if err != nil {
		logger.Error("Error adding new reward to database", "error", err, "reward_id", rewardID)
		http.Error(w, "Error adding new reward to database", http.StatusInternalServerError)
		return
	}

	logger.Info("Channel point reward created successfully", "reward_id", rewardID)
	w.WriteHeader(http.StatusOK)
}

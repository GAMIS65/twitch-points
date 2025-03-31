package api

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gamis65/twitch-points/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nicklaw5/helix/v2"
)

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

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.sessionStore.Get(r, "twitch-oauth-session")

		accessToken, ok := session.Values["access_token"].(string)
		expiry, expiryOk := session.Values["expiry"].(int64)

		if !ok || !expiryOk || accessToken == "" {
			http.Redirect(w, r, s.frontendURL, http.StatusSeeOther)
			return
		}

		// Check if the token is valid with Twitch API
		tokenValidity, err := isTokenValid(accessToken)
		if err != nil {
			fmt.Println(err)
		}

		if time.Now().Unix() > expiry || !tokenValidity {
			_, err := s.refreshAccessToken(r, w)
			if err != nil {
				http.Redirect(w, r, s.frontendURL, http.StatusSeeOther)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func isTokenValid(token string) (bool, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		return false, fmt.Errorf("couldn't make a request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (s *Server) addRewardHandler(w http.ResponseWriter, r *http.Request) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     s.oauthConfig.ClientID,
		ClientSecret: s.oauthConfig.ClientSecret,
	})

	session, _ := s.sessionStore.Get(r, "twitch-oauth-session")
	user_id, ok := session.Values["user_id"].(string)
	access_token, accessTokenOk := session.Values["access_token"].(string)

	if !ok {
		slog.Error("Error getting user id from session when adding a reward")
		http.Error(w, "Error getting streamers rewards from the database", http.StatusUnauthorized)
		return
	}

	if !accessTokenOk {
		slog.Error("Error getting access token from session when adding a reward")
		http.Error(w, "Error getting streamers rewards from the database", http.StatusUnauthorized)
		return
	}

	client.SetUserAccessToken(access_token)

	existingReward, err := s.db.GetRewardsByStreamer(r.Context(), pgtype.Text{String: user_id, Valid: true})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting streamers rewards from the database", "error", err, "id", user_id)
			http.Error(w, "Error getting streamers rewards from the database", http.StatusInternalServerError)
			return
		}
	}

	if len(existingReward) > 0 {
		err := s.db.DeleteRewardsByStreamerID(r.Context(), pgtype.Text{String: user_id, Valid: true})
		if err != nil {
			slog.Error("Error while trying to delete exsting rewards", "error", err, "id", user_id)
			http.Error(w, "Error while trying to delete existing rewards", http.StatusInternalServerError)
		}
	}

	resp, err := client.CreateCustomReward(&helix.ChannelCustomRewardsParams{
		BroadcasterID: user_id,
		Title:         "giveaway test",
		Cost:          1,
	})

	if err != nil {
		slog.Error("Failed to create a channel point reward", "error", err, "id", user_id)
		http.Error(w, "Failed to create a channel point reward", http.StatusInternalServerError)
		return
	}

	if len(resp.Data.ChannelCustomRewards) == 0 {
		slog.Error("Twitch API returned success but no rewards were created", "id", user_id)
		http.Error(w, "Failed to create a channel point reward", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != 200 {
		slog.Error("Failed to create a channel point reward", "status", resp.StatusCode, "error", resp.Error, "id", user_id)
		http.Error(w, "Failed to create a channel point reward", http.StatusInternalServerError)
		return
	}

	_, err = s.db.CreateReward(r.Context(), db.CreateRewardParams{
		RewardID:   resp.Data.ChannelCustomRewards[0].ID,
		StreamerID: pgtype.Text{String: user_id, Valid: true},
	})

	if err != nil {
		slog.Error("Error adding new reward to database", "error", err, "id", user_id)
		http.Error(w, "Error adding new reward to database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

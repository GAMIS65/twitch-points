package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/gamis65/twitch-points/internal/db"
	"github.com/gamis65/twitch-points/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/oauth2"
)

type AuthResponse struct {
	UserID       string `json:"user_id"`
	Login        string `json:"login"`
	DisplayName  string `json:"display_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Server) beginAuthHandler(w http.ResponseWriter, r *http.Request) {
	state := util.GenerateRandomState()

	session, _ := s.sessionStore.Get(r, "twitch-oauth-session")
	session.Values["state"] = state
	session.Save(r, w)

	url := s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) callbackHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessionStore.Get(r, "twitch-oauth-session")

	sessionStateValue := session.Values["state"]
	if sessionStateValue == nil {
		http.Error(w, "Missing state parameter in session", http.StatusBadRequest)
		return
	}

	sessionState, ok := sessionStateValue.(string)
	if !ok {
		http.Error(w, "Invalid state parameter type", http.StatusBadRequest)
		return
	}

	queryState := r.URL.Query().Get("state")
	if queryState != sessionState {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		slog.Error("Failed to exchange token", "error", err)
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	userData, err := s.getUserData(token.AccessToken)
	if err != nil {
		slog.Error("Error getting user data", "error", err)
		http.Error(w, "Error getting user data", http.StatusInternalServerError)
		return
	}

	if !util.IsDev() {
		// Sign in should be only available to channels with channel points
		if userData.BroadcasterType == "" {
			slog.Info("A user who is not an affiliate or a partner tried to sign in", "username", userData.Login)
			http.Error(w, "You must be an affiliate to log in", http.StatusUnauthorized)
			return
		}
	}

	session.Values["access_token"] = token.AccessToken
	session.Values["refresh_token"] = token.RefreshToken
	session.Values["expiry"] = token.Expiry.Unix()
	session.Values["user_id"] = userData.ID
	session.Options.MaxAge = int(token.Expiry.Unix() - time.Now().Unix())

	existingUser, err := s.db.GetStreamerByID(r.Context(), userData.ID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Database error when fetching user", "error", err, "id", userData.ID)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
	}

	if existingUser == (db.Streamer{}) {
		newUser, err := s.db.CreateStreamer(r.Context(), db.CreateStreamerParams{
			TwitchID:        userData.ID,
			Username:        userData.Login,
			AccessToken:     pgtype.Text{String: token.AccessToken, Valid: true},
			RefreshToken:    pgtype.Text{String: token.RefreshToken, Valid: true},
			ProfileImageUrl: pgtype.Text{String: userData.ProfileImageURL, Valid: true},
			Verified:        pgtype.Bool{Bool: false, Valid: true},
		})

		if err != nil {
			slog.Error("Error creating a new user", "error", err, "id", userData.ID, "username", userData.Login)
			http.Redirect(w, r, s.frontendURL+"/auth/twitch/login", http.StatusTemporaryRedirect)
			return
		}

		s.twitchEventSub.SubscribeToEvents([]db.Streamer{newUser})

		slog.Info("Created a new user", "id", userData.ID, "username", userData.Login)
	} else {
		_, err := s.db.UpdateStreamerTokens(r.Context(), db.UpdateStreamerTokensParams{
			TwitchID:     userData.ID,
			AccessToken:  pgtype.Text{String: token.AccessToken, Valid: true},
			RefreshToken: pgtype.Text{String: token.RefreshToken, Valid: true},
		})

		if err != nil {
			slog.Error("Error updating user tokens", "error", err, "id", userData.ID, "username", userData.Login)
		}

		slog.Info("User logged in", "id", userData.ID, "username", userData.Login)
	}

	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, s.frontendURL+"/addreward", http.StatusTemporaryRedirect)
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := s.sessionStore.Get(r, "twitch-oauth-session")

	session.Values["access_token"] = ""
	session.Values["refresh_token"] = ""
	session.Values["user_id"] = ""
	session.Values["state"] = ""
	session.Values["expiry"] = ""
	session.Options.MaxAge = -1

	session.Save(r, w)

	http.Redirect(w, r, s.frontendURL, http.StatusTemporaryRedirect)
}

func (s *Server) refreshAccessToken(r *http.Request, w http.ResponseWriter) (*oauth2.Token, error) {
	session, _ := s.sessionStore.Get(r, "twitch-oauth-session")

	refreshToken, ok := session.Values["refresh_token"].(string)
	if !ok || refreshToken == "" {
		return nil, fmt.Errorf("no refresh token available")
	}

	tokenSource := s.oauthConfig.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: refreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	session.Values["access_token"] = newToken.AccessToken
	session.Values["refresh_token"] = newToken.RefreshToken
	session.Values["expiry"] = newToken.Expiry.Unix()

	err = session.Save(r, w)
	if err != nil {
		fmt.Printf("Error saving session: %v", err)
	}

	return newToken, nil
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.sessionStore.Get(r, "twitch-oauth-session")

		accessToken, ok := session.Values["access_token"].(string)
		expiry, expiryOk := session.Values["expiry"].(int64)
		userId, userIdOk := session.Values["user_id"].(string)

		if !ok || !expiryOk || !userIdOk || accessToken == "" {
			slog.Error("Failed to authenticate user", "accessTokenOk", ok, "expiry", expiry, "expiryOk", expiryOk, "userId", userId, "userIdOk", userIdOk)
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

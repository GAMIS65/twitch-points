package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gamis65/twitch-points/internal/util"
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

	if util.IsDev() {
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

	session.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, s.frontendURL, http.StatusTemporaryRedirect)
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

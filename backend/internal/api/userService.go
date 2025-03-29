package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

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

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
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

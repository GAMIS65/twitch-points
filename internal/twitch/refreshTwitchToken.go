package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TwitchTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func GetRefreshTwitchToken(refreshToken, clientID, clientSecret string) (TwitchTokenResponse, error) {
	endpoint := "https://id.twitch.tv/oauth2/token"

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return TwitchTokenResponse{}, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TwitchTokenResponse{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TwitchTokenResponse{}, fmt.Errorf("error reading response body: %w", err)
	}

	var tokenResponse TwitchTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return TwitchTokenResponse{}, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return tokenResponse, nil
}

package twitch

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

func GetRefreshTwitchToken(clientID, clientSecret, refreshToken string) (*oauth2.Token, error) {
	twitchOAuth2Endpoint := oauth2.Endpoint{
		AuthURL:  "https://id.twitch.tv/oauth2/authorize",
		TokenURL: "https://id.twitch.tv/oauth2/token",
	}

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     twitchOAuth2Endpoint,
	}

	tokenSource := conf.TokenSource(context.Background(), &oauth2.Token{RefreshToken: refreshToken})

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("unable to refresh token: %v", err)
	}

	return newToken, nil
}

package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gamis65/twitch-points/internal/api"
	"github.com/gamis65/twitch-points/internal/db"
	"github.com/gamis65/twitch-points/internal/twitch"
	"github.com/gamis65/twitch-points/internal/util"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	eventSubs "github.com/joeyak/go-twitch-eventsub"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	twitchOauth "golang.org/x/oauth2/twitch"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
		os.Exit(1)
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	host := os.Getenv("HOST")
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	sessionKey := os.Getenv("SESSION_KEY")
	backendDomainName := os.Getenv("BACKEND_DOMAIN_NAME")
	dbURL := os.Getenv("DB_URL")
	myAccessToken := os.Getenv("MY_ACCESS_TOKEN")
	myChannelID := os.Getenv("MY_CHANNEL_ID")

	oauthConfig := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  backendDomainName + "/auth/twitch/callback",
		Scopes:       []string{"channel:read:redemptions", "channel:manage:redemptions"},
		Endpoint:     twitchOauth.Endpoint,
	}

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		slog.Error("Error connecting to the database", "error", err)
		os.Exit(1)
	}
	defer conn.Close()

	dbStore := db.NewStore(conn)
	if dbStore == nil {
		slog.Info("db url", "url", dbURL)
		log.Fatal("dbStore is nil")
	}

	events := []eventSubs.EventSubscription{
		eventSubs.SubStreamOnline,
		eventSubs.SubStreamOffline,
		eventSubs.SubChannelChannelPointsCustomRewardUpdate,
		eventSubs.SubChannelChannelPointsCustomRewardRedemptionAdd,
		eventSubs.SubChannelUpdate,
	}

	twitchEventSubClient := twitch.NewTwitchClient(clientID, clientSecret, dbStore, events, myChannelID, myAccessToken)
	twitchEventSubClient.Initialize()

	go func() {
		err := twitchEventSubClient.Connect()
		if err != nil {
			slog.Error("Failed to connect to twitch", "error", err)
		}
	}()

	sessionStore := sessions.NewCookieStore([]byte(sessionKey))
	if util.IsDev() {
		sessionStore.Options = &sessions.Options{
			Path:     "/",
			Domain:   "localhost",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}
	} else {
		sessionStore.Options = &sessions.Options{
			Path: "/",
			// TODO: change this
			Domain:   "." + frontendURL,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
	}

	server := api.NewServer(
		host,
		frontendURL,
		backendDomainName,
		oauthConfig,
		sessionStore,
		dbStore,
		twitchEventSubClient,
	)

	slog.Info("Server listening", "host", host)
	log.Fatal(server.Start())
}

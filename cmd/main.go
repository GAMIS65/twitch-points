package main

import (
	"context"
	"fmt"
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
	sessionKey := os.Getenv("SESSION_KEY")
	backendDomainName := os.Getenv("BACKEND_DOMAIN_NAME")
	cookieDomain := os.Getenv("COOKIE_DOMAIN")

	// twitch
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	twitchWebhookSecret := os.Getenv("TWITCH_WEBHOOK_SECRET")
	twitchWebhookURL := os.Getenv("TWITCH_WEBHOOK_URL")

	// DB
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	level := slog.LevelInfo
	if util.IsDev() {
		level = slog.LevelDebug
	}

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})

	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

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

	// Define the events to subscribe to
	events := []string{
		"stream.online",
		"stream.offline",
		"channel.update",
		"channel.channel_points_custom_reward_redemption.add",
		"channel.channel_points_custom_reward.update",
	}

	// Initialize the Twitch webhook client
	twitchWebhookClient, err := twitch.NewTwitchClient(
		clientID,
		clientSecret,
		twitchWebhookSecret,
		twitchWebhookURL,
		dbStore,
		events,
	)

	if err != nil {
		slog.Error("Failed to create Twitch webhook client", "error", err)
		os.Exit(1)
	}

	// Initialize the client and set up event handlers
	go func() {
		twitchWebhookClient.Initialize()
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
			Path:     "/",
			Domain:   cookieDomain,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
		}
	}

	server := api.NewServer(&api.ServerConfig{
		Host:          host,
		FrontendURL:   frontendURL,
		OAuthConfig:   oauthConfig,
		SessionStore:  sessionStore,
		DBStore:       dbStore,
		TwitchWebhook: twitchWebhookClient,
	})

	slog.Info("Server listening", "host", host)
	log.Fatal(server.Start())
}

package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gamis65/twitch-points/internal/api"
	"github.com/gamis65/twitch-points/internal/db"
	"github.com/gamis65/twitch-points/internal/util"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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
		clientID,
		clientSecret,
		sessionStore,
		dbStore,
	)

	slog.Info("Server listening", "host", host)
	log.Fatal(server.Start())
}

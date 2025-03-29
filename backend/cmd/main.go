package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gamis65/twitch-points/internal/api"
	"github.com/gamis65/twitch-points/internal/util"
	"github.com/gorilla/sessions"
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

	store := sessions.NewCookieStore([]byte(sessionKey))
	if util.IsDev() {
		store.Options = &sessions.Options{
			Path:     "/",
			Domain:   "localhost",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		}
	} else {
		store.Options = &sessions.Options{
			Path:     "/",
			Domain:   "." + frontendURL,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		}
	}

	server := api.NewServer(
		host,
		frontendURL,
		backendDomainName,
		clientID,
		clientSecret,
		store,
	)

	slog.Info("Server listening", "host", host)
	log.Fatal(server.Start())
}

package api

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

type Server struct {
	host         string
	frontendURL  string
	sessionStore *sessions.CookieStore
	oauthConfig  *oauth2.Config
}

func NewServer(host, frontendURL, backendDomainName, clientID, clientSecret string, sessionStore *sessions.CookieStore) *Server {
	return &Server{
		host:         host,
		frontendURL:  frontendURL,
		sessionStore: sessionStore,
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  backendDomainName + "/auth/twitch/callback",
			Scopes:       []string{"channel:read:redemptions", "channel:manage:redemptions"},
			Endpoint:     twitch.Endpoint,
		},
	}
}

func (s *Server) SetupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{s.frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	r.Get("/auth/twitch", s.beginAuthHandler)
	r.Get("/auth/twitch/callback", s.callbackHandler)
	r.Get("/logout/twitch", s.logoutHandler)

	return r
}

func (s *Server) Start() error {
	slog.Info("Starting server", "address", s.host)
	return http.ListenAndServe(s.host, s.SetupRoutes())
}

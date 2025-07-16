package api

import (
	"log/slog"
	"net/http"

	"github.com/gamis65/twitch-points/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	eventSub "github.com/gamis65/twitch-points/internal/twitch"
)

type Server struct {
	host          string
	frontendURL   string
	sessionStore  *sessions.CookieStore
	oauthConfig   *oauth2.Config
	db            *db.DBStore
	twitchWebhook *eventSub.TwitchWebhookClient
	logger        *slog.Logger
}

type ServerConfig struct {
	Host          string
	FrontendURL   string
	OAuthConfig   *oauth2.Config
	SessionStore  *sessions.CookieStore
	DBStore       *db.DBStore
	TwitchWebhook *eventSub.TwitchWebhookClient
	Logger        *slog.Logger
}

func NewServer(cfg *ServerConfig) *Server {
	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	return &Server{
		host:          cfg.Host,
		frontendURL:   cfg.FrontendURL,
		sessionStore:  cfg.SessionStore,
		oauthConfig:   cfg.OAuthConfig,
		db:            cfg.DBStore,
		twitchWebhook: cfg.TwitchWebhook,
		logger:        logger,
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
	r.With(s.authMiddleware).Post("/add-reward", s.addRewardHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(s.authMiddleware)
		r.Post("/add-reward", s.addRewardHandler)
		r.Get("/me", s.meHandler)
	})

	r.HandleFunc("/eventsub", s.twitchWebhook.GetHandler())

	r.Route("/giveaway", func(r chi.Router) {
		r.Get("/streamers", s.GetStreamersHandler)
		r.Get("/recent-entries", s.GetRecentEntriesHandler)
		r.Get("/participants-count", s.GetTotalParticipantsHandler)
		r.Get("/entries-count", s.GetTotalEntriesHandler)
		r.Get("/leaderboard", s.GetLeaderboardHandler)
	})
	return r
}

func (s *Server) Start() error {
	slog.Info("Starting server", "address", s.host)
	return http.ListenAndServe(s.host, s.SetupRoutes())
}

package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gamis65/twitch-points/internal/db"
	"github.com/gamis65/twitch-points/internal/util"
	"github.com/jackc/pgx/v5"
)

type TotalParticipantsResponse struct {
	TotalParticipants int `json:"total_participants"`
}

type TotalEntriesResponse struct {
	TotalEntries int `json:"total_entries"`
}

func (s *Server) GetStreamersHandler(w http.ResponseWriter, r *http.Request) {
	streamers, err := s.db.GetAllStreamers(r.Context())
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting streamers", "error", err)
			http.Error(w, "Error getting streamers", http.StatusInternalServerError)
			return
		}
	}

	if streamers == nil {
		emptyStreamers := make([]db.Streamer, 0)
		util.SendJSON(w, emptyStreamers)
		return
	}

	util.SendJSON(w, streamers)
}

func (s *Server) GetRecentEntriesHandler(w http.ResponseWriter, r *http.Request) {
	recentEntries, err := s.db.GetRecentRedemptionsWithUsernames(r.Context(), 10)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting recent redemptions", "error", err)
			http.Error(w, "Error getting recent redemptions", http.StatusInternalServerError)
			return
		}
	}

	if recentEntries == nil {
		emptyEntries := make([]db.Streamer, 0)
		util.SendJSON(w, emptyEntries)
		return
	}

	util.SendJSON(w, recentEntries)
}

func (s *Server) GetTotalParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	totalParticipants, err := s.db.GetTotalParticipantsCount(r.Context())
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting total participants count", "error", err)
			http.Error(w, "Error getting total participants count", http.StatusInternalServerError)
			return
		}
	}

	response := TotalParticipantsResponse{
		TotalParticipants: int(totalParticipants),
	}

	util.SendJSON(w, response)
}

func (s *Server) GetTotalEntriesHandler(w http.ResponseWriter, r *http.Request) {
	totalEntries, err := s.db.GetTotalRedemptionsCount(r.Context())
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting total entries count", "error", err)
			http.Error(w, "Error getting total entries count", http.StatusInternalServerError)
			return
		}
	}

	response := TotalEntriesResponse{
		TotalEntries: int(totalEntries),
	}

	util.SendJSON(w, response)
}

func (s *Server) GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	leaderboard, err := s.db.GetViewerLeaderboard(r.Context())
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting leaderboard data", "error", err)
			http.Error(w, "Error getting leaderboard data", http.StatusInternalServerError)
			return
		}
	}

	if leaderboard == nil {
		emptyLeaderboard := make([]db.Streamer, 0)
		util.SendJSON(w, emptyLeaderboard)
		return
	}

	util.SendJSON(w, leaderboard)
}

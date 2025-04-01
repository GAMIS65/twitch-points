package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gamis65/twitch-points/internal/db"
	"github.com/gamis65/twitch-points/internal/util"
	"github.com/jackc/pgx/v5"
)

func (s *Server) GetStreamers(w http.ResponseWriter, r *http.Request) {
	streamers, err := s.db.GetAllStreamers(r.Context())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			emptyStreamers := make([]db.Streamer, 0)
			util.SendJSON(w, emptyStreamers)
			return
		} else {
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

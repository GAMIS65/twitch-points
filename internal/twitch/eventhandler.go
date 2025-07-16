package twitch

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/LinneB/twitchwh"
	"github.com/gamis65/twitch-points/internal/db"
	"github.com/gamis65/twitch-points/internal/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type TwitchWebhookClient struct {
	client        *twitchwh.Client
	clientId      string
	clientSecret  string
	webhookSecret string
	webhookURL    string
	db            *db.DBStore
	events        []string
	logger        *slog.Logger
}

type StreamEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
}

type ChannelUpdateEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	Title                string `json:"title"`
}

type RewardUpdateEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	Title                string `json:"title"`
	Cost                 int    `json:"cost"`
}

type RewardRedemptionEvent struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	ID                   string `json:"id"` // Redemption ID
	Reward               struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Cost  int    `json:"cost"`
	} `json:"reward"`
	User struct {
		ID        string `json:"id"`
		Login     string `json:"login"`
		UserLogin string `json:"user_login"`
	} `json:"user"`
}

func NewTwitchClient(clientId string, clientSecret string, webhookSecret string, webhookURL string, dbStore *db.DBStore, eventsToSubscribeTo []string) (*TwitchWebhookClient, error) {
	client, err := twitchwh.New(twitchwh.ClientConfig{
		ClientID:      clientId,
		ClientSecret:  clientSecret,
		WebhookSecret: webhookSecret,
		WebhookURL:    webhookURL,
		Debug:         false,
	})

	if err != nil {
		return nil, err
	}

	return &TwitchWebhookClient{
		client:        client,
		clientId:      clientId,
		clientSecret:  clientSecret,
		webhookSecret: webhookSecret,
		webhookURL:    webhookURL,
		db:            dbStore,
		events:        eventsToSubscribeTo,
		logger:        slog.Default(),
	}, nil
}

// getEventLogger creates a contextualized logger with event specific fields
func (tc *TwitchWebhookClient) getEventLogger(eventType string, eventData any) *slog.Logger {
	logger := tc.logger.With(slog.String("eventType", eventType))

	// Add common fields based on event data structure
	switch data := eventData.(type) {
	case StreamEvent:
		logger = logger.With(
			slog.String("streamer_id", data.BroadcasterUserID),
			slog.String("streamer_username", data.BroadcasterUserLogin),
		)
	case RewardRedemptionEvent:
		logger = logger.With(
			slog.String("streamer_id", data.BroadcasterUserID),
			slog.String("streamer_username", data.BroadcasterUserLogin),
			slog.String("viewer_id", data.UserID),
			slog.String("viewer_username", data.UserLogin),
			slog.String("reward_id", data.Reward.ID),
			slog.String("reward_title", data.Reward.Title),
		)
	}

	return logger
}

func (tc *TwitchWebhookClient) Initialize() {
	// Stream live status
	tc.client.On("stream.online", tc.handleStreamOnline)
	tc.client.On("stream.offline", tc.handleStreamOffline)

	// Channel points
	tc.client.On("channel.channel_points_custom_reward_redemption.add", tc.handleRewardRedemption)
	tc.client.On("channel.update", tc.handleChannelUpdate)
	tc.client.On("channel.channel_points_custom_reward.update", tc.handleRewardUpdate)

	// Subs and cheers
	// TODO: Add a giveaway config
	// tc.client.On("channel.subscribe", tc.handleSubscription)
	// tc.client.On("channel.subscription.gift")
	// tc.client.On("channel.cheer")

	streamers, err := tc.db.GetAllStreamersWithTokens(context.Background())
	if err != nil {
		tc.logger.Error("Error getting streamers from the database", "error", err)
		return
	}

	// Refresh tokens for all streamers
	for i, streamer := range streamers {
		streamerLogger := tc.logger.With(
			slog.String("streamer_id", streamer.TwitchID),
			slog.String("streamer_username", streamer.Username),
		)

		newToken, err := GetRefreshTwitchToken(streamer.RefreshToken.String, tc.clientId, tc.clientSecret)
		if err != nil {
			streamerLogger.Error("Error refreshing token", "error", err)
			continue
		}

		_, err = tc.db.UpdateStreamerTokens(context.Background(), db.UpdateStreamerTokensParams{
			TwitchID:     streamer.TwitchID,
			AccessToken:  pgtype.Text{String: newToken.AccessToken, Valid: true},
			RefreshToken: pgtype.Text{String: newToken.RefreshToken, Valid: true},
		})

		if err != nil {
			streamerLogger.Error("Failed to refresh token", "error", err)
			continue
		}

		streamers[i].AccessToken.String = newToken.AccessToken
		streamers[i].RefreshToken.String = newToken.RefreshToken

		streamerLogger.Info("Refreshed streamer token")
	}

	tc.SubscribeToEvents(streamers)
}

func (tc *TwitchWebhookClient) SubscribeToEvents(streamers []db.Streamer) {
	for _, streamer := range streamers {
		streamerLogger := tc.logger.With(
			slog.String("streamer_id", streamer.TwitchID),
			slog.String("streamer_username", streamer.Username),
		)

		for _, event := range tc.events {
			streamerLogger.Info("Subscribing to an event", "event", event)

			err := tc.client.AddSubscription(event, "1", twitchwh.Condition{
				BroadcasterUserID: streamer.TwitchID,
			})

			if err != nil {
				streamerLogger.Error("Error subscribing to an event", "event", event, "error", err)
				continue
			}
		}
	}
}

func (tc *TwitchWebhookClient) handleStreamOnline(event json.RawMessage) {
	var eventData StreamEvent

	if err := json.Unmarshal(event, &eventData); err != nil {
		tc.logger.Error("Error parsing stream online event", "error", err)
		return
	}

	logger := tc.getEventLogger("stream.online", eventData)
	logger.Info("Streamer went live")
	util.SendWebHook(eventData.BroadcasterUserLogin + " went live")
}

func (tc *TwitchWebhookClient) handleStreamOffline(event json.RawMessage) {
	var eventData StreamEvent

	if err := json.Unmarshal(event, &eventData); err != nil {
		tc.logger.Error("Error parsing stream offline event", "error", err)
		return
	}

	logger := tc.getEventLogger("stream.offline", eventData)
	logger.Info("Streamer went offline")
}

func (tc *TwitchWebhookClient) handleRewardRedemption(event json.RawMessage) {
	var eventData RewardRedemptionEvent

	if err := json.Unmarshal(event, &eventData); err != nil {
		tc.logger.Error("Error parsing reward redemption event", "error", err)
		return
	}

	logger := tc.getEventLogger("reward.redemption", eventData)

	reward, err := tc.db.GetRewardsByStreamer(context.Background(), pgtype.Text{String: eventData.BroadcasterUserID, Valid: true})
	if err != nil {
		logger.Error("Error getting a reward for a streamer from db", "error", err)
		util.SendWebHook("Error getting a reward for a streamer from db " + eventData.BroadcasterUserLogin)
		return
	}

	if len(reward) == 0 {
		logger.Warn("No rewards found for streamer")
		return
	}

	// Check if the reward has the right ID
	if eventData.Reward.ID != reward[0].RewardID {
		return
	}

	viewer, err := tc.db.GetViewerByID(context.Background(), eventData.UserID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			logger.Error("Error getting viewer by id", "error", err)
			return
		}
	}

	emptyViewer := db.Viewer{}
	if viewer == emptyViewer {
		_, err = tc.db.CreateViewer(context.Background(), db.CreateViewerParams{
			TwitchID:     eventData.UserID,
			Username:     eventData.UserLogin,
			RegisteredIn: pgtype.Text{String: eventData.BroadcasterUserID, Valid: true},
		})

		if err != nil {
			logger.Error("Error creating a new viewer", "error", err)
			return
		}
	}

	_, err = tc.db.CreateRedemption(context.Background(), db.CreateRedemptionParams{
		MessageID:  eventData.ID,
		ViewerID:   pgtype.Text{String: eventData.UserID, Valid: true},
		StreamerID: pgtype.Text{String: eventData.BroadcasterUserID, Valid: true},
	})

	if err != nil {
		logger.Error("Error adding a redemption to db", "error", err)
		util.SendWebHook("Error adding a redemption to db " + eventData.BroadcasterUserLogin)
		return
	}

	logger.Info("User redeemed a reward")
	util.SendWebHook(eventData.UserLogin + " redeemed an entry in " + eventData.BroadcasterUserLogin)
}

func (tc *TwitchWebhookClient) handleChannelUpdate(event json.RawMessage) {
	var eventData ChannelUpdateEvent

	if err := json.Unmarshal(event, &eventData); err != nil {
		tc.logger.Error("Error parsing channel update event", "error", err)
		return
	}

	logger := tc.getEventLogger("channel.update", eventData)
	logger.Info("Channel updated", "title", eventData.Title)
}

func (tc *TwitchWebhookClient) handleRewardUpdate(event json.RawMessage) {
	var eventData RewardUpdateEvent

	if err := json.Unmarshal(event, &eventData); err != nil {
		tc.logger.Error("Error parsing reward update event", "error", err)
		return
	}

	logger := tc.getEventLogger("reward.update", eventData)
	logger.Warn("Streamer updated a channel point reward", "reward", eventData.Title, "cost", eventData.Cost)
}

// GetHandler returns the HTTP handler for webhook events
func (tc *TwitchWebhookClient) GetHandler() func(http.ResponseWriter, *http.Request) {
	return tc.client.Handler
}

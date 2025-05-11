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
	client         *twitchwh.Client
	clientId       string
	clientSecret   string
	webhookSecret  string
	webhookURL     string
	db             *db.DBStore
	events         []string
	myChannelID    string
	myAccessToken  string
	eventEndpoints map[string]string
}

func NewTwitchClient(clientId string, clientSecret string, webhookSecret string, webhookURL string, dbStore *db.DBStore, eventsToSubscribeTo []string, myChannelId string, myAccessToken string) (*TwitchWebhookClient, error) {
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
		myChannelID:   myChannelId,
		myAccessToken: myAccessToken,
		eventEndpoints: map[string]string{
			"stream.online":  "stream.online",
			"stream.offline": "stream.offline",
			"channel.update": "channel.update",
			"channel.channel_points_custom_reward_redemption.add": "channel.channel_points_custom_reward_redemption.add",
			"channel.channel_points_custom_reward.update":         "channel.channel_points_custom_reward.update",
		},
	}, nil
}

func (tc *TwitchWebhookClient) Initialize() {
	// Stream live status
	tc.client.On("stream.online", tc.handleStreamOnline)
	tc.client.On("stream.offline", tc.handleStreamOffline)

	// Channel points
	tc.client.On("channel.channel_points_custom_reward_redemption.add", tc.handleRewardRedemption)
	tc.client.On("channel.update", tc.handleChannelUpdate)
	tc.client.On("channel.channel_points_custom_reward.update", tc.handleRewardUpdate)

	streamers, err := tc.db.GetAllStreamersWithTokens(context.Background())
	if err != nil {
		slog.Error("Error getting streamers from the database", "error", err)
		return
	}

	// Refresh tokens for all streamers
	for i, streamer := range streamers {
		newToken, err := GetRefreshTwitchToken(streamer.RefreshToken.String, tc.clientId, tc.clientSecret)
		if err != nil {
			slog.Error("Error refreshing token", "error", err)
			continue
		}

		_, err = tc.db.UpdateStreamerTokens(context.Background(), db.UpdateStreamerTokensParams{
			TwitchID:     streamer.TwitchID,
			AccessToken:  pgtype.Text{String: newToken.AccessToken, Valid: true},
			RefreshToken: pgtype.Text{String: newToken.RefreshToken, Valid: true},
		})

		if err != nil {
			slog.Error("Failed to refresh token", "error", err, "id", streamer.TwitchID, "username", streamer.Username)
			continue
		}

		streamers[i].AccessToken.String = newToken.AccessToken
		streamers[i].RefreshToken.String = newToken.RefreshToken

		slog.Info("Refreshed streamer token", "id", streamer.TwitchID, "username", streamer.Username)
	}

	tc.SubscribeToEvents(streamers)
}

func (tc *TwitchWebhookClient) SubscribeToEvents(streamers []db.Streamer) {
	for _, streamer := range streamers {
		for _, event := range tc.events {
			slog.Info("Subscribing to an event", "streamerId", streamer.TwitchID, "streamerUsername", streamer.Username, "event", event)

			err := tc.client.AddSubscription(event, "1", twitchwh.Condition{
				BroadcasterUserID: streamer.TwitchID,
			})

			if err != nil {
				slog.Error("Error subscribing to an event", "event", event, "error", err, "streamerID", streamer.TwitchID)
				continue
			}

			slog.Info("Successfully subscribed to event", "event", event, "streamerID", streamer.TwitchID)
		}
	}
}

func (tc *TwitchWebhookClient) handleStreamOnline(event json.RawMessage) {
	var eventData struct {
		BroadcasterUserID    string `json:"broadcaster_user_id"`
		BroadcasterUserName  string `json:"broadcaster_user_name"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
	}

	if err := json.Unmarshal(event, &eventData); err != nil {
		slog.Error("Error parsing stream online event", "error", err)
		return
	}

	slog.Info("Streamer went live", "userId", eventData.BroadcasterUserID, "username", eventData.BroadcasterUserLogin)
	util.SendWebHook(eventData.BroadcasterUserLogin + " went live")
}

func (tc *TwitchWebhookClient) handleStreamOffline(event json.RawMessage) {
	var eventData struct {
		BroadcasterUserID    string `json:"broadcaster_user_id"`
		BroadcasterUserName  string `json:"broadcaster_user_name"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
	}

	if err := json.Unmarshal(event, &eventData); err != nil {
		slog.Error("Error parsing stream offline event", "error", err)
		return
	}

	slog.Info("Streamer went offline", "userId", eventData.BroadcasterUserID, "username", eventData.BroadcasterUserLogin)
}

func (tc *TwitchWebhookClient) handleRewardRedemption(event json.RawMessage) {
	var eventData struct {
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

	if err := json.Unmarshal(event, &eventData); err != nil {
		slog.Error("Error parsing reward redemption event", "error", err)
		return
	}

	reward, err := tc.db.GetRewardsByStreamer(context.Background(), pgtype.Text{String: eventData.BroadcasterUserID, Valid: true})
	if err != nil {
		slog.Error("Error getting a reward for a streamer from db", "error", err, "streamerId", eventData.BroadcasterUserID, "streamerUsername", eventData.BroadcasterUserName, "reward", eventData.Reward.Title, "viewerId", eventData.UserID, "viewerUsername", eventData.UserLogin)
		util.SendWebHook("Error getting a reward for a streamer from db " + eventData.BroadcasterUserLogin)
		return
	}

	if len(reward) == 0 {
		slog.Warn("No rewards found for streamer", "streamerId", eventData.BroadcasterUserID)
		return
	}

	// Check if the reward has the right ID
	if eventData.Reward.ID != reward[0].RewardID {
		return
	}

	viewer, err := tc.db.GetViewerByID(context.Background(), eventData.UserID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Error getting viewer by id", "error", err, "viewerId", eventData.UserID, "viewerUsername", eventData.UserLogin, "streamerId", eventData.BroadcasterUserID, "streamerUsername", eventData.BroadcasterUserName)
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
			slog.Error("Error creating a new viewer", "error", err, "viewerId", eventData.UserID, "viewerUsername", eventData.UserLogin, "streamerId", eventData.BroadcasterUserID, "streamerUsername", eventData.BroadcasterUserName)
			return
		}
	}

	_, err = tc.db.CreateRedemption(context.Background(), db.CreateRedemptionParams{
		MessageID:  eventData.ID,
		ViewerID:   pgtype.Text{String: eventData.UserID, Valid: true},
		StreamerID: pgtype.Text{String: eventData.BroadcasterUserID, Valid: true},
	})

	if err != nil {
		slog.Error("Error adding a redemption to db", "error", err, "streamerId", eventData.BroadcasterUserID, "viewerId", eventData.UserID)
		util.SendWebHook("Error adding a redemption to db " + eventData.BroadcasterUserLogin)
		return
	}

	slog.Info("User redeemed a reward", "userId", eventData.UserID, "username", eventData.UserLogin, "channel", eventData.BroadcasterUserLogin)
	util.SendWebHook(eventData.UserLogin + " redeemed an entry in " + eventData.BroadcasterUserLogin)
}

func (tc *TwitchWebhookClient) handleChannelUpdate(event json.RawMessage) {
	var eventData struct {
		BroadcasterUserID    string `json:"broadcaster_user_id"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
		Title                string `json:"title"`
	}

	if err := json.Unmarshal(event, &eventData); err != nil {
		slog.Error("Error parsing channel update event", "error", err)
		return
	}

	slog.Info("Channel updated", "channel", eventData.BroadcasterUserLogin, "title", eventData.Title)
}

func (tc *TwitchWebhookClient) handleRewardUpdate(event json.RawMessage) {
	var eventData struct {
		BroadcasterUserID    string `json:"broadcaster_user_id"`
		BroadcasterUserLogin string `json:"broadcaster_user_login"`
		Title                string `json:"title"`
		Cost                 int    `json:"cost"`
	}

	if err := json.Unmarshal(event, &eventData); err != nil {
		slog.Error("Error parsing reward update event", "error", err)
		return
	}

	slog.Warn("Streamer updated a channel point reward", "streamer", eventData.BroadcasterUserLogin, "reward", eventData.Title, "cost", eventData.Cost)
	util.SendWebHook(eventData.BroadcasterUserLogin + " updated a channel point reward, title=" + eventData.Title)
}

// GetHandler returns the HTTP handler for webhook events
func (tc *TwitchWebhookClient) GetHandler() func(http.ResponseWriter, *http.Request) {
	return tc.client.Handler
}

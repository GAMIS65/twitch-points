package twitch

import (
	"context"
	"log/slog"

	"github.com/gamis65/twitch-points/internal/db"
	"github.com/gamis65/twitch-points/internal/util"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joeyak/go-twitch-eventsub"
)

type TwitchEventSubClient struct {
	client        *twitch.Client
	clientId      string
	clientSecret  string
	db            *db.DBStore
	sessionID     string
	events        []twitch.EventSubscription
	myChannelID   string
	myAccessToken string
}

func NewTwitchClient(clientId string, clientSecret string, dbStore *db.DBStore, eventsToSubscribeTo []twitch.EventSubscription, myChannelId string, myAcessToken string) *TwitchEventSubClient {
	return &TwitchEventSubClient{
		client:        twitch.NewClient(),
		clientId:      clientId,
		clientSecret:  clientSecret,
		db:            dbStore,
		events:        eventsToSubscribeTo,
		myChannelID:   myChannelId,
		myAccessToken: myAcessToken,
	}
}

func (tc *TwitchEventSubClient) Initialize() {
	// Twitch websocket messages
	tc.client.OnWelcome(tc.handleWelcome)
	tc.client.OnRevoke(tc.handleRevoke)
	tc.client.OnReconnect(tc.handleReconnect)
	tc.client.OnError(func(err error) {
		slog.Error("Twitch error", "error", err)
	})

	// Stream live status
	tc.client.OnEventStreamOnline(tc.handleStreamOnline)
	tc.client.OnEventStreamOffline(tc.handleStreamOffline)

	// Channel points
	tc.client.OnEventChannelChannelPointsCustomRewardRedemptionAdd(tc.handleRewardRedemption)
	tc.client.OnEventChannelUpdate(tc.handleChannelUpdate)
	tc.client.OnEventChannelChannelPointsCustomRewardUpdate(tc.handleRewardUpdate)
}

func (tc *TwitchEventSubClient) handleWelcome(message twitch.WelcomeMessage) {
	slog.Info("Twitch welcome message", "status", message.Payload.Session.Status)

	tc.sessionID = message.Payload.Session.ID

	streamers, err := tc.db.GetAllStreamersWithTokens(context.Background())
	if err != nil {
		slog.Error("Error getting streamers from the database", "erro", err)
	}

	if len(streamers) < 1 {
		tc.subscribeToMyChannel()
	}

	for _, streamer := range streamers {
		newToken, err := GetRefreshTwitchToken(streamer.RefreshToken.String, tc.clientId, tc.clientSecret)
		if err != nil {
			slog.Error("Error refreshing token", "error", err)
		}

		_, err = tc.db.UpdateStreamerTokens(context.Background(), db.UpdateStreamerTokensParams{
			TwitchID:     streamer.TwitchID,
			AccessToken:  pgtype.Text{String: newToken.AccessToken, Valid: true},
			RefreshToken: pgtype.Text{String: newToken.RefreshToken, Valid: true},
		})

		if err != nil {
			slog.Error("Failed to refresh token", "error", err, "id", streamer.TwitchID, "username", streamer.Username)
		}
		slog.Info("Refreshed streamer token", "id", streamer.TwitchID, "username", streamer.Username)
	}

	tc.SubscribeToEvents(streamers)
}

func (tc *TwitchEventSubClient) subscribeToMyChannel() {
	for _, event := range tc.events {
		slog.Info("Subscribing to my channel event", "event", event)
		_, err := twitch.SubscribeEvent(twitch.SubscribeRequest{
			SessionID:   tc.sessionID,
			ClientID:    tc.clientId,
			AccessToken: tc.myAccessToken,
			Event:       event,
			Condition: map[string]string{
				"broadcaster_user_id": tc.myChannelID,
			},
		})

		if err != nil {
			slog.Error("Error subscribing to my channel event", "event", event, "error", err)
			return
		}
	}
}

func (tc *TwitchEventSubClient) SubscribeToEvents(streamers []db.Streamer) {
	for _, streamer := range streamers {
		for _, event := range tc.events {
			slog.Info("Subscribing to an event", "streamerId", streamer.TwitchID, "streamerUsername", streamer.Username, "event", event)
			_, err := twitch.SubscribeEvent(twitch.SubscribeRequest{
				SessionID:   tc.sessionID,
				ClientID:    tc.clientId,
				AccessToken: streamer.AccessToken.String,
				Event:       event,
				Condition: map[string]string{
					"broadcaster_user_id": streamer.TwitchID,
				},
			})

			if err != nil {
				slog.Error("Error subscribing to an event", "event", event, "error", err)
				return
			}
		}
	}
}

func (tc *TwitchEventSubClient) handleRevoke(message twitch.RevokeMessage) {
	slog.Warn("User revoked OAuth access", "userId", message.Payload.Subscription.Condition)
}

func (tc *TwitchEventSubClient) handleStreamOnline(event twitch.EventStreamOnline) {
	slog.Info("Streamer went live", "userId", event.BroadcasterUserId, "username", event.BroadcasterUserLogin)
	util.SendWebHook(event.BroadcasterUserLogin + " went live")
}

func (tc *TwitchEventSubClient) handleStreamOffline(event twitch.EventStreamOffline) {
	slog.Info("Streamer went offline", "userId", event.BroadcasterUserId, "username", event.BroadcasterUserLogin)
}

func (tc *TwitchEventSubClient) handleReconnect(message twitch.ReconnectMessage) {
	slog.Warn("Twitch WebSocket reconnected", "status", message.Payload.Session.Status)
	util.SendWebHook("Twitch WebSocket reconnected, status=" + message.Payload.Session.Status)
}

func (tc *TwitchEventSubClient) handleRewardRedemption(event twitch.EventChannelChannelPointsCustomRewardRedemptionAdd) {
	// TODO: Check for duplicates
	// TODO: Check for the reward ID in the database
	slog.Info("User redeemed a reward", "userId", event.UserID, "username", event.User.UserLogin, "channel", event.BroadcasterUserLogin)
	util.SendWebHook(event.UserLogin + " redeemed an entry in " + event.BroadcasterUserLogin)
}

func (tc *TwitchEventSubClient) handleChannelUpdate(event twitch.EventChannelUpdate) {
	slog.Info("Channel updated", "channel", event.BroadcasterUserLogin)
}

func (tc *TwitchEventSubClient) handleRewardUpdate(event twitch.EventChannelChannelPointsCustomRewardUpdate) {
	slog.Warn("Streamer updated a channel point reward", "streamer", event.BroadcasterUserLogin, "reward", event.Title, "cost", event.Cost)
	util.SendWebHook(event.BroadcasterUserLogin + " updated a channel point reward, title=" + event.Title)
}

func (tc *TwitchEventSubClient) Connect() error {
	return tc.client.Connect()
}

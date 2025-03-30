package twitch

import (
	"log/slog"

	"github.com/gamis65/twitch-points/internal/util"
	"github.com/joeyak/go-twitch-eventsub"
)

var (
	accessToken = "0"
)

type TwitchClient struct {
	client       *twitch.Client
	clientId     string
	clientSecret string
}

type TwitchStreamer struct {
	id           string
	accessToken  string
	refreshToken string
}

func NewTwitchClient(clientId string, clientSecret string) *TwitchClient {
	return &TwitchClient{
		client:       twitch.NewClient(),
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (tc *TwitchClient) Initialize() {
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

func (tc *TwitchClient) handleWelcome(message twitch.WelcomeMessage) {
	slog.Info("Twitch welcome message", "status", message.Payload.Session.Status)

	// TODO: Add new streamers when they log in
	// TODO: Fetch streamers from the database
	twitchStreamer1 := &TwitchStreamer{
		id:           "77738069",
		accessToken:  "0",
		refreshToken: "0",
	}

	twitchStreamer2 := &TwitchStreamer{
		id:           "77738068",
		accessToken:  "0",
		refreshToken: "0",
	}

	streamers := []TwitchStreamer{*twitchStreamer1, *twitchStreamer2}

	// TODO: Refresh access tokens
	// for _, streamer := range streamers {
	// 	newToken, err := GetRefreshTwitchToken(tc.clientId, tc.clientSecret, streamer.refreshToken)
	// 	if err != nil {
	// 		slog.Error("Error refreshing token", "error", err)
	// 	}
	// 	streamer.accessToken = newToken
	// }

	events := []twitch.EventSubscription{
		twitch.SubStreamOnline,
		twitch.SubStreamOffline,
		twitch.SubChannelChannelPointsCustomRewardRedemptionAdd,
		twitch.SubChannelUpdate,
	}

	tc.subscribeToEvents(message.Payload.Session.ID, streamers, events)
}

func (tc *TwitchClient) subscribeToEvents(sessionID string, streamers []TwitchStreamer, events []twitch.EventSubscription) {
	for _, streamer := range streamers {
		for _, event := range events {
			slog.Info("Subscribing to an event", "streamer", streamer, "event", event)
			_, err := twitch.SubscribeEvent(twitch.SubscribeRequest{
				SessionID:   sessionID,
				ClientID:    tc.clientId,
				AccessToken: accessToken,
				Event:       event,
				Condition: map[string]string{
					"broadcaster_user_id": streamer.id,
				},
			})

			if err != nil {
				slog.Error("Error subscribing to an event", "event", event, "error", err)
				return
			}
		}
	}
}

func (tc *TwitchClient) handleRevoke(message twitch.RevokeMessage) {
	slog.Warn("User revoked OAuth access", "userId", message.Payload.Subscription.Condition)
}

func (tc *TwitchClient) handleStreamOnline(event twitch.EventStreamOnline) {
	slog.Info("Streamer went live", "userId", event.BroadcasterUserId, "username", event.BroadcasterUserLogin)
}

func (tc *TwitchClient) handleStreamOffline(event twitch.EventStreamOffline) {
	slog.Info("Streamer went offline", "userId", event.BroadcasterUserId, "username", event.BroadcasterUserLogin)
}

func (tc *TwitchClient) handleReconnect(message twitch.ReconnectMessage) {
	slog.Warn("Twitch WebSocket reconnected", "status", message.Payload.Session.Status)
	util.SendWebHook("Twitch WebSocket reconnected, status=" + message.Payload.Session.Status)
}

func (tc *TwitchClient) handleRewardRedemption(event twitch.EventChannelChannelPointsCustomRewardRedemptionAdd) {
	// TODO: Check for duplicates
	// TODO: Check for the reward ID in the database
	slog.Info("User redeemed a reward", "userId", event.UserID, "username", event.User.UserLogin, "channel", event.BroadcasterUserLogin)
	util.SendWebHook(event.UserLogin + " redeemed an entry in " + event.BroadcasterUserLogin)
}

func (tc *TwitchClient) handleChannelUpdate(event twitch.EventChannelUpdate) {
	slog.Info("Channel updated", "channel", event.BroadcasterUserLogin)
}

func (tc *TwitchClient) handleRewardUpdate(event twitch.EventChannelChannelPointsCustomRewardUpdate) {
	slog.Warn("Streamer updated a channel point reward", "streamer", event.BroadcasterUserLogin, "reward", event.Title, "cost", event.Cost)
	util.SendWebHook(event.BroadcasterUserLogin + " updated a channel point reward, title=" + event.Title)
}

func (tc *TwitchClient) Connect() error {
	return tc.client.Connect()
}

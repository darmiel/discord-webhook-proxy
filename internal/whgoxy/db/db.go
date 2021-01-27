package db

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
)

type Database interface {
	// SaveWebhook inserts the specified webhook into the database or updates it if the _id is already present
	// returns an error if anything went wrong.
	SaveWebhook(w *discord.Webhook) (err error)

	// FindWebhook searches for a webhook by the given id (uid)
	// returns the webhook if found, otherwise an error if anything went wrong.
	FindWebhook(uid string, userID string) (w *discord.Webhook, err error)

	// FindWebhook searches for a webhook by the given id (uid) AND the matching secret
	// returns the webhook if found, otherwise an error if anything went wrong.
	FindWebhookWithSecret(uid string, userID string, secret string) (w *discord.Webhook, err error)

	// FindWebhooks returns all webhooks created by the user with the ID {userID}
	FindWebhooks(userID string) (w []*discord.Webhook, err error)

	// DeleteWebhook deletes the specified webhook if it was found.
	// Does not check further if the webhook existed before!
	DeleteWebhook(uid string, userID string) (err error)

	// CountWebhooksForUser returns the amount of webhooks
	CountWebhooksForUser(userID string) (count uint, err error)

	////

	// FindDiscordUser searches for a discord user
	FindDiscordUser(userID string) (dgd *discord.DiscordUser, err error)

	// SaveDiscordUser saves a discord user to the database
	SaveDiscordUser(dgd *discord.DiscordUser) (err error)

	////

	Disconnect() (err error)
}

var GlobalDatabase Database

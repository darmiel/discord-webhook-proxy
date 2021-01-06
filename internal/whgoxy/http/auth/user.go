package auth

import (
	discord2 "github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"golang.org/x/oauth2"
)

type User struct {
	DiscordUser *discord2.DiscordUser
	Token       *oauth2.Token
}

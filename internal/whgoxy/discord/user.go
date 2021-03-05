package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"strconv"
)

type DiscordUser struct {
	UserID        string      `json:"id" bson:"id"`
	Username      string      `json:"username" bson:"username"`
	Avatar        string      `json:"avatar" bson:"avatar""`
	Discriminator string      `json:"discriminator" bson:"discriminator"`
	PublicFlags   int         `json:"public_flags" bson:"public_flags"`
	Flags         int         `json:"flags" bson:"flags"`
	Locale        string      `json:"locale" bson:"locale"`
	MFAEnabled    bool        `json:"mfa_enabled" bson:"mfa_enabled"`
	Attributes    *attributes `json:"attributes" bson:"attributes"`
}

func (u *DiscordUser) GetFullName() string {
	return u.Username + "#" + u.Discriminator
}

func (u *DiscordUser) GetAvatarUrl() (url string) {
	url = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", u.UserID, u.Avatar)
	// system user?
	if u.Discriminator == "0000" {
		url = "/static/img/system.png"
	}
	return url
}

func NewUserByToken(token *oauth2.Token) (u *DiscordUser, err error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/v8/users/@me", nil)
	if err != nil {
		return nil, err
	}

	// Set header
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// check status
	if resp.StatusCode != 200 {
		return nil, errors.New("StatusCode was: " + strconv.Itoa(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	u = &DiscordUser{}
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *DiscordUser) Repair() (updated bool) {
	if u.Attributes == nil {
		u.Attributes = NewDefaultAttributes()
		updated = true
	}
	// repair attributes
	if u.Attributes.Repair() {
		updated = true
	}
	return
}

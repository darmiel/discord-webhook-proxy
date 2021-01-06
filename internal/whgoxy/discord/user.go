package discord

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type DiscordUser struct {
	UserID        string `json:"id"`
	Username      string `json:"username"`
	Avatar        string `json:"avatar"`
	Discriminator string `json:"discriminator"`
	PublicFlags   int    `json:"public_flags"`
	Flags         int    `json:"flags"`
	Locale        string `json:"locale"`
	MFAEnabled    bool   `json:"mfa_enabled"`
}

func (u *DiscordUser) GetAvatarUrl() (url string) {
	url = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", u.UserID, u.Avatar)
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

	log.Println("Response:", string(body))

	u = &DiscordUser{}
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, err
	}

	return u, nil
}

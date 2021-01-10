package auth

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/dchest/authcookie"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

var (
	authenticatedUsers = make(map[string]*User)
	oauthConfig        *oauth2.Config

	cookieSecret []byte
	cookieName   string
)

func InitOAuth2(conf config.OAuthConfig, router *mux.Router) {
	oauthConfig = &oauth2.Config{
		RedirectURL:  conf.RedirectURL,
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Scopes:       conf.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  conf.AuthURL,
			TokenURL: conf.TokenURL,
		},
	}

	cookieSecret = []byte(conf.CookieSecret)
	cookieName = conf.CookieName

	// routes
	router.HandleFunc("/login", handleLoginRoute)
	router.HandleFunc("/callback", handleCallbackRoute)
	router.HandleFunc("/logout", handleLogoutRoute)
}

func GetLoginCookie(r *http.Request) (value string, ok bool) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

///

func GetUser(r *http.Request) (u *User, ok bool) {
	// debug user
	// TODO: Remove me later
	u = &User{
		DiscordUser: &discord.DiscordUser{
			UserID:        "150347348088848384",
			Username:      "d2a",
			Avatar:        "408d6f884febd122f5252e2fc6d93c2e",
			Discriminator: "1325",
			PublicFlags:   256,
			Flags:         256,
			Locale:        "en-US",
			MFAEnabled:    true,
		},
		Token: &oauth2.Token{
			AccessToken:  "",
			TokenType:    "",
			RefreshToken: "",
			Expiry:       time.Time{},
		},
	}
	ok = true
	return

	// check if user sent a login cookie
	value, ok := GetLoginCookie(r)
	if !ok {
		return nil, false
	}
	log.Println("Checking if cookie is valid...")

	// check if cookie is valid
	if login := authcookie.Login(value, cookieSecret); login != "" {
		u, ok = authenticatedUsers[login]
		return u, ok
	} else {
		return nil, false
	}
}

func GetUserOrDie(r *http.Request, w http.ResponseWriter) (u *User, die bool) {
	u, ok := GetUser(r)
	if !ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, true
	}
	return u, false
}

///

func LoginUser(w http.ResponseWriter, u *User) {
	log.Println("Logging in user", u.DiscordUser.Username, "...")

	// generate cookie
	cookie := authcookie.NewSinceNow(
		u.DiscordUser.UserID,
		8*time.Hour,
		cookieSecret,
	)

	// add cookie
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   cookie,
		Expires: time.Now().Add(8 * time.Hour),
	})

	// add to map
	authenticatedUsers[u.DiscordUser.UserID] = u
}

func LogoutUser(w http.ResponseWriter, u *User) {
	log.Println("Logging out user", u.DiscordUser.Username, "...")

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: "",
	})

	delete(authenticatedUsers, u.DiscordUser.UserID)
}

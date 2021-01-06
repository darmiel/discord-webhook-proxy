package auth

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
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
	cookieSecret       []byte
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

	// routes
	router.HandleFunc("/login", handleLoginRoute)
	router.HandleFunc("/callback", handleCallbackRoute)
	router.HandleFunc("/logout", handleLogoutRoute)
}

func GetLoginCookie(r *http.Request) (value string, ok bool) {
	cookie, err := r.Cookie(config.ConfigOAuth.CookieName)
	if err != nil {
		return "", false
	}
	return cookie.Value, true
}

///

func GetUser(r *http.Request) (u *User, ok bool) {
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
		Name:    config.ConfigOAuth.CookieName,
		Value:   cookie,
		Expires: time.Now().Add(8 * time.Hour),
	})

	// add to map
	authenticatedUsers[u.DiscordUser.UserID] = u
}

func LogoutUser(w http.ResponseWriter, u *User) {
	log.Println("Logging out user", u.DiscordUser.Username, "...")

	http.SetCookie(w, &http.Cookie{
		Name:  config.ConfigOAuth.CookieName,
		Value: "",
	})

	delete(authenticatedUsers, u.DiscordUser.UserID)
}

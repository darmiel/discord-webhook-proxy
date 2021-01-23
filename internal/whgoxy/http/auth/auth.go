package auth

import (
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/dchest/authcookie"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

var (
	authUserCache = cache.New(12*time.Hour, 16*time.Hour)

	oauthConfig *oauth2.Config

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
	// check if user sent a login cookie
	value, ok := GetLoginCookie(r)
	if !ok {
		return nil, false
	}
	log.Println("Checking if cookie is valid...")

	// check if cookie is valid
	if login := authcookie.Login(value, cookieSecret); login != "" {
		// get user from cache
		var res interface{}
		res, ok = authUserCache.Get(login)
		// if found in cache: cast to user
		if ok {
			u, ok = res.(*User)
		}

		return
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
	dgu := u.DiscordUser
	log.Println("👋 Logging in ", dgu.GetFullName(), "("+dgu.UserID+")", "...")

	// generate cookie
	cookie := authcookie.NewSinceNow(
		dgu.UserID,
		8*time.Hour,
		cookieSecret,
	)

	// add cookie
	http.SetCookie(w, &http.Cookie{
		Name:    cookieName,
		Value:   cookie,
		Expires: time.Now().Add(8 * time.Hour),
	})

	// add to cache
	authUserCache.Set(dgu.UserID, u, cache.DefaultExpiration)
}

func LogoutUser(w http.ResponseWriter, u *User) {
	log.Println("🚪 Logging out ", u.DiscordUser.GetFullName(), "("+u.DiscordUser.UserID+")", "...")

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: "",
	})

	authUserCache.Delete(u.DiscordUser.UserID)
}

package auth

import (
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/config"
	"github.com/darmiel/whgoxy/internal/whgoxy/db"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"github.com/dchest/authcookie"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

var (
	AuthUserCache = cache.New(12*time.Hour, 16*time.Hour)

	oauthConfig *oauth2.Config

	cookieSecret []byte
	cookieName   string

	database db.Database
)

func InitOAuth2(conf config.OAuthConfig, router *mux.Router, db db.Database) {
	database = db

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

// TODO: Remove me later
const DebugMode = false

func GetUser(r *http.Request) (u *User, ok bool) {
	if DebugMode {
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
	}

	// check if user sent a login cookie
	value, ok := GetLoginCookie(r)
	if !ok {
		return nil, false
	}

	// check if cookie is valid
	if userID := authcookie.Login(value, cookieSecret); userID != "" {

		// get user from cache
		var res interface{}
		res, ok = AuthUserCache.Get(userID)
		// if found in cache: cast to user
		if ok {
			u, ok = res.(*User)
		} else {

			// get user from database
			dcu, _ := database.FindDiscordUser(userID)

			// login user and save to cache
			if dcu != nil && dcu.UserID == userID {
				// set valid data
				u = &User{
					DiscordUser: dcu,
					Token:       nil,
				}
				ok = true

				// add to cache
				AuthUserCache.Set(userID, u, cache.DefaultExpiration)
			}

		}

		return
	} else {
		return nil, false
	}
}

func GetUserOrDie(r *http.Request, w http.ResponseWriter) (u *User, die bool) {
	u, ok := GetUser(r)
	if !ok {
		w.WriteHeader(401)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil, true
	}
	return u, false
}

///

func LoginUser(w http.ResponseWriter, u *User) (success bool) {
	dgu := u.DiscordUser
	log.Println("👋 Logging in ", dgu.GetFullName(), "("+dgu.UserID+")", "...")

	// search in database for user
	if dbu, _ := db.GlobalDatabase.FindDiscordUser(dgu.UserID); dbu != nil {
		log.Printf("   └ User found in database. Using attributes: %+v\n", dbu.Attributes)
		// update attributes from database
		dgu.Attributes = dbu.Attributes
	}

	// check if user has permissions to login
	if !dgu.HasPermission(discord.PermissionLogin) {
		_, _ = fmt.Fprintln(w, "You don't have permissions to log in.")
		return false
	}

	// repair discord user
	if dgu.Repair() {
		log.Println("   └ User repaired.")
	}

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
	AuthUserCache.Set(dgu.UserID, u, cache.DefaultExpiration)

	// save to database
	if err := database.SaveDiscordUser(dgu); err != nil {
		log.Println("🤬 Error saving", u.DiscordUser.GetFullName(), ":", err)
	}

	return true
}

func LogoutUser(w http.ResponseWriter, u *User) {
	log.Println("🚪 Logging out ", u.DiscordUser.GetFullName(), "("+u.DiscordUser.UserID+")", "...")

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: "",
	})

	AuthUserCache.Delete(u.DiscordUser.UserID)
}

package auth

import (
	"context"
	"fmt"
	"github.com/darmiel/whgoxy/internal/whgoxy/discord"
	"math/rand"
	"net/http"
	"strings"
)

func generateState() string {
	var output strings.Builder

	charSet := "abcdedfghijklmnopqrstABCDEFGHIJKLMNOP"
	length := 20

	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteString(string(randomChar))
	}

	return output.String()
}

// handleLoginRoute creates an OAuth2 login URL and redirects the user to it.
func handleLoginRoute(w http.ResponseWriter, r *http.Request) {
	state := generateState()
	url := oauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleLogoutRoute removes the user from the cache
// and sets the cookie to an empty string ("")
func handleLogoutRoute(w http.ResponseWriter, r *http.Request) {
	user, ok := GetUser(r)
	if !ok {
		_, _ = fmt.Fprintf(w, "Error: %s", "You are not logged in.")
		return
	}

	// logout
	LogoutUser(w, user)

	// redirect
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func handleCallbackRoute(w http.ResponseWriter, r *http.Request) {
	// get discord code from query params
	code := r.URL.Query().Get("code")

	// request discord token
	token, err := oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	// parse user information
	// (id, username, ...)
	user, err := discord.NewUserByToken(token)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Invalid token or parse error: %s", err.Error())
		return
	}

	// set cookie
	LoginUser(w, &User{
		Token:       token,
		DiscordUser: user,
	})

	// redirect
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

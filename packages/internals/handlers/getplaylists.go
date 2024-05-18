package handlers

import (
	"net/http"
	"strings"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/tobyrushton/playlistpal/packages/internals/config"
	"github.com/tobyrushton/playlistpal/packages/web/templates/layouts"
	"github.com/tobyrushton/playlistpal/packages/web/templates/pages"
)

type PlaylistsHandler struct{}

func NewPlaylistsHandler() *PlaylistsHandler {
	return &PlaylistsHandler{}
}

func (h *PlaylistsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	redirectUrl := "http://localhost:8080/playlists"
	cfg := config.MustLoadConfig()

	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate),
		spotifyauth.WithClientID(cfg.SpotifyID),
	)

	if code == "" {
		url := auth.AuthURL("state")
		http.Redirect(w, r, url, http.StatusFound)
	}

	token, err := auth.Token(r.Context(), state, r)
	if err != nil {
		// incase the session has expired
		if strings.Contains(err.Error(), "oauth2: cannot fetch token: 400 Bad Request") {
			http.Redirect(w, r, "/playlists", http.StatusFound)
			return
		}
		http.Error(w, "Error getting token", http.StatusInternalServerError)
		return
	}

	client := spotify.New(auth.Client(r.Context(), token))
	user, err := client.CurrentUser(r.Context())
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	test, err := client.GetPlaylistsForUser(r.Context(), user.ID)
	if err != nil {
		http.Error(w, "Error getting playlists", http.StatusInternalServerError)
		return
	}

	c := pages.Playlists(test.Playlists)
	err = layouts.Layout(c, "Your playlists").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

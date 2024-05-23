package handlers

import (
	"net/http"

	"github.com/tobyrushton/playlistpal/internals/auth"
	"github.com/tobyrushton/playlistpal/web/templates/layouts"
	"github.com/tobyrushton/playlistpal/web/templates/pages"
)

type PlaylistsHandler struct{}

func NewPlaylistsHandler() *PlaylistsHandler {
	return &PlaylistsHandler{}
}

func (h *PlaylistsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a := auth.New()
	if !a.HasToken(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	client, err := a.GetAuthenticatedClient(r)
	if err != nil {
		http.Error(w, "Error getting client", http.StatusInternalServerError)
		return
	}

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

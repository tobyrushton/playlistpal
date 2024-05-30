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

	c := pages.Playlists()
	err := layouts.Layout(c, "Your playlists").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

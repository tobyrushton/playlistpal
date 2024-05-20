package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tobyrushton/playlistpal/packages/internals/auth"
	"github.com/tobyrushton/playlistpal/packages/web/templates/components"
	"github.com/zmb3/spotify/v2"
)

type AddHandler struct{}

func NewAddHandler() *AddHandler {
	return &AddHandler{}
}

func (h *AddHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "playlistID")
	songId := r.URL.Query().Get("songId")

	if songId == "" {
		http.Error(w, "No songId provided", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	client, err := auth.New().GetAuthenticatedClient(r)
	if err != nil {
		http.Error(w, "Error getting client", http.StatusInternalServerError)
		return
	}

	_, err = client.AddTracksToPlaylist(ctx, spotify.ID(playlistID), spotify.ID(songId))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error adding song to playlist", http.StatusInternalServerError)
		return
	}

	err = components.AddedSuccess().Render(ctx, w)
	if err != nil {
		http.Error(w, "Error rendering success message", http.StatusInternalServerError)
		return
	}
}

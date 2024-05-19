package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tobyrushton/playlistpal/packages/internals/auth"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/net/context"
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

	ctx := context.Background()

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

	w.WriteHeader(http.StatusOK)
}

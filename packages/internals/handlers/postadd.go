package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tobyrushton/playlistpal/packages/internals/config"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
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
	cfg := config.MustLoadConfig()

	config := &clientcredentials.Config{
		ClientID:     cfg.SpotifyID,
		ClientSecret: cfg.SpotifySecret,
		TokenURL:     spotifyauth.TokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		http.Error(w, "Error getting token", http.StatusInternalServerError)
		return
	}

	httpClient := spotifyauth.New().Client(ctx, token)
	client := spotify.New(httpClient)

	_, err = client.AddTracksToPlaylist(ctx, spotify.ID(playlistID), spotify.ID(songId))
	if err != nil {
		http.Error(w, "Error adding song to playlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

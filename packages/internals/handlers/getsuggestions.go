package handlers

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tobyrushton/playlistpal/packages/internals/config"
	"github.com/tobyrushton/playlistpal/packages/internals/finder"
	"github.com/tobyrushton/playlistpal/packages/web/templates/components"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type SuggestionsHandler struct{}

func NewSuggestionsHandler() *SuggestionsHandler {
	return &SuggestionsHandler{}
}

func (h *SuggestionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "playlistID")
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

	finder := finder.New(client, r, playlistID)
	suggestions, err := finder.Find()

	if err != nil {
		http.Error(w, "Error getting playlist", http.StatusInternalServerError)
		return
	}

	err = components.SongList(suggestions).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering song list", http.StatusInternalServerError)
		return
	}
}

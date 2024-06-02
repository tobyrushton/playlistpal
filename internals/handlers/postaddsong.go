package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tobyrushton/playlistpal/internals/auth"
	"github.com/tobyrushton/playlistpal/web/templates/components"
	"github.com/zmb3/spotify/v2"
)

type AddSongHandler struct{}

func NewAddSongHandler() *AddSongHandler {
	return &AddSongHandler{}
}

func (h *AddSongHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	songID := chi.URLParam(r, "songID")

	song, err := client.GetTrack(r.Context(), spotify.ID(songID))
	if err != nil {
		http.Error(w, "Error getting track", http.StatusInternalServerError)
		return
	}

	err = components.Song(spotify.SimpleTrack{
		Album:        song.Album,
		Artists:      song.Artists,
		Name:         song.Name,
		ExternalURLs: song.ExternalURLs,
		ID:           song.ID,
	}).Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

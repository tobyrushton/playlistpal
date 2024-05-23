package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tobyrushton/playlistpal/internals/auth"
	"github.com/tobyrushton/playlistpal/internals/finder"
	"github.com/tobyrushton/playlistpal/web/templates/components"
)

type SuggestionsHandler struct{}

func NewSuggestionsHandler() *SuggestionsHandler {
	return &SuggestionsHandler{}
}

func (h *SuggestionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "playlistID")

	client, err := auth.New().GetAuthenticatedClient(r)
	if err != nil {
		http.Error(w, "Error getting client", http.StatusInternalServerError)
		return
	}

	finder := finder.New(client, r, playlistID)
	suggestions, err := finder.Find()

	if err != nil || len(suggestions) == 0 {
		http.Error(w, "Error getting playlist", http.StatusInternalServerError)
		return
	}

	err = components.SongList(suggestions, playlistID).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering song list", http.StatusInternalServerError)
		return
	}
}

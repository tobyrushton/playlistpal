package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tobyrushton/playlistpal/packages/internals/auth"
	"github.com/tobyrushton/playlistpal/packages/internals/finder"
	"github.com/tobyrushton/playlistpal/packages/web/templates/components"
)

type SuggestionsHandler struct{}

func NewSuggestionsHandler() *SuggestionsHandler {
	return &SuggestionsHandler{}
}

func (h *SuggestionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	playlistID := chi.URLParam(r, "playlistID")

	client, err := auth.New().GetAuthenticatedClient(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error getting client", http.StatusInternalServerError)
		return
	}
	finder := finder.New(client, r, playlistID)
	suggestions, err := finder.Find()

	if err != nil {
		http.Error(w, "Error getting playlist", http.StatusInternalServerError)
		return
	}

	err = components.SongList(suggestions, playlistID).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering song list", http.StatusInternalServerError)
		return
	}
}

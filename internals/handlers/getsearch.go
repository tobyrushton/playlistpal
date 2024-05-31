package handlers

import (
	"net/http"

	"github.com/tobyrushton/playlistpal/internals/auth"
	"github.com/tobyrushton/playlistpal/web/templates/components"
	"github.com/zmb3/spotify/v2"
)

type SearchHandler struct{}

func NewSearchHandler() *SearchHandler {
	return &SearchHandler{}
}

func (h *SearchHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	q := r.URL.Query().Get("q")

	if q == "" {
		components.EmptySongDropDown().Render(r.Context(), w)
		return
	}

	searchResults, err := client.Search(r.Context(), q, spotify.SearchTypeTrack, spotify.Limit(5))
	if err != nil {
		http.Error(w, "Error searching", http.StatusInternalServerError)
		return
	}

	err = components.SongDropDown(searchResults.Tracks.Tracks).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

package handlers

import (
	"net/http"

	"github.com/tobyrushton/playlistpal/packages/web/templates/layouts"
	"github.com/tobyrushton/playlistpal/packages/web/templates/pages"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := pages.Index()
	err := layouts.Layout(c, "PlaylistPal").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

package handlers

import (
	"net/http"

	"github.com/tobyrushton/playlistpal/packages/internals/auth"
)

type LoginHandler struct{}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a := auth.New()

	state := r.URL.Query().Get("state")

	if state != "" {
		token, err := a.GetAuthClient().Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Error getting token", http.StatusInternalServerError)
			return
		}
		a.SetToken(w, token)
		http.Redirect(w, r, "/playlists", http.StatusFound)
		return
	}

	if a.HasToken(r) {
		http.Redirect(w, r, "/playlists", http.StatusFound)
		return
	}

	url := auth.New().GetAuthUrl()
	http.Redirect(w, r, url, http.StatusFound)
}

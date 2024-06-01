package handlers

import (
	"net/http"
	"strconv"

	"github.com/tobyrushton/playlistpal/internals/auth"
	"github.com/tobyrushton/playlistpal/internals/finder"
	"github.com/tobyrushton/playlistpal/web/templates/components"
	"github.com/zmb3/spotify/v2"
)

type CreatePlaylistHandler struct{}

func NewCreatePlaylistHandler() *CreatePlaylistHandler {
	return &CreatePlaylistHandler{}
}

func (h *CreatePlaylistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	r.ParseForm()
	songs := r.Form["songs"]
	name := r.FormValue("name")
	amount, _ := strconv.Atoi(r.FormValue("amount"))

	if songs[0] == "undefined" {
		http.Error(w, "Select a song to get started", http.StatusBadRequest)
		return
	}

	if name == "" {
		http.Error(w, "Choose a name for your playlist", http.StatusBadRequest)
		return
	}

	if amount > 50 || amount < 1 {
		http.Error(w, "Invalid amount of songs", http.StatusBadRequest)
		return
	}

	user, err := client.CurrentUser(r.Context())
	if err != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	// Create playlist
	playlist, err := client.CreatePlaylistForUser(r.Context(), user.ID, name, "", true, false)
	if err != nil {
		http.Error(w, "Error creating playlist", http.StatusInternalServerError)
		return
	}

	// get spotify IDs
	trackIDs := make([]spotify.ID, 0)
	for _, song := range songs {
		trackIDs = append(trackIDs, spotify.ID(song))
	}

	client.AddTracksToPlaylist(r.Context(), playlist.ID, trackIDs...)

	f := finder.New(client, r, string(playlist.ID))
	recs, err := f.FillNewPlaylist(amount)

	if err != nil {
		http.Error(w, "Error getting recommendations", http.StatusInternalServerError)
		return
	}

	// Add recommendations to playlist
	trackIDs = make([]spotify.ID, len(recs))
	for i, rec := range recs {
		trackIDs[i] = rec.ID
	}
	_, err = client.AddTracksToPlaylist(r.Context(), playlist.ID, trackIDs...)
	if err != nil {
		http.Error(w, "Error adding tracks to playlist", http.StatusInternalServerError)
		return
	}

	// update playlist to get the new image
	playlist, _ = client.GetPlaylist(r.Context(), playlist.ID)

	err = components.PlaylistSimple(*playlist).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

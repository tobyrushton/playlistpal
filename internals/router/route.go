package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tobyrushton/playlistpal/internals/handlers"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	fileServer := http.FileServer(http.Dir("./web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Group(func(r chi.Router) {
		r.Get("/", handlers.NewHomeHandler().ServeHTTP)
		r.Get("/playlists", handlers.NewPlaylistsHandler().ServeHTTP)
		r.Get("/api/suggestions/{playlistID}", handlers.NewSuggestionsHandler().ServeHTTP)
		r.Post("/api/add/{playlistID}", handlers.NewAddHandler().ServeHTTP)
		r.Get("/api/playlists", handlers.NewPlaylistsListHandler().ServeHTTP)
		r.Get("/login", handlers.NewLoginHandler().ServeHTTP)
		r.Get("/api/search", handlers.NewSearchHandler().ServeHTTP)
		r.Post("/api/add-song/{songID}", handlers.NewAddSongHandler().ServeHTTP)
		r.Post("/api/create-playlist", handlers.NewCreatePlaylistHandler().ServeHTTP)
	})

	return r
}

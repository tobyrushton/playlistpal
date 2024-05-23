package api

import (
	"net/http"

	"github.com/tobyrushton/playlistpal/internals/router"
)

func Handler(w http.ResponseWriter, req *http.Request) {
	r := router.NewRouter()
	r.ServeHTTP(w, req)
}

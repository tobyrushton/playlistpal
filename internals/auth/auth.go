package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/tobyrushton/playlistpal/internals/config"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type Auth struct {
	RedirectUrl string
}

func New() *Auth {
	if os.Getenv("env") == "production" {
		return &Auth{
			RedirectUrl: "https://playlistpal.tobyrushton.com/login",
		}
	} else {
		return &Auth{
			RedirectUrl: "http://localhost:8080/login",
		}
	}
}

func (a *Auth) GetAuthClient() *spotifyauth.Authenticator {
	cfg := config.MustLoadConfig()

	return spotifyauth.New(
		spotifyauth.WithRedirectURL(a.RedirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopeUserReadPrivate, spotifyauth.ScopePlaylistModifyPrivate, spotifyauth.ScopePlaylistModifyPublic),
		spotifyauth.WithClientID(cfg.SpotifyID),
		spotifyauth.WithClientSecret(cfg.SpotifySecret),
	)
}

func (a *Auth) GetAuthUrl() string {
	auth := a.GetAuthClient()

	return auth.AuthURL(a.GenerateState())
}

func (a *Auth) GenerateState() string {
	bytes := make([]byte, 128)
	rand.Read(bytes)

	return base64.URLEncoding.EncodeToString(bytes)[:128]
}

func (a *Auth) GetAuthenticatedClient(r *http.Request) (*spotify.Client, error) {
	if !a.HasToken(r) {
		return nil, errors.New("no token")
	}
	token, err := a.GetToken(r)
	if err != nil {
		return nil, err
	}

	httpClient := a.GetAuthClient().Client(r.Context(), token)
	return spotify.New(httpClient), nil
}

// checks cookies for token
func (a *Auth) HasToken(r *http.Request) bool {
	_, err := r.Cookie("spotify_token")
	return err == nil
}

// sets token in cookies
func (a *Auth) SetToken(w http.ResponseWriter, token *oauth2.Token) {
	jsonToken, _ := json.Marshal(token)
	encodedJson := base64.StdEncoding.EncodeToString(jsonToken)
	http.SetCookie(w, &http.Cookie{
		Name:     "spotify_token",
		Value:    encodedJson,
		HttpOnly: true,
		Expires:  token.Expiry,
	})
}

func (a *Auth) GetToken(r *http.Request) (*oauth2.Token, error) {
	cookie, err := r.Cookie("spotify_token")
	if err != nil {
		return nil, err
	}

	decodedJson, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	var token oauth2.Token
	err = json.Unmarshal([]byte(decodedJson), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

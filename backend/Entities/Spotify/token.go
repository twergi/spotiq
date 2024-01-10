package Spotify

import (
	"fmt"
	"log"
	"net/http"
	entities "spotiq_backend/Entities"
	"strings"
	"time"
)

var spotifyTokenURL string = "https://accounts.spotify.com/api/token"
var spotifyLoginURL string = "https://accounts.spotify.com/authorize"

type TokenData struct {
	AccessToken  string
	RefreshToken string
	TokenType    string
	ExpiresAt    time.Time
}

type SpotifyToken interface {
	GetData() TokenData
	CreateAuthURL() string
	Refresh()
	IsExpired() bool
	IsValid() bool
}

type spotifyToken struct {
	tokenURL string
	loginURL string
	data     TokenData
}

func (sToken *spotifyToken) GetData() TokenData {
	return sToken.data
}

func (sToken *spotifyToken) CreateAuthURL() string {
	var responseType string = "code"
	var redirectURI = fmt.Sprintf("%s/callback/", entities.RedirectURI)
	var scopes []string = []string{
		"user-read-playback-state",
		"user-read-currently-playing",
		"user-modify-playback-state",
	}
	var scope string = strings.Join(scopes, " ")

	request, err := http.NewRequest("GET", sToken.loginURL, nil)

	if err != nil {
		log.Panicln(err.Error())
	}

	query := request.URL.Query()
	query.Add("client_id", entities.ClientID)
	query.Add("scope", scope)
	query.Add("response_type", responseType)
	query.Add("redirect_uri", redirectURI)

	request.URL.RawQuery = query.Encode()
	return request.URL.String()
}

func (sToken *spotifyToken) IsExpired() bool {
	return sToken.data.ExpiresAt.Add(-time.Second * 10).Before(time.Now())
}

func (sToken *spotifyToken) IsValid() bool {
	if sToken.data.AccessToken == "" {
		return false
	}

	var expired bool = sToken.IsExpired()

	if expired {
		sToken.Refresh()
		expired = sToken.IsExpired()
	}

	return expired
}

func (sToken *spotifyToken) Refresh() {
	var grantType string = "refresh_token"

	response, err := http.Post(
		sToken.tokenURL,
		"application/x-www-form-urlencoded",
	)
}

func NewToken() SpotifyToken {
	return &spotifyToken{
		tokenURL: spotifyTokenURL,
		loginURL: spotifyLoginURL,
		data:     TokenData{},
	}
}

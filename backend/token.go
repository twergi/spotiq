package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	tokenURL = "https://accounts.spotify.com/api/token"
	loginURL = "https://accounts.spotify.com/authorize"
)

type Token struct {
	access    string
	refresh   string
	ttype     string
	expiresAt time.Time

	redirectURI            string
	clientID, clientSecret string
}

func NewToken(redirectURI, clientID, clientSecret string) *Token {
	return &Token{
		redirectURI:  redirectURI,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

func (t *Token) createAuthURL() (string, error) {
	responseType := "code"
	redirectURI := t.redirectURI + "/callback"
	scopes := []string{
		"user-read-playback-state",
		"user-read-currently-playing",
		"user-modify-playback-state",
	}

	authURL, err := url.Parse(loginURL)
	if err != nil {
		return "", err
	}

	query := authURL.Query()
	query.Add("scope", strings.Join(scopes, " "))
	query.Add("response_type", responseType)
	query.Add("redirect_uri", redirectURI)
	query.Add("client_id", t.clientID)

	authURL.RawQuery = query.Encode()

	return authURL.String(), nil
}

func (t *Token) ready() bool {
	return !t.expiresAt.IsZero()
}

func (t *Token) expired() bool {
	if !t.ready() {
		return true
	}

	return time.Now().Add(-time.Second * 10).After(t.expiresAt)
}

func (t *Token) getData() (map[string]any, bool) {
	if !t.ready() {
		return nil, false
	}

	return map[string]any{
		"access_token":  t.access,
		"token_type":    t.ttype,
		"refresh_token": t.refresh,
		"expires_at":    t.expiresAt,
	}, true
}

func (t *Token) validateToken() error {
	if !t.ready() {
		return nil
	}

	if t.expired() {
		return t.refreshToken()
	}

	return nil
}

type refreshResponseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Error        string `json:"error"`
}

func (t *Token) refreshToken() error {
	grantType := "refresh_token"

	bodyData := map[string]any{
		"grant_type":    grantType,
		"refresh_token": t.refresh,
	}
	bodyRaw, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}

	resp, err := http.Post(tokenURL, "application/json", bytes.NewReader(bodyRaw))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("response status code not 200:" + string(body))
	}

	responseBody := refreshResponseBody{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return err
	}

	if responseBody.Error != "" {
		return errors.New(responseBody.Error)
	}

	if responseBody.ExpiresIn != 0 {
		t.expiresAt = _calculateExpiresAt(time.Second * time.Duration(responseBody.ExpiresIn))
	}
	if responseBody.AccessToken != "" {
		t.access = responseBody.AccessToken
	}
	if responseBody.RefreshToken != "" {
		t.refresh = responseBody.RefreshToken
	}
	if responseBody.TokenType != "" {
		t.ttype = responseBody.TokenType
	}

	return nil
}

func (t *Token) _createTokenData(code string) error {
	grantType := "authorization_code"
	redirectURI := t.redirectURI + "/callback"

	bodyData := url.Values{
		"client_id":     []string{t.clientID},
		"client_secret": []string{t.clientSecret},
		"redirect_uri":  []string{redirectURI},
		"grant_type":    []string{grantType},
		"code":          []string{code},
	}

	resp, err := http.PostForm(tokenURL, bodyData)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("response status code not 200:" + string(body))
	}

	responseBody := refreshResponseBody{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return err
	}

	if responseBody.Error != "" {
		return errors.New(responseBody.Error)
	}

	t.expiresAt = _calculateExpiresAt(time.Second * time.Duration(responseBody.ExpiresIn))
	t.access = responseBody.AccessToken
	t.refresh = responseBody.RefreshToken
	t.ttype = responseBody.TokenType

	return nil
}

func _calculateExpiresAt(expiresIn time.Duration) time.Time {
	return time.Now().Add(expiresIn)
}

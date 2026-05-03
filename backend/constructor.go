package main

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
)

var (
	authURL *url.URL
)

func init() {
	var err error

	authURL, err = url.Parse("https://accounts.spotify.com/authorize")
	if err != nil {
		panic(err)
	}
}

type HTTPConstructor struct {
	conf *Config
}

func NewConstructor(
	conf *Config,
) *HTTPConstructor {
	return &HTTPConstructor{
		conf: conf,
	}
}

func (c *HTTPConstructor) AuthURL() string {
	query := authURL.Query()
	query.Add("scope", strings.Join([]string{
		"user-read-playback-state",
		"user-read-currently-playing",
		"user-modify-playback-state",
	}, " "))
	query.Add("response_type", "code")
	query.Add("redirect_uri", c.conf.RedirectURI+"/callback")
	query.Add("client_id", c.conf.ClientID)

	authURL.RawQuery = query.Encode()

	return authURL.String()
}

func (c *HTTPConstructor) CallbackURL() string {
	return c.conf.RedirectURI + "/callback"
}

func (c *HTTPConstructor) TokenURL() string {
	return "https://accounts.spotify.com/api/token"
}

func (c *HTTPConstructor) FinishLoginRequest(ctx context.Context, code string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.TokenURL(), strings.NewReader(url.Values{
		"client_id":     []string{c.conf.ClientID},
		"client_secret": []string{c.conf.ClientSecret},
		"redirect_uri":  []string{c.CallbackURL()},
		"grant_type":    []string{"authorization_code"},
		"code":          []string{code},
	}.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (c *HTTPConstructor) BaseURL() string {
	return "https://api.spotify.com/v1"
}

func (c *HTTPConstructor) GetCurrentQueueRequest(ctx context.Context, accessToken string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL()+"/me/player/queue", nil)
	if err != nil {
		return nil, err
	}

	c.authorizeBearer(req, accessToken)

	return req, nil
}

func (c *HTTPConstructor) AddToQueueRequest(ctx context.Context, songURI, deviceID, accessToken string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL()+"/me/player/queue", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("uri", songURI)
	query.Add("device_id", deviceID)
	req.URL.RawQuery = query.Encode()

	c.authorizeBearer(req, accessToken)

	return req, nil
}

func (c *HTTPConstructor) authorizeBearer(req *http.Request, accessToken string) {
	req.Header.Add("Authorization", "Bearer "+accessToken)
}

func (c *HTTPConstructor) authorizeBasic(req *http.Request) {
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.conf.ClientID+":"+c.conf.ClientSecret)))
}

func (c *HTTPConstructor) SearchRequest(ctx context.Context, query, accessToken string) (*http.Request, error) {
	url, err := url.Parse(c.BaseURL() + "/search")
	if err != nil {
		return nil, err
	}

	q := url.Query()
	q.Add("q", query)
	q.Add("type", "track")
	q.Add("limit", "10")

	url.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	c.authorizeBearer(req, accessToken)

	return req, nil
}

func (c *HTTPConstructor) GetDevicesRequest(ctx context.Context, accessToken string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL()+"/me/player/devices", nil)
	if err != nil {
		return nil, err
	}

	c.authorizeBearer(req, accessToken)

	return req, nil
}

func (c *HTTPConstructor) RefreshTokenRequest(ctx context.Context, refreshToken string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.TokenURL(), strings.NewReader(url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refreshToken},
	}.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.authorizeBasic(req)

	return req, nil
}

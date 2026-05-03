package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

const (
	baseURL = "https://api.spotify.com/v1"
)

type Router struct {
	token  *Token
	device *Device
}

func NewRouter(redirectURI, clientID, clientSecret string) *Router {
	return &Router{
		token:  NewToken(redirectURI, clientID, clientSecret),
		device: NewDevice(),
	}
}

func (r *Router) clearData() {
	r.token = &Token{}
	r.device = &Device{}
}

func (r *Router) _makeAuthHeaders() (name, value string) {
	return "Authorization", "Bearer " + r.token.access
}

func (r *Router) getData() map[string]any {
	resp := make(map[string]any, 2)

	if data, ok := r.token.getData(); ok {
		resp["token"] = data
	} else {
		resp["token"] = nil
	}

	if data, ok := r.device.getData(); ok {
		resp["device"] = data
	} else {
		resp["device"] = nil
	}

	return resp
}

func (r *Router) createTokenData(code string) error {
	return r.token._createTokenData(code)
}

func (r *Router) getDevices() ([]byte, error) {
	r.token.validateToken()

	req, err := http.NewRequest(http.MethodGet, baseURL+"/me/player/devices", nil)
	if err != nil {
		return nil, err
	}

	authHeader, authHValue := r._makeAuthHeaders()
	req.Header.Add(authHeader, authHValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code not 200:" + string(body))
	}

	return body, nil
}

func (r *Router) getCurrentQueue() ([]byte, error) {
	r.token.validateToken()

	url := baseURL + "/me/player/queue"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	authHeader, authHValue := r._makeAuthHeaders()
	req.Header.Add(authHeader, authHValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code not 200:" + string(body))
	}

	return body, nil
}

type SongInfo struct {
	URI string
}

func (r *Router) addToQueue(songInfo SongInfo) error {
	r.token.validateToken()

	url := baseURL + "/me/player/queue"

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	query := req.URL.Query()
	query.Add("uri", songInfo.URI)
	query.Add("device_id", r.device.id)
	req.URL.RawQuery = query.Encode()

	authHeader, authHValue := r._makeAuthHeaders()
	req.Header.Set(authHeader, authHValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("response status code not 200:" + string(body) + ", status: " + resp.Status)
	}

	return nil
}

func (r *Router) getCurrentlyPlaying() ([]byte, error) {
	r.token.validateToken()

	url := baseURL + "/me/player/currently-playing"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	authHeader, authHValue := r._makeAuthHeaders()
	req.Header.Add(authHeader, authHValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code not 200:" + string(body))
	}

	return body, nil
}

type SearchResponseBody struct {
	Error  string `json:"error"`
	Tracks []any  `json:"tracks"`
}

func (r *Router) search(query string) ([]byte, error) {
	r.token.validateToken()

	url, err := url.Parse(baseURL + "/search")
	if err != nil {
		return nil, err
	}

	q := url.Query()
	q.Add("q", query)
	q.Add("type", "track")
	q.Add("limit", "10")

	url.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}

	authHeader, authHValue := r._makeAuthHeaders()
	req.Header.Add(authHeader, authHValue)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respBody := map[string]json.RawMessage{}
	err = json.Unmarshal(bodyRaw, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody["tracks"], nil
}

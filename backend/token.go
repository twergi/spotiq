package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type TokenService struct {
	access    string
	refresh   string
	ttype     string
	expiresAt time.Time

	conf   *Config
	c      *HTTPConstructor
	client *http.Client

	lock   sync.RWMutex
	ctx    context.Context
	cancel func()
}

func NewTokenService(
	ctx context.Context,
	conf *Config,
	c *HTTPConstructor,
	client *http.Client,
) *TokenService {
	return &TokenService{
		conf:   conf,
		c:      c,
		ctx:    ctx,
		client: client,
	}
}

func (t *TokenService) Clear() {
	if !t.Ready() {
		return
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	t.cancel()
	t.access = ""
	t.refresh = ""
	t.ttype = ""
	t.expiresAt = time.Time{}
	t.cancel = nil
}

func (t *TokenService) Ready() bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return !t.expired()
}

func (t *TokenService) expired() bool {
	return t.expiresAt.IsZero() ||
		time.Now().Add(-time.Second*10).After(t.expiresAt)
}

func (t *TokenService) refreshToken(ctx context.Context) error {
	t.lock.RLock()

	req, err := t.c.RefreshTokenRequest(ctx, t.refresh)
	if err != nil {
		return err
	}

	t.lock.RUnlock()

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer discard(req.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("refresh token error:" + string(body))
	}

	responseBody := AuthPayload{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return err
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	if responseBody.Error != "" {
		return errors.New(responseBody.Error)
	}

	payload := responseBody.TokenPayload
	if payload.RefreshToken == "" {
		payload.RefreshToken = t.refresh
	}

	t.setToken(payload)

	return nil
}

func (t *TokenService) SetToken(payload TokenPayload) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.setToken(payload)
}

func (t *TokenService) setToken(payload TokenPayload) {
	t.expiresAt = time.Now().Add(time.Second * time.Duration(payload.ExpiresIn))
	t.access = payload.AccessToken
	t.refresh = payload.RefreshToken
	t.ttype = payload.TokenType
	t.setupRefreshTimer()
}

func (t *TokenService) setupRefreshTimer() {
	if t.cancel != nil {
		t.cancel()
	}

	ctx, cancel := context.WithCancel(t.ctx)
	t.cancel = cancel
	deadline := time.Until(t.expiresAt.Add(-time.Minute * 2))

	go func() {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "refresh timer context cancelled")
			return
		case <-time.After(deadline):
			for {
				slog.InfoContext(ctx, "refreshing token")
				err := t.refreshToken(ctx)
				if err != nil {
					slog.ErrorContext(ctx, "error refreshing token, retrying...", "error", err.Error())
					time.Sleep(time.Second * 5)
					if ctx.Err() != nil {
						return
					}
					continue
				}
				slog.InfoContext(ctx, "token refreshed")
				return
			}
		}
	}()
}

func (t *TokenService) GetAccessToken() string {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.access
}

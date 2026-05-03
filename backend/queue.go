package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type Queue struct {
	tokenService  *TokenService
	deviceService *DeviceService
	c             *HTTPConstructor
	conf          *Config
	client        *http.Client

	currentQueue CurrentQueue
}

func NewQueue(
	tokenService *TokenService,
	deviceService *DeviceService,
	c *HTTPConstructor,
	conf *Config,
	client *http.Client,
) *Queue {
	return &Queue{
		tokenService:  tokenService,
		deviceService: deviceService,
		c:             c,
		conf:          conf,
		client:        client,
	}
}

func (q *Queue) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "queue refresher context cancelled")
			return ctx.Err()
		case <-time.After(time.Second * 15):
		}

		slog.InfoContext(ctx, "refreshing queue...")
		err := q.refreshQueue(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "error refreshing queue", "error", err.Error())
			continue
		}
		slog.InfoContext(ctx, "queue refreshed")
	}
}

func (q *Queue) GetState() QueueState {
	return QueueState{
		TokenSet:  q.tokenService.Ready(),
		DeviceSet: q.deviceService.Ready(),
	}
}

func (q *Queue) ClearState() {
	q.tokenService.Clear()
	q.deviceService.Clear()
}

func (q *Queue) FinishLogin(ctx context.Context, code string) error {
	req, err := q.c.FinishLoginRequest(ctx, code)
	if err != nil {
		return err
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	responseBody := AuthPayload{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK || responseBody.Error != "" {
		errstr := responseBody.Error
		if errstr == "" {
			errstr = string(body)
		}

		return errors.New("finish login error: " + errstr)
	}

	q.tokenService.SetToken(responseBody.TokenPayload)

	return nil
}

func (q *Queue) GetQueue(ctx context.Context) CurrentQueue {
	return q.currentQueue
}

func (q *Queue) refreshQueue(ctx context.Context) error {
	if !q.ready() {
		slog.InfoContext(ctx, "not yet ready to refresh queue")
		return nil
	}

	req, err := q.c.GetCurrentQueueRequest(ctx, q.tokenService.GetAccessToken())
	if err != nil {
		return err
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer discard(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("error retrieving current queue:" + string(body))
	}

	res := CurrentQueue{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	q.currentQueue = res

	return nil
}

func (q *Queue) ready() bool {
	state := q.GetState()
	return state.TokenSet && state.DeviceSet
}

func (q *Queue) Enqueue(ctx context.Context, songURI string) error {
	req, err := q.c.AddToQueueRequest(
		ctx,
		songURI,
		q.deviceService.GetCurrentDeviceID(),
		q.tokenService.GetAccessToken(),
	)
	if err != nil {
		return err
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return err
	}
	defer discard(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return errors.New("add to queue error:" + string(body))
	}

	return nil
}

func (q *Queue) Search(ctx context.Context, query string) ([]Track, error) {
	req, err := q.c.SearchRequest(ctx, query, q.tokenService.GetAccessToken())
	if err != nil {
		return nil, err
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer discard(resp.Body)

	bodyRaw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	payload := SearchPayload{}
	err = json.Unmarshal(bodyRaw, &payload)
	if err != nil {
		return nil, err
	}

	return payload.Tracks.Items, nil
}

func (q *Queue) GetDevices(ctx context.Context) ([]Device, error) {
	req, err := q.c.GetDevicesRequest(ctx, q.tokenService.GetAccessToken())
	if err != nil {
		return nil, err
	}

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer discard(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("get devices error:" + string(body))
	}

	var payload DevicesPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		return nil, err
	}

	return payload.Devices, nil
}

func (q *Queue) SetDevice(device Device) {
	q.deviceService.SetDevice(device)
}

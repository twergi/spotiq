package main

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
)

type Router struct {
	mux *chi.Mux

	q    *Queue
	c    *HTTPConstructor
	conf *Config
}

func NewRouter(
	conf *Config,
	q *Queue,
	c *HTTPConstructor,
) *Router {
	r := &Router{
		q:    q,
		c:    c,
		conf: conf,
	}

	r.initMux(conf)

	return r
}

func (rr *Router) initMux(conf *Config) {
	mux := chi.NewMux()

	cors := cors.New(cors.Options{
		AllowedOrigins:   conf.CORSAllowedOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	mux.Use(cors.Handler)

	mux.Get("/state", rr.getState)
	mux.Get("/login", rr.login)
	mux.Get("/logout", rr.logout)
	mux.Get("/callback", rr.callback)
	mux.Get("/current_queue", rr.getQueue)
	mux.Post("/current_queue", rr.enqueue)
	mux.Get("/search", rr.search)
	mux.Get("/devices", rr.getDevices)
	mux.Post("/devices", rr.setDevice)

	rr.mux = mux
}

func (rr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rr.mux.ServeHTTP(w, r)
}

func (rr *Router) getState(w http.ResponseWriter, r *http.Request) {
	state := rr.q.GetState()
	body, err := json.Marshal(state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func (rr *Router) login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, rr.c.AuthURL(), http.StatusTemporaryRedirect)
}

func (rr *Router) logout(w http.ResponseWriter, r *http.Request) {
	rr.q.ClearState()
	http.Redirect(w, r, rr.conf.FrontURL, http.StatusTemporaryRedirect)
}

func (rr *Router) callback(w http.ResponseWriter, r *http.Request) {
	err := rr.q.FinishLogin(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, rr.conf.FrontURL, http.StatusTemporaryRedirect)
}

func (rr *Router) getQueue(w http.ResponseWriter, r *http.Request) {
	queue := rr.q.GetQueue(r.Context())

	body, err := json.Marshal(queue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func (rr *Router) enqueue(w http.ResponseWriter, r *http.Request) {
	defer discard(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	song := EnqueuePayload{}
	err = json.Unmarshal(body, &song)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = rr.q.Enqueue(r.Context(), song.URI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (rr *Router) search(w http.ResponseWriter, r *http.Request) {
	tracks, err := rr.q.Search(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(tracks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func (rr *Router) getDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := rr.q.GetDevices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(devices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(body)
}

func (rr *Router) setDevice(w http.ResponseWriter, r *http.Request) {
	data := Device{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer discard(r.Body)

	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rr.q.SetDevice(data)

	w.WriteHeader(http.StatusOK)
}

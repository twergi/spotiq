package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
)

func main() {
	var (
		backendIP          = os.Getenv("BACKEND_IP")
		backendPort        = os.Getenv("BACKEND_PORT")
		backendPortForward = os.Getenv("BACKEND_PORT_FORWARD")
		frontIP            = os.Getenv("FRONTEND_IP")
		frontPort          = os.Getenv("FRONTEND_PORT_FORWARD")
		frontURL           = "http://" + frontIP + ":" + frontPort
		redirectURI        = "http://" + frontIP + ":" + backendPortForward
		clientID           = os.Getenv("CLIENT_ID")
		clientSecret       = os.Getenv("CLIENT_SECRET")
	)

	api := NewRouter(redirectURI, clientID, clientSecret)

	mux := chi.NewMux()
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1:3000",
			"127.0.0.1:3000",
			"http://127.0.0.1",
			"127.0.0.1",
			frontIP + ":" + frontPort,
			"http://" + frontIP + ":" + frontPort,
			frontIP,
			"http://" + frontIP,
		},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	mux.Use(cors.Handler)

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := json.Marshal(api.getData())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(body)
	})

	mux.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		url, err := api.token.createAuthURL()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	mux.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		api.clearData()
		http.Redirect(w, r, frontURL, http.StatusTemporaryRedirect)
	})

	mux.Get("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		err := api.createTokenData(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, frontURL, http.StatusTemporaryRedirect)
	})

	mux.Get("/current_queue", func(w http.ResponseWriter, r *http.Request) {
		body, err := api.getCurrentQueue()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(body)
	})

	mux.Post("/current_queue", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		song := SongInfo{}
		err = json.Unmarshal(body, &song)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = api.addToQueue(song)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	mux.Get("/currently_playing", func(w http.ResponseWriter, r *http.Request) {
		body, err := api.getCurrentlyPlaying()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(body)
	})

	mux.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		body, err := api.search(r.URL.Query().Get("q"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(body)
	})

	mux.Get("/devices", func(w http.ResponseWriter, r *http.Request) {
		body, err := api.getDevices()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(body)
	})

	mux.Post("/devices", func(w http.ResponseWriter, r *http.Request) {
		data := DeviceData{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		api.device.setData(data)

		w.WriteHeader(http.StatusOK)
	})

	server := http.Server{}

	server.Handler = mux

	server.Addr = backendIP + ":" + backendPort

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	errgr, errctx := errgroup.WithContext(ctx)
	errgr.Go(func() error {
		<-errctx.Done()
		return server.Shutdown(ctx)
	})

	errgr.Go(func() error {
		slog.InfoContext(errctx, fmt.Sprintf("running server on: %s", server.Addr))
		return server.ListenAndServe()
	})

	err := errgr.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}

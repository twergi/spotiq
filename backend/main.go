package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	client := http.DefaultClient
	conf := NewConfigFromEnv()
	constructor := NewConstructor(conf)
	token := NewTokenService(ctx, conf, constructor, client)
	device := NewDeviceService()
	queue := NewQueue(token, device, constructor, conf, client)
	router := NewRouter(conf, queue, constructor)
	server := NewServer(router, conf.BackendHostPort)

	errgr, errctx := errgroup.WithContext(ctx)
	errgr.Go(func() error {
		<-errctx.Done()
		return server.Shutdown(ctx)
	})

	errgr.Go(func() error {
		slog.InfoContext(errctx, fmt.Sprintf("running server on: %s", server.Addr))
		return server.ListenAndServe()
	})

	errgr.Go(func() error {
		return queue.Start(ctx)
	})

	err := errgr.Wait()
	if err != nil {
		log.Fatalln(err)
	}
}

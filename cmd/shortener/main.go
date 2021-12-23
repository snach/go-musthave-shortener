package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"snach/go-musthave-shortener/cmd/shortener/config"
	"snach/go-musthave-shortener/cmd/shortener/handlers"
	"snach/go-musthave-shortener/cmd/shortener/repository"
	"syscall"
	"time"
)

func serve(ctx context.Context) (err error) {
	conf, err := config.MakeConf()
	if err != nil {
		panic(err)
	}

	repo, err := repository.NewRepository(conf.FileStoragePath)
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    conf.ServerAddress,
		Handler: handlers.NewRouter(conf.BaseURL, repo),
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started")
	<-ctx.Done()
	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed: %+s", err)
	}

	log.Printf("server exited properly")
	if err == http.ErrServerClosed {
		err = nil
	}
	return
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()

	if err := serve(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}

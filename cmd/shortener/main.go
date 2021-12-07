package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"snach/go-musthave-shortener/cmd/shortener/handlers"
	"time"
)

var shortToFull = make(map[int]string)
var mapCounter = 1

func serve(ctx context.Context) (err error) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handlers.NewRouter(shortToFull, mapCounter),
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
	signal.Notify(c, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()

	if err := serve(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
}

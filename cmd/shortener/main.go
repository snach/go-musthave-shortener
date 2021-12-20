package main

import (
	"context"
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
	"net/http"
	"os"
	"os/signal"
	"snach/go-musthave-shortener/cmd/shortener/handlers"
	"snach/go-musthave-shortener/cmd/shortener/repository"
	"syscall"
	"time"
)

type ServerConfig struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"storage.txt"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func makeConf() (ServerConfig, error) {
	var conf ServerConfig
	if err := env.Parse(&conf); err != nil {
		return ServerConfig{}, err
	}

	address := flag.String("a", conf.ServerAddress, "address server (or env var SERVER_ADDRESS)")
	fileStoragePath := flag.String("f", conf.FileStoragePath, "path to storage file (or env var FILE_STORAGE_PATH)")
	baseUrl := flag.String("b", conf.BaseURL, "base url ajh shortened link (or env var BASE_URL)")
	flag.Parse()

	return ServerConfig{
		ServerAddress:   *address,
		FileStoragePath: *fileStoragePath,
		BaseURL:         *baseUrl,
	}, nil
}

func serve(ctx context.Context) (err error) {
	conf, err := makeConf()
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

package main

import (
	"log"
	"net/http"
	"snach/go-musthave-shortener/cmd/shortener/handlers"
)

func main() {
	var shortToFull = make(map[int]string)
	var mapCounter = 1

	log.Fatal(http.ListenAndServe(":8080", handlers.NewRouter(shortToFull, mapCounter)))
}

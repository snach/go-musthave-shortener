package main

import (
	"log"
	"net/http"
	"snach/go-musthave-shortener/cmd/shortener/handlers"
)

func main() {
	var shortToFull = make(map[int]string)
	var mapCounter = 1

	http.HandleFunc("/", handlers.ShortenerHandler(shortToFull, mapCounter))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

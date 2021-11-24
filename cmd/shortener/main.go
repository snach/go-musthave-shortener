package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var shortToFull = make(map[int]string)
var mapCounter = 1

func ShortenerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		shortToFull[mapCounter] = string(body)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "http://localhost:8080/%d", mapCounter)

		mapCounter++
		return

	} else if r.Method == http.MethodGet {
		shortId, err := strconv.Atoi(strings.Trim(r.RequestURI, "/"))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if fullUrl, ok := shortToFull[shortId]; ok {
			w.Header().Set("Location", fullUrl)
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
	}

	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	// маршрутизация запросов обработчику
	http.HandleFunc("/", ShortenerHandler)
	// запуск сервера с адресом localhost, порт 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}

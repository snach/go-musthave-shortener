package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func ShortenerHandler(shortToFull map[int]string, mapCounter int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			shortToFull[mapCounter] = string(body)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, "http://localhost:8080/%d", mapCounter)

			mapCounter++
			return

		} else if r.Method == http.MethodGet {
			shortID, err := strconv.Atoi(strings.Trim(r.RequestURI, "/"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if fullURL, ok := shortToFull[shortID]; ok {
				w.Header().Set("Location", fullURL)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}
		}

		w.WriteHeader(http.StatusBadRequest)
	}
}

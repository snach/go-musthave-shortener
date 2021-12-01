package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"strconv"
)

func NewRouter(shortToFull map[int]string, mapCounter int) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/", CreateShortUrlHandler(shortToFull, mapCounter))
	r.Get("/{id}", GetFullUrlHandler(shortToFull))
	return r
}

func CreateShortUrlHandler(shortToFull map[int]string, mapCounter int) http.HandlerFunc {
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
		}
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetFullUrlHandler(shortToFull map[int]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			shortID, err := strconv.Atoi(chi.URLParam(r, "id"))
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

package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"net/http"
	"snach/go-musthave-shortener/cmd/shortener/repository"
	"strconv"
)

func NewRouter(repo repository.Repositorier) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/", CreateShortURLHandler(repo))
	r.Post("/api/shorten", CreateShortURLJSONHandler(repo))
	r.Get("/{id}", GetFullURLHandler(repo))
	return r
}

func CreateShortURLHandler(repo repository.Repositorier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		index, err := repo.Save(string(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "http://localhost:8080/%d", index)
	}
}

type RequestCreateShortURLJSON struct {
	URL string `json:"url"`
}

type ResponseCreateShortURLJSON struct {
	Result string `json:"result"`
}

func CreateShortURLJSONHandler(repo repository.Repositorier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			http.Error(w, "Bad Content-Type header, need application/json", http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		request := RequestCreateShortURLJSON{}
		if err := json.Unmarshal(body, &request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		index, err := repo.Save(request.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := ResponseCreateShortURLJSON{Result: "http://localhost:8080/" + strconv.Itoa(index)}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(responseJSON)
	}
}

func GetFullURLHandler(repo repository.Repositorier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fullURL, err := repo.Get(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", fullURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	}
}

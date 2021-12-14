package repository

import (
	"errors"
	"strconv"
)

//go:generate mockery --name=Repositorier --structname=RepositorierMock
type Repositorier interface {
	Get(shortURLID string) (string, error)
	Save(url string) (int, error)
}

type Repository struct {
	Storage    map[int]string
	CurrentInd int
}

func (r *Repository) Get(shortURLID string) (string, error) {
	shortID, err := strconv.Atoi(shortURLID)
	if err != nil {
		return "", err
	}

	if fullURL, ok := r.Storage[shortID]; ok {
		return fullURL, nil
	} else {
		return "", errors.New("No full url for short url index " + shortURLID)
	}

}

func (r *Repository) Save(url string) (int, error) {
	r.CurrentInd++
	r.Storage[r.CurrentInd] = url
	return r.CurrentInd, nil
}

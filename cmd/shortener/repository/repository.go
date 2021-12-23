package repository

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
)

//go:generate mockery --name=Repositorier --structname=RepositorierMock
type Repositorier interface {
	Get(shortURLID string) (string, error)
	Save(url string) (int, error)
}

type ShortToFullURL struct {
	Index   int
	FullURL string
}

type Repository struct {
	Storage    map[int]string
	CurrentInd int
	FileName   string
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

	file, err := os.OpenFile(r.FileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	itemJSON, err := json.Marshal(ShortToFullURL{Index: r.CurrentInd, FullURL: url})
	if err != nil {
		return 0, err
	}

	_, err = file.Write(append(itemJSON, '\n'))
	if err != nil {
		return 0, err
	}

	return r.CurrentInd, nil
}

func NewRepository(filename string) (*Repository, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	storageMap := make(map[int]string)
	maxIndex := 0

	reader := bufio.NewReader(file)

	for {
		lineBytes, err := reader.ReadBytes('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		itemMap := ShortToFullURL{}
		err = json.Unmarshal(lineBytes, &itemMap)
		if err != nil {
			return nil, err
		}

		storageMap[itemMap.Index] = itemMap.FullURL

		if itemMap.Index > maxIndex {
			maxIndex = itemMap.Index
		}
	}
	return &Repository{Storage: storageMap, CurrentInd: maxIndex, FileName: filename}, nil
}

package repository

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"testing"
)

const JSONStringForStorage = "{\"Index\":1,\"FullURL\":\"https://stackoverflow.com/\"}\n" +
	"{\"Index\":2,\"FullURL\":\"https://stepik.org/\"}\n" +
	"{\"Index\":3,\"FullURL\":\"https://hh.ru/\"}\n"

func TestNewRepositoryEmptyStorage(t *testing.T) {
	filename := "TestNewRepositoryFile.txt"
	repo, err := NewRepository(filename)
	defer os.Remove(filename)
	assert.NoError(t, err)
	assert.Equal(t, 0, repo.CurrentInd)
	assert.Equal(t, filename, repo.FileName)

	assert.Len(t, repo.Storage, 0)
}

func TestNewRepositoryNotEmptyStorage(t *testing.T) {
	filename := "TestNewRepositoryFile.txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	assert.NoError(t, err)
	file.Write([]byte(JSONStringForStorage))
	file.Close()
	defer os.Remove(filename)

	repo, err := NewRepository(filename)
	assert.NoError(t, err)
	assert.Equal(t, 3, repo.CurrentInd)
	assert.Equal(t, filename, repo.FileName)

	needStorageMap := map[int]string{
		1: "https://stackoverflow.com/",
		2: "https://stepik.org/",
		3: "https://hh.ru/",
	}
	assert.Len(t, repo.Storage, 3)
	assert.Equal(t, needStorageMap[1], repo.Storage[1])
	assert.Equal(t, needStorageMap[2], repo.Storage[2])
	assert.Equal(t, needStorageMap[3], repo.Storage[3])
}

func TestGetRepository(t *testing.T) {
	tests := []struct {
		name       string
		storage    map[int]string
		shortURLID string
		fullURL    string
		isErr      bool
	}{
		{
			name:       "positive test: exist url in storage",
			storage:    map[int]string{1: "https://stepik.org/"},
			shortURLID: "1",
			fullURL:    "https://stepik.org/",
			isErr:      false,
		},
		{
			name:       "negative test: bad index in url",
			storage:    map[int]string{1: "https://stepik.org/"},
			shortURLID: "abc",
			fullURL:    "",
			isErr:      true,
		},
		{
			name:       "negative test: no url in storage",
			storage:    map[int]string{1: "https://stepik.org/"},
			shortURLID: "2",
			fullURL:    "",
			isErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repository{Storage: tt.storage, CurrentInd: 0}
			fullURL, err := repo.Get(tt.shortURLID)
			assert.Equal(t, tt.fullURL, fullURL)
			if tt.isErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}

func TestSaveRepository(t *testing.T) {
	filename := "TestNewRepositoryFile.txt"
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	assert.NoError(t, err)
	file.Write([]byte(JSONStringForStorage))
	assert.NoError(t, err)
	file.Close()
	defer os.Remove(filename)

	repo, err := NewRepository(filename)
	assert.NoError(t, err)
	repo.Save("https://meduza.io/")

	assert.Equal(t, 4, repo.CurrentInd)
	assert.Equal(t, "https://meduza.io/", repo.Storage[4])

	file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDONLY, 0777)
	assert.NoError(t, err)

	bytes, err := io.ReadAll(file)
	assert.NoError(t, err)
	strs := strings.Split(string(bytes), "\n")
	assert.Equal(t, "{\"Index\":4,\"FullURL\":\"https://meduza.io/\"}", strs[3])
}

package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

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

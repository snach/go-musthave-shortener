package repository

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRepository(t *testing.T) {
	tests := []struct {
		name       string
		storage    map[int]string
		shortUrlId string
		fullUrl    string
		isErr      bool
	}{
		{
			name:       "positive test: exist url in storage",
			storage:    map[int]string{1: "https://stepik.org/"},
			shortUrlId: "1",
			fullUrl:    "https://stepik.org/",
			isErr:      false,
		},
		{
			name:       "negative test: bad index in url",
			storage:    map[int]string{1: "https://stepik.org/"},
			shortUrlId: "abc",
			fullUrl:    "",
			isErr:      true,
		},
		{
			name:       "negative test: no url in storage",
			storage:    map[int]string{1: "https://stepik.org/"},
			shortUrlId: "2",
			fullUrl:    "",
			isErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Repository{Storage: tt.storage, CurrentInd: 0}
			fullUrl, err := repo.Get(tt.shortUrlId)
			assert.Equal(t, tt.fullUrl, fullUrl)
			if tt.isErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}

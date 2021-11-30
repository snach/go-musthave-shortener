package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMethodShortenerHandler(t *testing.T) {
	type want struct {
		statusCode int
		location   string
		bodyLen    int
	}
	tests := []struct {
		name       string
		url        string
		urlMap     map[int]string
		mapCounter int
		want       want
	}{
		{
			name:       "positive test: exist url in urlMap",
			url:        "/1",
			urlMap:     map[int]string{1: "https://stepik.org/"},
			mapCounter: 2,
			want: want{
				statusCode: 307,
				location:   "https://stepik.org/",
				bodyLen:    0,
			},
		},
		{
			name:       "negative test: no url in urlMap",
			url:        "/100",
			urlMap:     map[int]string{1: "https://stepik.org/"},
			mapCounter: 2,
			want: want{
				statusCode: 400,
				location:   "",
				bodyLen:    0,
			},
		},
		{
			name:       "negative test: bad index in url",
			url:        "/abc",
			urlMap:     map[int]string{1: "https://stepik.org/"},
			mapCounter: 2,
			want: want{
				statusCode: 500,
				location:   "",
				bodyLen:    44,
			},
		},
		{
			name:       "negative test: no index in url",
			url:        "/",
			urlMap:     map[int]string{1: "https://stepik.org/"},
			mapCounter: 2,
			want: want{
				statusCode: 500,
				location:   "",
				bodyLen:    41,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ShortenerHandler(tt.urlMap, tt.mapCounter))
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)

			resBody, err := ioutil.ReadAll(res.Body)
			require.NoError(t, err)
			err = res.Body.Close()
			require.NoError(t, err)
			assert.Len(t, resBody, tt.want.bodyLen)

			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
		})

	}
}

func TestPostMethodShortenerHandler(t *testing.T) {

}

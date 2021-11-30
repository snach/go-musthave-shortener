package handlers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestPostMethodShortenerHandler(t *testing.T) {
	tests := []struct {
		name       string
		bodyReader io.Reader
		statusCode int
	}{
		{
			name:       "positive test",
			bodyReader: strings.NewReader("https://stackoverflow.com/"),
			statusCode: 201,
		},
		{
			name:       "negotive test: error in body reader",
			bodyReader: errReader(0),
			statusCode: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", tt.bodyReader)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ShortenerHandler(map[int]string{1: "https://stepik.org/"}, 2))
			h.ServeHTTP(w, request)
			res := w.Result()
			assert.Equal(t, tt.statusCode, res.StatusCode)
			if res.StatusCode == 201 {
				resBody, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)
				err = res.Body.Close()
				require.NoError(t, err)
				assert.Equal(t, "http://localhost:8080/2", string(resBody))
			}
		})
	}
}

func TestUnsupportedMethodShortenerHandler(t *testing.T) {
	request := httptest.NewRequest(http.MethodHead, "/", nil)
	w := httptest.NewRecorder()
	h := http.HandlerFunc(ShortenerHandler(map[int]string{1: "https://stepik.org/"}, 2))
	h.ServeHTTP(w, request)
	res := w.Result()
	assert.Equal(t, 400, res.StatusCode)
}

func TestIntegrationMapCounterIncrementShortenerHandler(t *testing.T) {
	h := http.HandlerFunc(ShortenerHandler(make(map[int]string), 1))

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://stackoverflow.com/")))
	assert.Equal(t, 201, w.Result().StatusCode)

	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://stepik.org/")))
	assert.Equal(t, 201, w.Result().StatusCode)

	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://hh.ru/")))
	assert.Equal(t, 201, w.Result().StatusCode)

	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/3", nil))
	res := w.Result()

	assert.Equal(t, 307, res.StatusCode)
	assert.Equal(t, "https://hh.ru/", res.Header.Get("Location"))

}

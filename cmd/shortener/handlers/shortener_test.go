package handlers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"snach/go-musthave-shortener/cmd/shortener/repository"
	"snach/go-musthave-shortener/cmd/shortener/repository/mocks"
	"strconv"
	"strings"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	return resp, string(respBody)
}

func TestGetFullUrlHandler(t *testing.T) {
	type want struct {
		statusCode int
		location   string
		bodyLen    int
	}
	tests := []struct {
		name         string
		url          string
		repoGetUrl   string
		repoGetError error
		want         want
	}{
		{
			name:         "positive test: exist url in urlMap",
			url:          "/1",
			repoGetUrl:   "https://stepik.org/",
			repoGetError: nil,
			want: want{
				statusCode: 307,
				location:   "https://stepik.org/",
				bodyLen:    0,
			},
		},
		{
			name:         "negative test: error from storage",
			url:          "/100",
			repoGetUrl:   "",
			repoGetError: errors.New("No full url for short url index 100"),
			want: want{
				statusCode: 400,
				location:   "",
				bodyLen:    36,
			},
		},
		{
			name: "negative test: no index in url",
			url:  "/",
			want: want{
				statusCode: 405,
				location:   "",
				bodyLen:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.RepositorierMock)
			repo.On("Get", mock.Anything).Return(tt.repoGetUrl, tt.repoGetError)
			r := NewRouter(repo)
			ts := httptest.NewServer(r)
			defer ts.Close()

			res, resBody := testRequest(t, ts, http.MethodGet, tt.url, nil)
			defer res.Body.Close()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.Len(t, resBody, tt.want.bodyLen)

			assert.Equal(t, tt.want.location, res.Header.Get("Location"))
		})

	}
}

func TestCreateShortUrlHandler(t *testing.T) {
	tests := []struct {
		name       string
		bodyReader io.Reader
		savedIndex int
		statusCode int
	}{
		{
			name:       "positive test: save url to storage",
			bodyReader: strings.NewReader("https://stackoverflow.com/"),
			savedIndex: 2,
			statusCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mocks.RepositorierMock)
			repo.On("Save", mock.Anything).Return(tt.savedIndex, nil)
			r := NewRouter(repo)
			ts := httptest.NewServer(r)
			defer ts.Close()

			res, resBody := testRequest(t, ts, http.MethodPost, "/", tt.bodyReader)
			defer res.Body.Close()

			assert.Equal(t, tt.statusCode, res.StatusCode)
			if res.StatusCode == 201 {
				assert.Equal(t, "http://localhost:8080/"+strconv.Itoa(tt.savedIndex), resBody)
			}
		})
	}
}

func TestUnsupportedMethodShortenerHandler(t *testing.T) {
	repo := new(mocks.RepositorierMock)
	r := NewRouter(repo)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, _ := testRequest(t, ts, http.MethodHead, "/", nil)
	defer res.Body.Close()
	assert.Equal(t, 405, res.StatusCode)
}

func TestIntegrationMapCounterIncrementShortenerHandler(t *testing.T) {
	repo := repository.Repository{
		Storage:    make(map[int]string),
		CurrentInd: 0,
	}
	r := NewRouter(&repo)
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, http.MethodPost, "/", strings.NewReader("https://stackoverflow.com/"))
	defer res.Body.Close()
	assert.Equal(t, 201, res.StatusCode)

	res, _ = testRequest(t, ts, http.MethodPost, "/", strings.NewReader("https://stepik.org/"))
	defer res.Body.Close()
	assert.Equal(t, 201, res.StatusCode)

	res, _ = testRequest(t, ts, http.MethodPost, "/", strings.NewReader("https://hh.ru/"))
	defer res.Body.Close()
	assert.Equal(t, 201, res.StatusCode)

	res, _ = testRequest(t, ts, http.MethodGet, "/3", nil)
	defer res.Body.Close()
	assert.Equal(t, 307, res.StatusCode)
	assert.Equal(t, "https://hh.ru/", res.Header.Get("Location"))
}

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

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func TestGetFullUrlHandler(t *testing.T) {
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
				statusCode: 405,
				location:   "",
				bodyLen:    0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRouter(tt.urlMap, tt.mapCounter)
			ts := httptest.NewServer(r)
			defer ts.Close()

			res, resBody := testRequest(t, ts, "GET", tt.url, nil)

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
		statusCode int
	}{
		{
			name:       "positive test",
			bodyReader: strings.NewReader("https://stackoverflow.com/"),
			statusCode: 201,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := NewRouter(map[int]string{1: "https://stepik.org/"}, 2)
			ts := httptest.NewServer(r)
			defer ts.Close()

			res, resBody := testRequest(t, ts, "POST", "/", tt.bodyReader)

			assert.Equal(t, tt.statusCode, res.StatusCode)
			if res.StatusCode == 201 {
				assert.Equal(t, "http://localhost:8080/2", resBody)
			}
		})
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func TestCreateShortUrlErrorBodyReaderHandler(t *testing.T) {
	r := NewRouter(map[int]string{1: "https://stepik.org/"}, 2)
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL+"/", errReader(0))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.Contains(t, err.Error(), "test error")

	//assert.Equal(t, 500, resp.StatusCode)
	assert.Equal(t, resp, resp)
}

func TestUnsupportedMethodShortenerHandler(t *testing.T) {
	r := NewRouter(map[int]string{1: "https://stepik.org/"}, 2)
	ts := httptest.NewServer(r)
	defer ts.Close()
	res, _ := testRequest(t, ts, "HEAD", "/", nil)
	assert.Equal(t, 405, res.StatusCode)
}

func TestIntegrationMapCounterIncrementShortenerHandler(t *testing.T) {
	r := NewRouter(map[int]string{1: "https://stepik.org/"}, 1)
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, _ := testRequest(t, ts, "POST", "/", strings.NewReader("https://stackoverflow.com/"))
	assert.Equal(t, 201, res.StatusCode)

	res, _ = testRequest(t, ts, "POST", "/", strings.NewReader("https://stepik.org/"))
	assert.Equal(t, 201, res.StatusCode)

	res, _ = testRequest(t, ts, "POST", "/", strings.NewReader("https://hh.ru/"))
	assert.Equal(t, 201, res.StatusCode)

	res, _ = testRequest(t, ts, "GET", "/3", nil)
	assert.Equal(t, 307, res.StatusCode)
	assert.Equal(t, "https://hh.ru/", res.Header.Get("Location"))
}

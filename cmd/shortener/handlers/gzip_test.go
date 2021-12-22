package handlers

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"snach/go-musthave-shortener/cmd/shortener/repository/mocks"
	"strings"
	"testing"
)

func TestGzipInRequestHandler(t *testing.T) {
	repo := new(mocks.RepositorierMock)
	repo.On("Save", mock.Anything).Return(1, nil)
	r := NewRouter(testBaseURL, repo)
	ts := httptest.NewServer(r)
	defer ts.Close()

	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	gzipWriter.Write([]byte("https://stackoverflow.com/"))
	gzipWriter.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/", &buf)
	require.NoError(t, err)
	req.Header.Set("Content-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	resp.Body.Close()
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Encoding"), "")
}

func TestGzipInResponseHandler(t *testing.T) {
	repo := new(mocks.RepositorierMock)
	repo.On("Save", mock.Anything).Return(1, nil)
	r := NewRouter(testBaseURL, repo)
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, err := http.NewRequest(http.MethodPost, ts.URL+"/", strings.NewReader("https://stackoverflow.com/"))
	require.NoError(t, err)
	req.Header.Set("Accept-Encoding", "gzip")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	gzipReader, err := gzip.NewReader(resp.Body)
	assert.NoError(t, err)
	respBody, err := ioutil.ReadAll(gzipReader)
	assert.NoError(t, err)
	resp.Body.Close()

	assert.Equal(t, testBaseURL+"/1", string(respBody))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, resp.Header.Get("Content-Encoding"), "gzip")
}

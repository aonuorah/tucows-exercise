package service

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	Client = clientMockOk()

	table := []struct {
		path     string
		expected string
	}{
		{"/", "testdata/test-1.html"},
		{"/?category=100", "testdata/test-2.html"},
		{"/?category=100&grayscale=true", "testdata/test-3.html"},
	}

	for _, row := range table {
		req, err := http.NewRequest("GET", row.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(rootHandler)

		handler.ServeHTTP(rr, req)
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check the content type is what we expect.
		want := "text/html; charset=utf-8"
		if got := rr.Header().Get("Content-Type"); got != want {
			t.Errorf("handler returned wrong content type: got %v want %v",
				got, want)
		}

		// Check the response body is what we expect.
		expected, _ := os.ReadFile(row.expected)
		if rr.Body.String() != string(expected) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), string(expected))
		}
	}
}

type clientMock struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *clientMock) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func clientMockOk() *clientMock {
	return &clientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if req.Host == "picsum.photos" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("raw photo bytes")),
				}, nil
			}

			if req.Host == "api.forismatic.com" {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("{\"quoteText\":\"test quote\"}")),
				}, nil
			}

			return &http.Response{}, nil
		},
	}
}

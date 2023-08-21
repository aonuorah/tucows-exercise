package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/aonuorah/tucows-exercise/internal/service"
	"github.com/spf13/afero"
)

const (
	defaultPhotofile = "tmp/image.jpg"
	defaultQuotefile = "tmp/quote.txt"
)

func TestInvalidArgs(t *testing.T) {
	if os.Getenv("CHECK_ARGS") == "1" {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		service.Client = clientMockOk()
		service.FS = afero.NewMemMapFs()

		args := strings.Split(os.Getenv("INVALID_ARGS"), " ")
		os.Args = append([]string{"main"}, args...)
		main()
		return
	}

	table := []struct {
		args string
	}{
		{"--output-quote"},
		{"--output-image"},
		{"--category invalid"},
	}

	for _, row := range table {
		cmd := exec.Command(os.Args[0], "-test.run=TestInvalidArgs")
		cmd.Env = append(os.Environ(), "CHECK_ARGS=1", "INVALID_ARGS="+row.args)

		if err := cmd.Run(); err == nil {
			t.Fatalf("invalid argument '%v' expected to fail, but completed successfully", row.args)
		}
	}
}

func TestValidArgsWithAvailableApi(t *testing.T) {
	table := []struct {
		args      string
		imagefile string
		quotefile string
	}{
		{"", defaultPhotofile, defaultQuotefile},
		{"--grayscale", defaultPhotofile, defaultQuotefile},
		{"--category 100", defaultPhotofile, defaultQuotefile},
		{"--output-quote test-1.txt", defaultPhotofile, "test-1.txt"},
		{"--output-image test-2.jpg", "test-2.jpg", defaultQuotefile},
		{"--output-quote test-3.txt --output-image test-3.jpg", "test-3.jpg", "test-3.txt"},
	}

	service.Client = clientMockOk()
	service.FS = afero.NewMemMapFs()

	for _, row := range table {
		args := strings.Split(row.args, " ")
		os.Args = append([]string{"main"}, args...)

		main()

		gotPhoto, _ := afero.ReadFile(service.FS, row.imagefile)
		expected := "raw photo bytes"
		if string(gotPhoto) != expected {
			t.Fatalf("process with args '%v' failed! expected photo '%v' but got '%v'", row.args, expected, string(gotPhoto))
		}

		gotQuote, _ := afero.ReadFile(service.FS, row.quotefile)
		expected = "test quote"
		if string(gotQuote) != "test quote" {
			t.Fatalf("process with args '%v' failed! expected quote '%v' but got '%v'", row.args, expected, string(gotQuote))
		}
	}
}

func TestValidArgsWithFailingApi(t *testing.T) {
	table := []struct {
		args      string
		imagefile string
		quotefile string
	}{
		{"", defaultPhotofile, defaultQuotefile},
		{"--output-quote test-1.txt", defaultPhotofile, "test-1.txt"},
		{"--output-image test-2.jpg", "test-2.jpg", defaultQuotefile},
		{"--output-quote test-3.txt --output-image=test-3.jpg", "test-3.jpg", "test-3.txt"},
	}

	service.Client = clientMockFailing()
	service.FS = afero.NewMemMapFs()

	// prepare default photo and quote in FS memory
	expectedPhoto := "default photo raw bytes"
	service.WriteToFile([]byte(expectedPhoto), "web/img/img-1-gs.jpg")

	expectedQuote := "default quote"
	service.WriteToFile([]byte(expectedQuote), "web/data/quote-1.txt")

	for _, row := range table {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		args := strings.Split(row.args, " ")
		os.Args = append([]string{"main"}, args...)

		main()

		gotPhoto, _ := afero.ReadFile(service.FS, row.imagefile)
		if string(gotPhoto) != expectedPhoto {
			t.Fatalf("process with args '%v' failed! expected quote '%v' but got '%v'", row.args, expectedPhoto, string(gotPhoto))
		}

		gotQuote, _ := afero.ReadFile(service.FS, row.quotefile)
		if string(gotQuote) != expectedQuote {
			t.Fatalf("process with args '%v' failed! expected quote '%v' but got '%v'", row.args, expectedQuote, string(gotQuote))
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

func clientMockFailing() *clientMock {
	return &clientMock{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}
}

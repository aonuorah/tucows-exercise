package service

import (
	"os"
	"testing"

	"github.com/aonuorah/tucows-exercise/internal/photo"
	"github.com/aonuorah/tucows-exercise/internal/quote"
)

func TestBuildHTML(t *testing.T) {

	table := []struct {
		category  int
		quote     string
		grayscale bool
		bytes     string
		base64Img string
		width     int
		height    int
		htmlfile  string
	}{
		{
			category:  100,
			quote:     "test quote",
			grayscale: false,
			bytes:     "raw photo bytes",
			base64Img: "cmF3IHBob3RvIGJ5dGVz",
			width:     600,
			height:    400,
			htmlfile:  "testdata/test-2.html",
		},
		{
			category:  10000,
			quote:     "test quote",
			grayscale: true,
			bytes:     "test",
			base64Img: "dGVzdA==",
			width:     1000,
			height:    600,
			htmlfile:  "testdata/test-4.html",
		},
	}

	for _, row := range table {
		p := &photo.Info{
			Bytes:     []byte(row.bytes),
			Grayscale: row.grayscale,
			Width:     row.width,
			Height:    row.height,
		}
		q := &quote.Info{
			Text:     row.quote,
			Category: row.category,
		}

		got, _ := BuildHTML(p, q)
		want, _ := os.ReadFile(row.htmlfile)

		if string(got) != string(want) {
			t.Errorf("mismatch: got %v want %v",
				string(got), string(want))
		}
	}
}

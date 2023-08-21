package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aonuorah/tucows-exercise/internal/photo"
	"github.com/aonuorah/tucows-exercise/internal/quote"
	"github.com/spf13/afero"
)

var (
	// The default http.Client used by the photo and quote providers
	// client has a 10s timeout to prevent long running http calls
	Client httpclient = &http.Client{Timeout: 10 * time.Second}
)

type httpclient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Gets a new photo and quote from their respective api providers
// Note: this method returns a default photo or quote if the respective api provider returns an error.
func GetNewItem(ctx context.Context) (*photo.Info, *quote.Info) {
	pCh := make(chan func() (*photo.Info, error))
	qCh := make(chan func() (*quote.Info, error))

	go func(ctx context.Context) {
		info, err := photo.Get(ctx, Client)
		pCh <- func() (*photo.Info, error) { return info, err }
	}(ctx)

	go func(ctx context.Context) {
		info, err := quote.Get(ctx, Client)
		qCh <- func() (*quote.Info, error) { return info, err }
	}(ctx)

	photo, err := (<-pCh)()
	if err != nil {
		fmt.Println(err)
		photo = defaultPhoto()
	}

	quote, err := (<-qCh)()
	if err != nil {
		fmt.Println(err)
		quote = defaultQuote()
	}

	return photo, quote
}

// default photo loaded in memory
var defPhoto *photo.Info

// gets the default photo used when the photo provider fails
// this method loads the defPhoto variable when access for the first time
func defaultPhoto() *photo.Info {
	if defPhoto == nil {
		bytes, _ := afero.ReadFile(FS, "web/img/img-1-gs.jpg")
		defPhoto = &photo.Info{
			Bytes:     bytes,
			Grayscale: true,
			Width:     photo.DefaultImageWidth,
			Height:    photo.DefaultImageHeight,
		}
	}

	return defPhoto
}

// default quote loaded in memory
var defQuote *quote.Info

// gets the default quote used when the quote provider fails
// this method loads the defQuote variable when access for the first time
func defaultQuote() *quote.Info {
	if defQuote == nil {
		bytes, _ := afero.ReadFile(FS, "web/data/quote-1.txt")
		return &quote.Info{
			Text:     string(bytes),
			Category: quote.DefaultCategory,
		}
	}

	return defQuote
}

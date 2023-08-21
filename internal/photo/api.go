package photo

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// The default image width used in the photo provider request
	DefaultImageWidth = 600

	// The default image height used in the photo provider request
	DefaultImageHeight = 400
)

var (
	// base url of the photo provider url
	// url sample - https://picsum.photos/{width}/{height}?grayscale
	photoApiBaseUrl = fmt.Sprintf("%s/%d/%d", "https://picsum.photos", DefaultImageWidth, DefaultImageHeight)
)

type httpclient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Makes an HTTP request to the photo provider url and returns an object with properties of the response photo.
func Get(ctx context.Context, client httpclient) (*Info, error) {
	base, err := url.Parse(photoApiBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("error building photo url: %v", err)
	}

	// Build query params
	grayscale := GrayscaleFromCtx(ctx)
	if grayscale {
		params := url.Values{}
		params.Add("grayscale", "true")
		base.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating photo request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("photo request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("photo request failed with status: %v", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading photo request body: %v", err)
	}

	return &Info{
		Bytes:     bytes,
		Grayscale: grayscale,
		Width:     DefaultImageWidth,
		Height:    DefaultImageHeight,
	}, nil
}

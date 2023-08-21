package quote

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	// The default category used in the quote provider request
	DefaultCategory = 100000

	// base url of the quote provider url
	quoteApiBaseUrl = "https://api.forismatic.com/api/1.0/"
)

type httpclient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Makes an HTTP request to the quote provider url and returns an object with properties of the response quote.
func Get(ctx context.Context, client httpclient) (*Info, error) {
	base, err := url.Parse(quoteApiBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("error building quote url: %v", err)
	}

	// Build query params
	category := CategoryFromCtx(ctx)
	params := url.Values{}
	params.Add("method", "getQuote")
	params.Add("format", "json")
	params.Add("key", strconv.Itoa(category))
	params.Add("lang", "en")
	base.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating quote request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("quote request failed: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("quote request failed with status: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading quote request body: %v", err)
	}

	body = []byte(strings.ReplaceAll(string(body), `\`, ""))

	var q Info
	if err := json.Unmarshal(body, &q); err != nil {
		return nil, fmt.Errorf("error unmarshalling quote body: %v", err)
	}

	q.Category = category

	return &q, nil
}

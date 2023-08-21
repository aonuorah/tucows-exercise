package photo

import (
	"context"
	"encoding/base64"
)

// type used to embed values into contexts
type ctxkey struct {
	name string
}

var (
	// key name for 'grayscale' variable set in context.Context
	ctxKeyGrayscale = ctxkey{"grayscale"}
)

// Returns a new context with the 'grayscale' variable added
func CtxWithGrayscale(ctx context.Context, val bool) context.Context {
	return context.WithValue(ctx, ctxKeyGrayscale, val)
}

// Gets the grayscale value from the specified context
func GrayscaleFromCtx(ctx context.Context) bool {
	val, _ := ctx.Value(ctxKeyGrayscale).(bool)
	return val
}

// struct to hold properties of the photo returned from the photo provider
type Info struct {
	Bytes     []byte
	Grayscale bool
	Width     int
	Height    int
	Filename  string
}

// Returns a base64 encoded string of the photo
func (i Info) Base64Encoding() string {
	return base64.StdEncoding.EncodeToString(i.Bytes)
}

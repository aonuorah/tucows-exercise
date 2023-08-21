package quote

import (
	"context"
)

// type used to embed values into contexts
type ctxkey struct {
	name string
}

var (
	// key name for 'category' variable set in context.Context
	ctxKeyCategory = ctxkey{"category"}
)

// Returns a new context with the 'category' variable added.
// Note: if an invalid value is specified, then the default is set instead.
func CtxWithCategory(ctx context.Context, val int) context.Context {
	if val < 1 || val > 999999 {
		val = DefaultCategory
	}

	return context.WithValue(ctx, ctxKeyCategory, val)
}

// Gets the category value from the specified context
func CategoryFromCtx(ctx context.Context) int {
	val, _ := ctx.Value(ctxKeyCategory).(int)
	return val
}

// struct to hold properties of the quote returned from the quote provider
type Info struct {
	Text     string `json:"quoteText,omitempty"`
	Category int    `json:"-"`
}

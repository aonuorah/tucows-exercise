package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aonuorah/tucows-exercise/internal/photo"
	"github.com/aonuorah/tucows-exercise/internal/quote"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Registers the service route 'patterns' and their corresponding handlers
func RegisterHandlers(mux *chi.Mux) {
	mux.Use(middleware.Timeout(30 * time.Second))

	mux.Get("/", rootHandler)
}

// the root handler to serve the main html page
func rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := ctxWithParams(r)

	photo, quote := GetNewItem(ctx)
	body, err := BuildHTML(photo, quote)
	if err != nil {
		fmt.Printf("error building html from template: %v", err)
		http.Error(w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(body)
}

// Gets a new context with supported query params added
func ctxWithParams(r *http.Request) context.Context {
	grayscale, _ := strconv.ParseBool(r.URL.Query().Get("grayscale"))
	category, _ := strconv.Atoi(r.URL.Query().Get("category"))

	ctx := photo.CtxWithGrayscale(r.Context(), grayscale)
	ctx = quote.CtxWithCategory(ctx, category)
	return ctx
}

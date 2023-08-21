package main

import (
	"net/http"

	"github.com/aonuorah/tucows-exercise/internal/service"
	"github.com/go-chi/chi/v5"
)

func main() {
	mux := chi.NewRouter()
	service.RegisterHandlers(mux)
	http.ListenAndServe(":8080", mux)
}

package controller

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/api/address/search", h.AddressSearchHandler)
	r.Post("/api/address/geocode", h.AddressGeocodeHandler)

	return r
}

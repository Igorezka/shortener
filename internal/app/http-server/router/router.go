package router

import (
	"github.com/Igorezka/shortener/internal/app/config"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/create"
	"github.com/Igorezka/shortener/internal/app/http-server/handlers/url/get"
	"github.com/Igorezka/shortener/internal/app/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New(cfg *config.Config, store *storage.Store) chi.Router {
	r := chi.NewRouter()

	r.Get("/{id}", get.New(store))
	r.Group(func(r chi.Router) {
		r.Use(middleware.AllowContentType("text/plain"))
		r.Post("/", create.New(cfg, store))
	})

	return r
}
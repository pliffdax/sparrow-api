package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pliffdax/sparrow-api/internal/http/handlers"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	r.Get("/health", handlers.HealthCheck())

	r.Route("/user", func(r chi.Router) {
		r.Post("/", handlers.CreateUser())
		r.Get("/{user_id}", handlers.GetUser())       // id with fetching logic will be added later
		r.Delete("/{user_id}", handlers.DeleteUser()) // id with deleting logic will be added later
	})
	r.Get("/users", handlers.ListUsers())

	r.Route("/category", func(r chi.Router) {
		r.Post("/", handlers.CreateUser())
		r.Delete("/{id}", handlers.DeleteUser()) // id with deleting logic will be added later
	})
	r.Get("/category", handlers.ListCategories())

	r.Route("/record", func(r chi.Router) {
		r.Post("/", handlers.CreateRecord())
		r.Get("/{record_id}", handlers.GetRecord())       // id with fetching logic will be added later
		r.Delete("/{record_id}", handlers.DeleteRecord()) // id with deleting logic will be added later
	})
	r.Get("/record", handlers.QueryRecords()) // ?user_id=&category_id=

	return r
}

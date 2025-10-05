package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/pliffdax/sparrow-api/internal/http/handlers"
	"github.com/pliffdax/sparrow-api/internal/storage"
)

func NewRouter(
	us storage.UserStore,
	cs storage.CategoryStore,
	rs storage.RecordStore,
) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", handlers.HealthCheck())

	r.Route("/users", func(r chi.Router) {
		r.Get("/", handlers.ListUsers(us))
		r.Get("/{id}", handlers.GetUser(us))
		r.Post("/", handlers.CreateUser(us))
		r.Delete("/{id}", handlers.DeleteUser(us))
	})

	r.Route("/categories", func(r chi.Router) {
		r.Get("/", handlers.ListCategories(cs))
		r.Post("/", handlers.CreateCategory(cs))
		r.Delete("/{id}", handlers.DeleteCategory(cs))
	})

	r.Route("/records", func(r chi.Router) {
		r.Get("/", handlers.QueryRecords(rs)) // ?user_id=&category_id=
		r.Get("/{id}", handlers.GetRecord(rs))
		r.Post("/", handlers.CreateRecord(rs, us, cs))
		r.Delete("/{id}", handlers.DeleteRecord(rs))
	})

	return r
}

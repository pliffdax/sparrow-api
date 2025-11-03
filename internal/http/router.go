package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/pliffdax/sparrow-api/internal/http/handlers"
	authmw "github.com/pliffdax/sparrow-api/internal/http/middleware"
	"github.com/pliffdax/sparrow-api/internal/storage"
)

func NewRouter(us storage.UserStore, cs storage.CategoryStore, rs storage.RecordStore) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", handlers.HealthCheck())

	aus, ok := us.(storage.AuthUserStore)
	if !ok {
		panic("UserStore does not implement AuthUserStore")
	}
	auth := handlers.AuthHandler{Users: aus}

	r.Post("/auth/register", auth.Register())
	r.Post("/auth/login", auth.Login())

	r.Group(func(r chi.Router) {
		r.Use(authmw.AuthRequired)

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
			r.Get("/", handlers.QueryRecords(rs))
			r.Get("/{id}", handlers.GetRecord(rs))
			r.Post("/", handlers.CreateRecord(rs, us, cs))
			r.Delete("/{id}", handlers.DeleteRecord(rs))
		})
	})

	return r
}

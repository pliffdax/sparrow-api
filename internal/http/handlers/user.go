package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pliffdax/sparrow-api/internal/storage"
	"github.com/pliffdax/sparrow-api/internal/util"
)

type createUserReq struct {
	Name string `json:"name"`
}

func CreateUser(us storage.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createUserReq
		if err := util.DecodeJSON(r, &req); err != nil && req.Name == "" {
			http.Error(w, "invalid body: need {\"name\":\"...\"}", http.StatusBadRequest)
			return
		}
		u, _ := us.Create(req.Name)
		util.WriteJSON(w, http.StatusCreated, u)
	}
}
func GetUser(us storage.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		u, ok := us.GetByID(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		util.WriteJSON(w, http.StatusOK, u)
	}
}
func DeleteUser(us storage.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if !us.Delete(id) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
func ListUsers(us storage.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.WriteJSON(w, http.StatusOK, us.List())
	}
}

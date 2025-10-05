package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pliffdax/sparrow-api/internal/storage"
	"github.com/pliffdax/sparrow-api/internal/util"
)

type createCategoryReq struct {
	Title string `json:"title"`
}

func CreateCategory(cs storage.CategoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createCategoryReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && req.Title == "" {
			http.Error(w, "invalid body: need {\"title\":\"...\"}", http.StatusBadRequest)
			return
		}
		cat, _ := cs.Create(req.Title)
		util.WriteJSON(w, http.StatusCreated, cat)
	}
}
func ListCategories(cs storage.CategoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		util.WriteJSON(w, http.StatusOK, cs.List())
	}
}
func DeleteCategory(cs storage.CategoryStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if !cs.Delete(id) {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

package handlers

import "net/http"

func CreateCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}
func ListCategories() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}
func DeleteCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}

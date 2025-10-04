package handlers

import "net/http"

func CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}
func GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}
func DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}
func ListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}

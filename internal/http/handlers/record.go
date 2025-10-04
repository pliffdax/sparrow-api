package handlers

import "net/http"

func CreateRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}
func GetRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}
func DeleteRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}
func QueryRecords() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotImplemented) }
}

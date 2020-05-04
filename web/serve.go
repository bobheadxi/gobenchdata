package web

import "net/http"

// ListenAndServe serves the web app on the given address
func ListenAndServe(addr string) error {
	handler := http.NewServeMux()
	handler.Handle("/", http.FileServer(HTTP))
	return http.ListenAndServe(addr, handler)
}

package handlers

import "net/http"

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	}
}

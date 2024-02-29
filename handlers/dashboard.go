package handlers

import "net/http"

func DashboardHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := w.Write([]byte("Dashboard"))
		if err != nil {
			return
		}
	}
}

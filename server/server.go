package server

import (
	"net/http"
)

func StartServer(urlContent string, downloadError error) {
	http.HandleFunc("/content", func(w http.ResponseWriter, r *http.Request) {
		if downloadError != nil {
			http.Error(w, downloadError.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(urlContent))
	})

	http.ListenAndServe(":8080", nil)
}

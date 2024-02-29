package main

import (
	"log"
	"net/http"
)

func middlewareCors(next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})

	return handler
}

func main() {
	mux := http.NewServeMux()
	corsMux := middlewareCors(mux)
	server := &http.Server{
		Addr:    ":8081",
		Handler: corsMux,
	}
	log.Fatal(server.ListenAndServe())
}

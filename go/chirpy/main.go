package main

import (
	"log"
	"net/http"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

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

	mux.Handle("/healthz", http.HandlerFunc(healthz))
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	mux.Handle("assets/logo.png", http.FileServer(http.Dir("./assets/logo.png")))

	server := &http.Server{
		Addr:    ":8081",
		Handler: corsMux,
	}

	log.Fatal(server.ListenAndServe())
}

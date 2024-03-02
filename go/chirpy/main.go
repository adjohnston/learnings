package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func metrics(c *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %d", c.fileserverHits)))
	}
}

func resetMetrics(c *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits = 0
		w.WriteHeader(http.StatusOK)
	}
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}

	d := json.NewDecoder(r.Body)
	p := params{}
	err := d.Decode(&p)

	if err != nil {
		type response struct {
			Error string `json:"error"`
		}

		resp, _ := json.Marshal(response{Error: "Something went wrong"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	if len(p.Body) > 140 {
		type response struct {
			Error string `json:"error"`
		}

		resp, _ := json.Marshal(response{Error: "Chirp is too long"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(resp)
		return
	}

	type response struct {
		Valid bool `json:"valid"`
	}

	resp, _ := json.Marshal(response{Valid: true})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func middlewareMetricsInc(cfg *apiConfig, next http.Handler) http.Handler {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})

	return handler
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
	r := chi.NewRouter()
	corsMux := middlewareCors(r)
	hits := apiConfig{fileserverHits: 0}

	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)
	apiRouter.Get("/healthz", healthz)
	apiRouter.Get("/metrics", metrics(&hits))
	apiRouter.Handle("/reset", http.HandlerFunc(resetMetrics(&hits)))
	apiRouter.Post("/validate_chirp", validateChirp)

	adminRouter := chi.NewRouter()
	r.Mount("/admin", adminRouter)
	r.Get("/admin/metrics", hits.metricsHandler)

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	r.Handle("/app", middlewareMetricsInc(&hits, fsHandler))
	r.Handle("/app/*", middlewareMetricsInc(&hits, fsHandler))

	server := &http.Server{
		Addr:    ":8081",
		Handler: corsMux,
	}

	log.Fatal(server.ListenAndServe())
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/adjohnston/chirpy/internal/db"
	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileserverHits int
	DB             *db.DB
}

func sanitiseChirp(original string) (chirp string) {
	bannedWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	words := strings.Split(original, " ")

	for i, w := range words {
		if bannedWords[strings.ToLower(w)] {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func metrics(c *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, fmt.Sprintf("Hits: %d", c.fileserverHits))
	}
}

func resetMetrics(c *apiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits = 0
		w.WriteHeader(http.StatusOK)
	}
}

func getChirps(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := db.GetChirps()

		if err != nil {
			type response struct {
				Error string `json:"error"`
			}

			respondWithJSON(w, http.StatusInternalServerError, response{Error: "Something went wrong"})
			return
		}

		respondWithJSON(w, http.StatusOK, chirps)
	}
}

func createChirps(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

			respondWithJSON(w, http.StatusBadRequest, response{Error: "Something went wrong"})
			return
		}

		if len(p.Body) > 140 {
			type response struct {
				Error string `json:"error"`
			}

			respondWithJSON(w, http.StatusBadRequest, response{Error: "Chirp is too long"})
			return
		}

		santisedChirp := sanitiseChirp(p.Body)
		newChirp, err := db.CreateChirp(santisedChirp)

		if err != nil {
			type response struct {
				Error string `json:"error"`
			}

			respondWithJSON(w, http.StatusInternalServerError, response{Error: "Something went wrong"})
			return
		}

		type response struct {
			Body string `json:"body"`
			ID   int    `json:"id"`
		}

		respondWithJSON(w, http.StatusCreated, response{ID: newChirp.ID, Body: newChirp.Body})
	}
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
	db, err := db.NewDB("database.json")

	if err != nil {
		log.Fatal(err)
	}

	hits := apiConfig{fileserverHits: 0, DB: db}

	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)
	apiRouter.Get("/healthz", healthz)
	apiRouter.Get("/metrics", metrics(&hits))
	apiRouter.Handle("/reset", http.HandlerFunc(resetMetrics(&hits)))
	apiRouter.Get("/chirps", getChirps(db))
	apiRouter.Post("/chirps", createChirps(db))

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

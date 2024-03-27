package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/adjohnston/chirpy/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
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

func respondWithErr(w http.ResponseWriter, code int, error string) {
	type response struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, response{Error: error})
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

func createUser(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		d := json.NewDecoder(r.Body)
		p := params{}
		decodeErr := d.Decode(&p)
		hashedPassword, hashedPasswordErr := db.HashPassword(p.Password)
		newUser, createUserErr := db.CreateUser(p.Email, hashedPassword)

		if decodeErr != nil || hashedPasswordErr != nil || createUserErr != nil {
			respondWithErr(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		type response struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		}

		respondWithJSON(w, http.StatusCreated, response{ID: newUser.ID, Email: newUser.Email})
	}
}

func updateUser(db *db.DB, c apiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		token := strings.Split(bearerToken, " ")[1]
		userId, err := c.Validate(token)

		if err != nil {
			respondWithErr(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		type params struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		d := json.NewDecoder(r.Body)
		p := params{}
		err = d.Decode(&p)
		hashedPassword, hashedPasswordErr := db.HashPassword(p.Password)

		if err != nil || hashedPasswordErr != nil {
			respondWithErr(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		updatedUser, err := db.UpdateUser(userId, p.Email, hashedPassword)

		if err != nil {
			respondWithErr(w, http.StatusInternalServerError, "Unable to update user")
			return
		}

		type response struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		}

		respondWithJSON(w, http.StatusOK, response{ID: updatedUser.ID, Email: updatedUser.Email})
	}
}

func login(db *db.DB, c apiConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type requestBody struct {
			Email            string `json:"email"`
			Password         string `json:"password"`
			ExpiresInSeconds int    `json:"expires_in_seconds"`
		}

		d := json.NewDecoder(r.Body)
		json := requestBody{}
		err := d.Decode(&json)

		if err != nil {
			respondWithErr(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		user, err := db.LoginUser(json.Email, json.Password)

		if err != nil {
			respondWithErr(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		token, err := c.CreateAuthToken(user.ID, json.ExpiresInSeconds)

		if err != nil {
			respondWithErr(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		type response struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
			Token string `json:"token"`
		}

		respondWithJSON(w, http.StatusOK, response{ID: user.ID, Email: user.Email, Token: token})
	}
}

func getChirps(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := db.GetChirps()

		if err != nil {
			respondWithErr(w, http.StatusInternalServerError, "Something went wrong")
			return
		}

		respondWithJSON(w, http.StatusOK, chirps)
	}
}

func getChirp(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idParam)

		if err != nil {
			respondWithErr(w, http.StatusBadRequest, "Unable to parse ID")
			return
		}

		chirp, err := db.GetChirpByID(id)

		if err != nil {
			respondWithErr(w, http.StatusNotFound, "Chirp not found")
			return
		}

		type response struct {
			ID   int    `json:"id"`
			Body string `json:"body"`
		}

		respondWithJSON(w, http.StatusOK, response{ID: chirp.ID, Body: chirp.Body})
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
			respondWithErr(w, http.StatusBadRequest, "Something went wrong")
			return
		}

		if len(p.Body) > 140 {
			respondWithErr(w, http.StatusBadRequest, "Chirp is too long")
			return
		}

		santisedChirp := sanitiseChirp(p.Body)
		newChirp, err := db.CreateChirp(santisedChirp)

		if err != nil {
			respondWithErr(w, http.StatusInternalServerError, "Something went wrong")
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
	godotenv.Load()
	r := chi.NewRouter()
	corsMux := middlewareCors(r)
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *debug {
		fmt.Println("Debug mode enabled")
		err := os.Remove("database.json")
		if err != nil {
			fmt.Println("Error removing database.json", err)
		}
	}

	db, err := db.NewDB("database.json")

	if err != nil {
		log.Fatal(err)
	}

	config := apiConfig{fileserverHits: 0, jwtSecret: os.Getenv("JWT_SECRET"), DB: db}

	apiRouter := chi.NewRouter()
	r.Mount("/api", apiRouter)
	apiRouter.Get("/healthz", healthz)
	apiRouter.Get("/metrics", metrics(&config))
	apiRouter.Handle("/reset", http.HandlerFunc(resetMetrics(&config)))
	apiRouter.Get("/chirps", getChirps(db))
	apiRouter.Get("/chirps/{id}", getChirp(db))
	apiRouter.Post("/chirps", createChirps(db))
	apiRouter.Post("/users", createUser(db))
	apiRouter.Put("/users", updateUser(db, config))
	apiRouter.Post("/login", login(db, config))

	adminRouter := chi.NewRouter()
	r.Mount("/admin", adminRouter)
	r.Get("/admin/metrics", config.metricsHandler)

	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	r.Handle("/app", middlewareMetricsInc(&config, fsHandler))
	r.Handle("/app/*", middlewareMetricsInc(&config, fsHandler))

	server := &http.Server{
		Addr:    ":8081",
		Handler: corsMux,
	}

	log.Fatal(server.ListenAndServe())
}

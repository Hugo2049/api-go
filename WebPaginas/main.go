package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux" // es la librería para manejar rutas y parámetros en la API
	"github.com/rs/cors"     //permite clientes en otros dominios
)

// match representa un partido de fútbol
type Match struct {
	ID        int    `json:"id"`
	HomeTeam  string `json:"homeTeam"`
	AwayTeam  string `json:"awayTeam"`
	MatchDate string `json:"matchDate"`
}

// MatchStore administra los partidos en memoria
type MatchStore struct {
	mu      sync.RWMutex
	matches map[int]Match
	nextID  int
}

// NewMatchStore crea un nuevo almacenamiento de partidos
func NewMatchStore() *MatchStore {
	return &MatchStore{
		matches: make(map[int]Match),
		nextID:  1,
	}
}

// createMatch agrega un nuevo partido
func (ms *MatchStore) CreateMatch(match Match) Match {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	match.ID = ms.nextID
	ms.matches[ms.nextID] = match
	ms.nextID++
	return match
}

// GetAllMatches devuelve todos los partidos
func (ms *MatchStore) GetAllMatches() []Match {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	matches := make([]Match, 0, len(ms.matches))
	for _, match := range ms.matches {
		matches = append(matches, match)
	}
	return matches
}

// GetMatchByID obtiene un partido por su ID
func (ms *MatchStore) GetMatchByID(id int) (Match, bool) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	match, exists := ms.matches[id]
	return match, exists
}

// UpdateMatch actualiza un partido existente
func (ms *MatchStore) UpdateMatch(id int, updatedMatch Match) (Match, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	if _, exists := ms.matches[id]; !exists {
		return Match{}, false
	}
	updatedMatch.ID = id
	ms.matches[id] = updatedMatch
	return updatedMatch, true
}

// DeleteMatch elimina un partido
func (ms *MatchStore) DeleteMatch(id int) bool {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	_, exists := ms.matches[id]
	if !exists {
		return false
	}
	delete(ms.matches, id)
	return true
}

// Manejadores de la API
func createMatchHandler(store *MatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var match Match
		if err := json.NewDecoder(r.Body).Decode(&match); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdMatch := store.CreateMatch(match)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdMatch)
	}
}

func getAllMatchesHandler(store *MatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matches := store.GetAllMatches()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(matches)
	}
}

func getMatchByIDHandler(store *MatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de partido inválido", http.StatusBadRequest)
			return
		}

		match, exists := store.GetMatchByID(id)
		if !exists {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
	}
}

func updateMatchHandler(store *MatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de partido inválido", http.StatusBadRequest)
			return
		}

		var updatedMatch Match
		if err := json.NewDecoder(r.Body).Decode(&updatedMatch); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		match, updated := store.UpdateMatch(id, updatedMatch)
		if !updated {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
	}
}

func deleteMatchHandler(store *MatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de partido inválido", http.StatusBadRequest)
			return
		}

		deleted := store.DeleteMatch(id)
		if !deleted {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	store := NewMatchStore()
	r := mux.NewRouter()

	// Middleware CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Rutas de partidos
	r.HandleFunc("/api/matches", createMatchHandler(store)).Methods("POST")
	r.HandleFunc("/api/matches", getAllMatchesHandler(store)).Methods("GET")
	r.HandleFunc("/api/matches/{id}", getMatchByIDHandler(store)).Methods("GET")
	r.HandleFunc("/api/matches/{id}", updateMatchHandler(store)).Methods("PUT")
	r.HandleFunc("/api/matches/{id}", deleteMatchHandler(store)).Methods("DELETE")

	// Aplicar middleware CORS
	handler := c.Handler(r)

	//iniciar servidor en el puerto 8080
	puerto := 8080
	fmt.Printf("Servidor iniciando en el puerto %d...\n", puerto)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", puerto), handler))
}

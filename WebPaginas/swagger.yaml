package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Match representa un partido de fútbol
type Match struct {
	ID          int    `json:"id"`
	HomeTeam    string `json:"homeTeam"`
	AwayTeam    string `json:"awayTeam"`
	MatchDate   string `json:"matchDate"`
	HomeGoals   int    `json:"homeGoals,omitempty"`
	AwayGoals   int    `json:"awayGoals,omitempty"`
	YellowCards int    `json:"yellowCards,omitempty"`
	RedCards    int    `json:"redCards,omitempty"`
	ExtraTime   bool   `json:"extraTime,omitempty"`
}

// MatchStore gestiona partidos con persistencia en JSON
type MatchStore struct {
	mu          sync.RWMutex
	matches     map[int]Match
	nextID      int
	storagePath string
}

// NewMatchStore crea un nuevo almacén de partidos
func NewMatchStore(storagePath string) *MatchStore {
	store := &MatchStore{
		matches:     make(map[int]Match),
		nextID:      1,
		storagePath: storagePath,
	}

	// Cargar datos existentes del archivo
	store.loadFromDisk()
	return store
}

// loadFromDisk carga datos desde el archivo JSON
func (ms *MatchStore) loadFromDisk() {
	data, err := ioutil.ReadFile(ms.storagePath)
	if err != nil {
		// Si el archivo no existe, no es un error
		if os.IsNotExist(err) {
			return
		}
		log.Printf("Error al leer el archivo de datos: %v", err)
		return
	}

	var storage struct {
		Matches map[int]Match `json:"matches"`
		NextID  int           `json:"nextID"`
	}

	if err := json.Unmarshal(data, &storage); err != nil {
		log.Printf("Error al deserializar datos: %v", err)
		return
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.matches = storage.Matches
	ms.nextID = storage.NextID
}

// saveToDisk guarda datos en el archivo JSON
func (ms *MatchStore) saveToDisk() error {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	storage := struct {
		Matches map[int]Match `json:"matches"`
		NextID  int           `json:"nextID"`
	}{
		Matches: ms.matches,
		NextID:  ms.nextID,
	}

	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar datos: %v", err)
	}

	if err := ioutil.WriteFile(ms.storagePath, data, 0644); err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	return nil
}

// CreateMatch añade un nuevo partido
func (ms *MatchStore) CreateMatch(match Match) Match {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	match.ID = ms.nextID
	ms.matches[ms.nextID] = match
	ms.nextID++

	// Guardar cambios en disco
	go ms.saveToDisk()

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

// GetMatchByID recupera un partido por su ID
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

	// Guardar cambios en disco
	go ms.saveToDisk()

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

	// Guardar cambios en disco
	go ms.saveToDisk()

	return true
}

// RegisterGoal añade un gol a un partido
func (ms *MatchStore) RegisterGoal(id int) (Match, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	match, exists := ms.matches[id]
	if !exists {
		return Match{}, false
	}

	match.HomeGoals++
	ms.matches[id] = match

	// Guardar cambios en disco
	go ms.saveToDisk()

	return match, true
}

// RegisterYellowCard añade una tarjeta amarilla a un partido
func (ms *MatchStore) RegisterYellowCard(id int) (Match, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	match, exists := ms.matches[id]
	if !exists {
		return Match{}, false
	}

	match.YellowCards++
	ms.matches[id] = match

	// Guardar cambios en disco
	go ms.saveToDisk()

	return match, true
}

// RegisterRedCard añade una tarjeta roja a un partido
func (ms *MatchStore) RegisterRedCard(id int) (Match, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	match, exists := ms.matches[id]
	if !exists {
		return Match{}, false
	}

	match.RedCards++
	ms.matches[id] = match

	// Guardar cambios en disco
	go ms.saveToDisk()

	return match, true
}

// SetExtraTime establece tiempo extra para un partido
func (ms *MatchStore) SetExtraTime(id int) (Match, bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	match, exists := ms.matches[id]
	if !exists {
		return Match{}, false
	}

	match.ExtraTime = true
	ms.matches[id] = match

	// Guardar cambios en disco
	go ms.saveToDisk()

	return match, true
}

// Manejadores de API (sin cambios)
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

func registerGoalHandler(store *MatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de partido inválido", http.StatusBadRequest)
			return
		}

		match, updated := store.RegisterGoal(id)
		if !updated {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
	}
}

func registerYellowCardHandler(store *MatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de partido inválido", http.StatusBadRequest)
			return
		}

		match, updated := store.RegisterYellowCard(id)
		if !updated {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
	}
}

func registerRedCardHandler(store *MatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de partido inválido", http.StatusBadRequest)
			return
		}

		match, updated := store.RegisterRedCard(id)
		if !updated {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
	}
}

func setExtraTimeHandler(store *MatchStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "ID de partido inválido", http.StatusBadRequest)
			return
		}

		match, updated := store.SetExtraTime(id)
		if !updated {
			http.Error(w, "Partido no encontrado", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(match)
	}
}

func main() {
	// Ruta al archivo de almacenamiento JSON
	storagePath := "matches.json"

	// Crear el almacén con persistencia
	store := NewMatchStore(storagePath)

	// Configurar el router
	r := mux.NewRouter()

	// Middleware CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Rutas para partidos
	r.HandleFunc("/api/matches", createMatchHandler(store)).Methods("POST")
	r.HandleFunc("/api/matches", getAllMatchesHandler(store)).Methods("GET")
	r.HandleFunc("/api/matches/{id}", getMatchByIDHandler(store)).Methods("GET")
	r.HandleFunc("/api/matches/{id}", updateMatchHandler(store)).Methods("PUT")
	r.HandleFunc("/api/matches/{id}", deleteMatchHandler(store)).Methods("DELETE")

	// Rutas PATCH para operaciones adicionales
	r.HandleFunc("/api/matches/{id}/goals", registerGoalHandler(store)).Methods("PATCH")
	r.HandleFunc("/api/matches/{id}/yellowcards", registerYellowCardHandler(store)).Methods("PATCH")
	r.HandleFunc("/api/matches/{id}/redcards", registerRedCardHandler(store)).Methods("PATCH")
	r.HandleFunc("/api/matches/{id}/extratime", setExtraTimeHandler(store)).Methods("PATCH")

	// Aplicar middleware CORS
	handler := c.Handler(r)

	// Iniciar servidor
	port := 8080
	fmt.Printf("Servidor iniciando en el puerto %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), handler))
}

package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type CounterStore interface {
	Inc()
	Total() int
}

type MemoryStore struct {
	mu    sync.Mutex
	total int
}

func (s *MemoryStore) Inc() {
	s.mu.Lock()
	s.total++
	s.mu.Unlock()
}

func (s *MemoryStore) Total() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.total
}

type Stats struct {
	Requests int    `json:"requests"`
	Uptime   int    `json:"uptime"`
	Instance string `json:"instance"`
}

func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		value = strings.Trim(value, `"'`)

		os.Setenv(key, value)
	}

	return scanner.Err()
}

func instanceID() string {
	id := os.Getenv("INSTANCE_ID")
	if id == "" {
		id, _ = os.Hostname()
	}
	return id
}

func countRequests(store CounterStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		store.Inc()
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := loadEnv(".env"); err != nil {
		log.Println("pas de .env, valeurs par défaut")
	}

	port := os.Getenv("PING_LISTEN_PORT")
	if port == "" {
		port = "8080"
	}

	var store CounterStore = &MemoryStore{}
	start := time.Now()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "web/index.html")
	})

	mux.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/style.css")
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(r.Header)
	})

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s := Stats{
			Requests: store.Total(),
			Uptime:   int(time.Since(start).Seconds()),
			Instance: instanceID(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	log.Println("serveur lancé sur le port", port)
	log.Fatal(http.ListenAndServe(":"+port, countRequests(store, mux)))
}

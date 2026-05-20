package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"
)

var (
	store = map[string]string{}
	mu    sync.RWMutex
	chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func shortCode(n int) string {
	b := make([]rune, n)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[idx.Int64()]
	}
	return string(b)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed); return
	}
	var body struct{ URL string `json:"url"` }
	json.NewDecoder(r.Body).Decode(&body)
	if body.URL == "" { http.Error(w, "url required", 400); return }
	code := shortCode(6)
	mu.Lock(); store[code] = body.URL; mu.Unlock()
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(map[string]string{"short": "http://localhost:8080/" + code})
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]
	mu.RLock(); url, ok := store[code]; mu.RUnlock()
	if !ok { http.NotFound(w, r); return }
	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/", redirectHandler)
	fmt.Println("URL shortener running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

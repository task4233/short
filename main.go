package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	m       map[string]string
	handler http.Handler
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) init() {
	s.m = make(map[string]string)
	s.m["test"] = "https://example.com"

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/urls", s.URLs)
	r.Post("/register", s.Register)
	r.Get("/", s.Index)
	s.handler = r
}

// Data stores mapping data for redirect
type Data struct {
	Key string
	URL string
}

// Index(GET /) redirects stored data
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[INDEX]: ", r.URL.Path[1:])
	u, ok := s.m[r.URL.Path[1:]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found\n")
		return
	}
	http.Redirect(w, r, u, http.StatusMovedPermanently)
}

// GET /urls
func (s *Server) URLs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[URLs]: ", r.URL.Path[1:])
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.m)
}

// POST /register
// {"key": "test", "URL": "https://example.com"}
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[Register]: ", r.URL.Path[1:])
	var data Data
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "bad request body\n")
		return
	}

	// if key is not set, generate md5 hash
	if data.Key == "" {
		hash := md5.New()
		key := fmt.Sprintf("%x", hash.Sum([]byte(data.URL)))
		s.m[key] = data.URL

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "register successed!\nkey: %s, URL: %s\n", key, data.URL)
		return
	}

	// if key has already been used, its request is bad request
	if _, ok := s.m[data.Key]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "URL has already been registered\n")
		return
	}

	s.m[data.Key] = data.URL

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "register successed!\nkey: %s, URL: %s\n", data.Key, data.URL)
}

func main() {
	s := NewServer()
	s.init()
	http.ListenAndServe(":8080", s.handler)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	mux := http.NewServeMux()

	mux.HandleFunc("/urls", s.URLs)
	mux.HandleFunc("/register", s.Register)
	mux.HandleFunc("/", s.Index)
	s.handler = mux
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
		fmt.Fprintf(w, "not found")
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
	} else if data.Key == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "key must not be empty\n")
		return
	}

	if s.m[data.Key] != "" {
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

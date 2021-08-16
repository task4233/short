package main

import (
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
	mux.HandleFunc("/", s.Index)
	s.handler = mux
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("path: ", r.URL.Path[1:])
	u, ok := s.m[r.URL.Path[1:]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "not found")
		return
	}
	http.Redirect(w, r, u, http.StatusMovedPermanently)
}

func main() {
	s := NewServer()
	s.init()
	http.ListenAndServe(":8080", s.handler)
}

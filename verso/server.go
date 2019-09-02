package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func (s *Server) Initialize() {
	s.setRouters()
}

func (s *Server) setRouters() {
	fmt.Println("setroute")
	s.Post("/reverse", HandleReverse)
}

func (s *Server) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	fmt.Println("POST rergieste")
	s.router.HandleFunc(path, f).Methods("POST")
}

func (s *Server) Run(apiHost string) {
	fmt.Println("RUN", apiHost)
	http.ListenAndServe(apiHost, nil)
}

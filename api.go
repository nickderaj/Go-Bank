package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string
}

type ApiServer struct {
	listenAddr string
}

func handleRequest(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			e := WriteJSON(w, http.StatusInternalServerError, ApiError{err.Error()})
			if e != nil {
				return
			}
		}
	}
}

func NewAPIServer(listenAddr string) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", handleRequest(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account/{id}", handleRequest(s.handleGetAccount)).Methods("GET")
	router.HandleFunc("/account/{id}", handleRequest(s.handleDeleteAccount)).Methods("DELETE")
	router.HandleFunc("/transfer", handleRequest(s.handleTransfer)).Methods("POST")

	log.Println("Listening on", s.listenAddr)
	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal("Error starting the server: ", err)
		return
	}
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)

	account := newAccount("Nick", "De Raj")
	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)

	return nil
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

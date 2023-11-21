package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

type ApiServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", handleRequest(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/account", handleRequest(s.handleGetAccounts)).Methods("GET")
	router.HandleFunc("/account/{id}", withJWTAuth(handleRequest(s.handleGetAccount), s.store)).Methods("GET")
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
	id, err := parseId(r)
	if err != nil {
		return err
	}

	account, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := newAccount(createAccountReq.FirstName, createAccountReq.LastName)

	tokenString, err := createJWT(account)

	if err != nil {
		return err
	}
	fmt.Println("Token string: ", tokenString)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := parseId(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted:": id})
}

func (s *ApiServer) handleGetAccounts(w http.ResponseWriter, _ *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error closing body: ", err)
		}
	}(r.Body)

	return WriteJSON(w, http.StatusOK, transferReq)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
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

func parseId(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid account id: %s", idStr)
	}
	return id, nil
}

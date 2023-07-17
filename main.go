package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/peksinsara/e-voting-RDBMS/function"
	"github.com/peksinsara/e-voting-RDBMS/user"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/register", user.RegisterUser).Methods("POST")
	router.HandleFunc("/login", user.LoginUser).Methods("POST")

	adminRouter := router.PathPrefix("/admin").Subrouter()

	adminRouter.HandleFunc("/candidate", user.AddCandidate).Methods("POST")
	adminRouter.HandleFunc("/alldata", user.GetAllData).Methods("GET")
	adminRouter.HandleFunc("/candidate/{id}", user.DeleteCandidate).Methods("DELETE")

	router.HandleFunc("/vote", function.CastVote).Methods("POST")

	addr := ":8000"

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	log.Printf("Successfully started")
	log.Printf("Running application E-voting-RDBMS on port %s", addr)
	log.Fatal(server.ListenAndServe())
}

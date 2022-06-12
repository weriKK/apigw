package main

import (
	"apigw/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/services", service.RegistrationHandler).Methods("POST")

	http.ListenAndServe(":8081", r)
}

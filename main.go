package main

import (
	"apigw/internal/service"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func readyHandler(rw http.ResponseWriter, req *http.Request) {
	req.Body.Close()
	rw.WriteHeader(http.StatusOK)
}

func notReadyHandler(rw http.ResponseWriter, req *http.Request) {
	req.Body.Close()
	rw.WriteHeader(http.StatusInternalServerError)
}

func timeoutHandler(rw http.ResponseWriter, req *http.Request) {
	req.Body.Close()
	time.Sleep(2 * time.Second)
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/services", service.RegistrationHandler).Methods("POST")

	r.HandleFunc("/readyz", readyHandler).Methods("GET")
	r.HandleFunc("/notreadyz", notReadyHandler).Methods("GET")
	r.HandleFunc("/timeoutz", timeoutHandler).Methods("GET")

	http.ListenAndServe(":8081", r)
}

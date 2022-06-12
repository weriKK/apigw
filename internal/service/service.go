package service

import (
	"apigw/internal/util"
	"fmt"
	"net/http"
)

type ServiceRegistration struct {
	Name          string            `json:"name"`
	ReadinessPath string            `json:"readinessPath,omitempty"`
	Endpoints     []ServiceEndpoint `json:"endpoints"`
}

type ServiceEndpoint struct {
	ApiRoot string   `json:"apiRoot"`
	Methods []string `json:"methods"`
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {

	var sr *ServiceRegistration
	var err error

	if sr, err = parseServiceRegistration(r); err != nil {
		util.JSONError(w, err, http.StatusBadRequest)
		return
	}

	fmt.Println(sr.Name)

	w.WriteHeader(http.StatusOK)
}

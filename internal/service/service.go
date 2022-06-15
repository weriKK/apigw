package service

import (
	"apigw/internal/util"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ServiceRegistration struct {
	Name          string            `json:"name"`
	Host          string            `json:"host"`
	Port          uint              `json:"port"`
	ReadinessPath string            `json:"readinessPath,omitempty"`
	Endpoints     []ServiceEndpoint `json:"endpoints"`
}

type ServiceEndpoint struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
}

// http://apigw:8080/{Name}/{Path} -> http://{Host}:{Port}/{Name}/{Path}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {

	var sr *ServiceRegistration
	var err error

	if sr, err = parseServiceRegistration(r); err != nil {
		util.JSONError(w, err, http.StatusBadRequest)
		return
	}

	RegisterService(sr)

	w.WriteHeader(http.StatusOK)
}

func RegisterService(sr *ServiceRegistration) {

	apiRoot := fmt.Sprintf("http://%s:%d", sr.Host, sr.Port)

	si := ServiceInformation{
		Name:         sr.Name,
		ApiRoot:      apiRoot,
		ReadinessURI: fmt.Sprintf("%s/%s", apiRoot, sr.ReadinessPath),
		Status:       ServiceStatusUnknown,
		Endpoints:    sr.Endpoints,
	}

	_DB.SaveService(&si)
}

func init() {
	CheckServices()
}

func CheckServices() {

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	// quit is worthless now, can't stop this!

	go func() {
		for {
			select {
			case <-ticker.C:
				for name, info := range _DB.FindServices() {
					go ReadinessCheck(name, info.ReadinessURI)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

}

func ReadinessCheck(name, url string) {
	client := http.Client{
		Timeout: 250 * time.Millisecond,
	}

	newStatus := ServiceStatusUnknown

	resp, err := client.Get(url)
	if err != nil {
		newStatus = ServiceStatusUnreachable
	} else if resp.StatusCode != http.StatusOK {
		newStatus = ServiceStatusNotReady
	} else {
		newStatus = ServiceStatusReady
	}

	log.Printf("Service status update: %s [%s]", newStatus, name)
	_DB.UpdateServiceStatus(name, newStatus)
}

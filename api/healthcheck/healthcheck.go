// Package healthcheck to check application disponibility
package healthcheck

import (
	"go-boilerplate/common/response"
	"go-boilerplate/repository"
	"net/http"
)

type simpleResponse struct {
	Status string `json:"status"`
}

// SimpleHandler handles simple healthcheck requests
func SimpleHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	response.Write(w, simpleResponse{
		Status: "OK",
	}, http.StatusOK)
}

// CompleteHandler handles complete healthcheck requests
func CompleteHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	healthy := repository.Healthcheck()
	if !healthy.Healthy() {
		response.Write(w, healthy, http.StatusServiceUnavailable)
		return
	}
	response.Write(w, healthy, http.StatusOK)
}

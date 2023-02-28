package api_test

import (
	"go-boilerplate/test"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	testCases := []test.APITestCase{
		{
			Name:   "complete healthcheck",
			Route:  "http://localhost:9000/healthcheck",
			Method: http.MethodGet,
			Status: http.StatusOK,
			Body:   `{"DB":"OK","SQS":"OK","HTTP":"OK"}`,
		},
		{
			Name:   "simple healthcheck",
			Route:  "http://localhost:9000/healthcheck/status",
			Method: http.MethodGet,
			Status: http.StatusOK,
			Body:   `{"status":"OK"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, tc.Run)
	}
}

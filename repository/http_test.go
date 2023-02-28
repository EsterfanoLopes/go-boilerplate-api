package repository_test

import (
	"fmt"
	"go-boilerplate/repository"
	"go-boilerplate/test"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	retryDisabled = iota
	retryEnabled
)

func TestMain(m *testing.M) {
	err := repository.Setup()
	if err != nil {
		fmt.Printf("error starting http tests %s \n", err)
		os.Exit(-1)
	}
	os.Exit(m.Run())
}

type anything struct {
	ID   int    `json:"id"`
	Desc string `json:"desc"`
}

func TestExecuteAndParseHTTPResponse(t *testing.T) {
	retrySlowWait := 3
	test.MockHTTP(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			b := []byte("ERROR")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(b)
			return
		}
		if r.URL.Path == "/not-found" {
			b := []byte("NOT_FOUND")
			w.WriteHeader(http.StatusNotFound)
			w.Write(b)
			return
		}
		if r.URL.Path == "/unauthorized" {
			b := []byte("UNAUTHORIZED")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(b)
			return
		}
		if r.URL.Path == "/anything" {
			b := []byte(`{"id":1,"desc":"anything"}`)
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		if r.URL.Path == "/slow" {
			time.Sleep(1 * time.Second)
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.URL.Path == "/slow-retry" {
			retrySlowWait -= 1
			time.Sleep(time.Duration(retrySlowWait) * time.Second)
			b := []byte{}
			if retrySlowWait < 2 {
				b = []byte(`{"id":1,"desc":"anything"}`)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(b)
			return
		}
		if r.URL.Path == "/empty-response" {
			w.WriteHeader(http.StatusOK)
			return
		}
	})

	testCases := []struct {
		name          string
		method        string
		path          string
		body          io.Reader
		header        *http.Header
		timeout       time.Duration
		expected      anything
		expectedError string
		retry         int
	}{
		{
			name:    "get executed successfully",
			method:  http.MethodGet,
			path:    "anything",
			timeout: 1 * time.Second,
			expected: anything{
				ID:   1,
				Desc: "anything",
			},
			expectedError: "",
			retry:         retryDisabled,
		},
		{
			name:    "post executed successfully",
			method:  http.MethodPost,
			path:    "anything",
			timeout: 1 * time.Second,
			body:    strings.NewReader(`{"a":1,"b":"c"}`),
			expected: anything{
				ID:   1,
				Desc: "anything",
			},
			expectedError: "",
			retry:         retryDisabled,
		},
		{
			name:          "not found resource",
			method:        http.MethodGet,
			path:          "not-found",
			timeout:       1 * time.Second,
			expected:      anything{},
			expectedError: "resource not found",
			retry:         retryDisabled,
		},
		{
			name:          "unauthorized resource",
			method:        http.MethodGet,
			path:          "unauthorized",
			timeout:       1 * time.Second,
			expected:      anything{},
			expectedError: "unauthorized resource",
			retry:         retryDisabled,
		},
		{
			name:          "error executing get request",
			method:        http.MethodGet,
			path:          "error",
			timeout:       1 * time.Second,
			expected:      anything{},
			expectedError: "error executing GET http://127.0.0.1:8001/error - 500 - ERROR",
			retry:         retryDisabled,
		},
		{
			name:          "timeout executing get request",
			method:        http.MethodGet,
			path:          "slow",
			timeout:       500 * time.Millisecond,
			expected:      anything{},
			expectedError: `Get "http://127.0.0.1:8001/slow": context deadline exceeded`,
			retry:         retryEnabled,
		},
		{
			name:    "get executed successfully after retrying",
			method:  http.MethodGet,
			path:    "slow-retry",
			timeout: 1500 * time.Millisecond,
			expected: anything{
				ID:   1,
				Desc: "anything",
			},
			expectedError: "",
			retry:         retryEnabled,
		},
		{
			name:          "empty body response when some data is expected",
			method:        http.MethodGet,
			path:          "empty-response",
			timeout:       1 * time.Second,
			expected:      anything{},
			expectedError: "",
			retry:         retryDisabled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := anything{}
			err := repository.ExecuteAndParseHTTPResponse(tc.method, fmt.Sprintf("http://127.0.0.1:8001/%s", tc.path),
				&result, tc.body, tc.header, tc.timeout, tc.retry)
			if tc.expectedError != "" && tc.expectedError != err.Error() {
				t.Errorf("unexpected error executing get %s", err)
				return
			}
			if tc.expectedError == "" && err != nil {
				t.Errorf("unexpected error executing get %s", err)
				return
			}
			if cmp.Diff(tc.expected, result) != "" {
				t.Errorf("unexpected result of get request %s", cmp.Diff(tc.expected, result))
				return
			}
		})
	}
}

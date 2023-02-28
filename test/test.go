// Package test keep main test functions to be used by its appendages
package test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/google/go-cmp/cmp"
)

// MockHTTP provides a mocked server with the given handler func
func MockHTTP(t *testing.T, handler http.HandlerFunc) {
	mock := httptest.NewUnstartedServer(
		http.HandlerFunc(handler),
	)
	mock.Listener.Close()
	l, err := net.Listen("tcp", "127.0.0.1:8001")
	if err != nil {
		t.Fatalf("error binding 127.0.0.1:8001 %+v", err)
		return
	}
	mock.Listener = l
	mock.Start()
	t.Cleanup(func() {
		mock.Close()
	})
	time.Sleep(500 * time.Millisecond)
}

// MockAuthorization mocks authorization in account-api
func MockAuthorization(t *testing.T, accountID string) {
	MockHTTP(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/user-info" && r.Method == http.MethodGet {
			body := []byte(fmt.Sprintf(`{"name":"Updated","uuid":"%s","authorities":[],"email":"updated@mailinator.com"}`, accountID))
			w.WriteHeader(http.StatusOK)
			w.Write(body)
			return
		}
		w.WriteHeader(http.StatusNotImplemented)
	})
}

// APITestCase basic struct to test cases
type APITestCase struct {
	Name    string
	Route   string
	Method  string
	Status  int
	Payload string
	Body    string
	Headers http.Header
}

// Run execute test cases
func (tc APITestCase) Run(t *testing.T) {
	var resp *http.Response
	var err error
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest(tc.Method, tc.Route, bytes.NewBuffer([]byte(tc.Payload)))
	if err != nil {
		t.Errorf("Error in request %s", err)
	}
	if tc.Headers != nil {
		req.Header = tc.Headers
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Errorf("io error calling %s %s", tc.Name, err)
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("body read error calling %s %s", tc.Name, err)
		return
	}

	body := string(bytes)
	if resp.StatusCode != tc.Status {
		t.Errorf("status error calling %s %d %s", tc.Name, resp.StatusCode, body)
		return
	}

	if resp.StatusCode/100 != 3 {
		if diff := cmp.Diff(body, tc.Body); diff != "" {
			t.Errorf("unexpected response body in '%s'\n %s", tc.Name, diff)
			return
		}
	}
}

// AssertError asserts if an error message is the expected
func AssertError(t *testing.T, err error, expectedError string) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	if msg != expectedError {
		t.Errorf("unexpected error message %s", msg)
	}
}

// AssertErrorPrefix asserts if an error message starts with the expected
func AssertErrorPrefix(t *testing.T, err error, expectedError string) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	if expectedError == "" && msg != "" {
		t.Errorf("unexpected error message prefix %s", msg)
	}

	if !strings.HasPrefix(msg, expectedError) {
		t.Errorf("unexpected error message prefix %s", msg)
	}
}

// AssertErrorType asserts if an error message is the expected
func AssertErrorType(t *testing.T, err error, expectedError error) {
	if err == nil && expectedError == nil {
		return
	}
	if !errors.Is(err, expectedError) {
		t.Errorf("unexpected error type %s", err)
	}
}

// FreezeTime freezes the time for test purposes
func FreezeTime(t *testing.T, ref time.Time) {
	patch := monkey.Patch(time.Now, func() time.Time { return ref })
	t.Cleanup(func() {
		patch.Unpatch()
	})
}

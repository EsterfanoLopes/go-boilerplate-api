package response_test

import (
	"errors"
	"go-boilerplate/common/response"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWrite(t *testing.T) {
	body := struct {
		X string `json:"x"`
		Y int    `json:"y"`
	}{
		X: "1",
		Y: 2,
	}
	rw := httptest.NewRecorder()
	response.Write(rw, body, http.StatusOK)
	rw.Flush()
	if rw.Code != http.StatusOK {
		t.Errorf("unexpected status code %d", rw.Code)
	}
	if rw.Header().Get("content-type") != "application/json" {
		t.Errorf("unexpected content type %s", rw.Header().Get("content-type"))
	}
	if rw.Body.String() != `{"x":"1","y":2}` {
		t.Errorf("unexpected body %s", rw.Body.String())
	}
}

func TestWriteError(t *testing.T) {
	testCases := []struct {
		name           string
		err            error
		message        string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "unknown error",
			err:            errors.New("timeout bla blabla"),
			message:        "error message",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"code":"CAS001","error":"timeout bla blabla"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rw := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "http://test.com.br/v1/test", nil)
			response.WriteError(rw, r, tc.err, tc.message)
			rw.Flush()

			if rw.Code != tc.expectedStatus {
				t.Errorf("unexpected status code %d", rw.Code)
				return
			}
			if rw.Body.String() != tc.expectedBody {
				t.Errorf("unexpected body %s", rw.Body.String())
				return
			}
		})
	}
}

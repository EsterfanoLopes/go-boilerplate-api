package repository

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"go-boilerplate/common"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/rs/dnscache"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

const (
	retryDisabled = iota
	retryEnabled
)

var (
	httpTimeout             = common.Config.GetInt("httpTimeoutSeconds")
	httpHealthcheckEndpoint = common.Config.Get("httpHealthcheckEndpoint")
	httpMaxRetries          = common.Config.GetInt("httpMaxRetries")
)

var (
	// ErrUnauthorizedResource used to say that a resource can't be accessed by authorization policies
	ErrUnauthorizedResource = errors.New("unauthorized resource")
	// ErrUnprocessableEntityResource used to say that the request data is unprocessable
	ErrUnprocessableEntityResource = errors.New("unprocessable request")
)

var httpTimeoutDuration = time.Duration(httpTimeout) * time.Second

var customHTTPTransport = &http.Transport{
	Proxy:                 http.ProxyFromEnvironment,
	DialContext:           timeoutDialer(httpTimeoutDuration, httpTimeoutDuration),
	MaxIdleConns:          common.Config.GetInt("httpMinConnections") * 10,
	MaxIdleConnsPerHost:   common.Config.GetInt("httpMaxConnections"),
	IdleConnTimeout:       httpTimeoutDuration,
	TLSHandshakeTimeout:   httpTimeoutDuration,
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	ExpectContinueTimeout: httpTimeoutDuration,
}

func timeoutDialer(connectionTimeut time.Duration, readTimeut time.Duration) func(ctx context.Context, net, addr string) (c net.Conn, err error) {
	r := &dnscache.Resolver{}

	go func() {
		t := time.NewTicker(5 * time.Minute)
		defer t.Stop()
		for range t.C {
			r.Refresh(true)
		}
	}()

	return func(ctx context.Context, netw, addr string) (conn net.Conn, err error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		ips, err := r.LookupHost(ctx, host)
		if err != nil {
			return nil, err
		}
		for _, ip := range ips {
			dialer := net.Dialer{Timeout: connectionTimeut}
			conn, err = dialer.DialContext(ctx, netw, net.JoinHostPort(ip, port))
			if err == nil {
				break
			}
		}

		if err != nil {
			return
		}

		conn.SetDeadline(time.Now().Add(readTimeut))

		return
	}
}

// HTTP is an http client instance
var HTTP *http.Client

var httpReady = int32(0)

func httpIsReady() {
	atomic.StoreInt32(&httpReady, 1)
}

func isHTTPReady() bool {
	return atomic.LoadInt32(&httpReady) == 1
}

func setupHTTP() {
	if isHTTPReady() {
		return
	}
	HTTP = &http.Client{
		Transport: customHTTPTransport,
		Timeout:   httpTimeoutDuration,
	}
	HTTP = httptrace.WrapClient(
		HTTP,
		httptrace.RTWithServiceName("go-boilerplate-http-client"),
		httptrace.RTWithResourceNamer(func(req *http.Request) string {
			return strings.Split(req.URL.String(), "?")[0]
		}),
	)
	httpIsReady()
}

// ExecuteAndParseHTTPResponse executes the given request with default http instance
func ExecuteAndParseHTTPResponse(
	method, url string,
	result interface{},
	body io.Reader,
	header *http.Header,
	timeout time.Duration,
	attempt int) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	if header != nil {
		request.Header = *header
	}
	response, err := HTTP.Do(request)
	defer CloseBody(response)
	if err != nil {
		if attempt >= retryEnabled && attempt <= httpMaxRetries {
			common.Logger.Warnf("retrying, method: %s, url: %s, attempt: %d, timeout: %v, err: %v", method, url, attempt, timeout, err)
			attempt++
			time.Sleep(time.Duration(attempt) * time.Second)
			return ExecuteAndParseHTTPResponse(method, url, result, body, header, timeout, attempt)
		}
		return err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode == http.StatusNotFound {
		common.Logger.Errorf("error executing %s %s - %d - %s", method, url, response.StatusCode, string(bytes))
		return ErrNotFound
	}
	if response.StatusCode == http.StatusUnauthorized {
		common.Logger.Errorf("error executing %s %s - %d - %s", method, url, response.StatusCode, string(bytes))
		return ErrUnauthorizedResource
	}
	if response.StatusCode == http.StatusForbidden {
		common.Logger.Errorf("error executing %s %s - %d - %s", method, url, response.StatusCode, string(bytes))
		return ErrUnauthorizedResource
	}
	if response.StatusCode == http.StatusUnprocessableEntity {
		common.Logger.Errorf("error executing %s %s - %d - %s", method, url, response.StatusCode, string(bytes))
		return ErrUnprocessableEntityResource
	}
	if response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("error executing %s %s - %d - %s", method, url, response.StatusCode, string(bytes))
	}

	if common.Config.Get("httpResponseDebug") == "true" {
		common.Logger.Debugf("response from %s %s - %s", method, url, string(bytes))
	}

	if result != nil {
		if len(bytes) == 0 {
			common.Logger.Warnf("response from %s %s has an empty body", method, url)
			return nil
		}

		err = json.Unmarshal(bytes, &result)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAndParseHTTPResponse perform a GET request in the given url and parse its results
func GetAndParseHTTPResponse(url string, result interface{}, header *http.Header, timeout time.Duration) error {
	return ExecuteAndParseHTTPResponse(http.MethodGet, url, result, nil, header, timeout, retryEnabled)
}

// PostAndParseHTTPResponse perform a POST request in the given url with given body an parse its results
func PostAndParseHTTPResponse(url string, result interface{}, body io.Reader, header *http.Header, timeout time.Duration) error {
	return ExecuteAndParseHTTPResponse(http.MethodPost, url, result, body, header, timeout, retryDisabled)
}

// PatchAndParseHTTPResponse perform a PATCH request in the given url with given body an parse its results
func PatchAndParseHTTPResponse(url string, result interface{}, body io.Reader, header *http.Header, timeout time.Duration) error {
	return ExecuteAndParseHTTPResponse(http.MethodPatch, url, result, body, header, timeout, retryDisabled)
}

// PutAndParseHTTPResponse perform a PUT request in the given url with given body an parse its results
func PutAndParseHTTPResponse(url string, result interface{}, body io.Reader, header *http.Header, timeout time.Duration) error {
	return ExecuteAndParseHTTPResponse(http.MethodPut, url, result, body, header, timeout, retryDisabled)
}

// DeleteAndParseHTTPResponse perform a DELETE request in the given url with given body an parse its results
func DeleteAndParseHTTPResponse(url string, result interface{}, body io.Reader, header *http.Header, timeout time.Duration) error {
	return ExecuteAndParseHTTPResponse(http.MethodDelete, url, result, body, header, timeout, retryDisabled)
}

// CloseBody closes the given response body
func CloseBody(response *http.Response) {
	if response == nil {
		return
	}
	err := response.Body.Close()
	if err != nil {
		common.Logger.Errorf("error closing response body", err)
	}
}

func httpHealthcheck() error {
	if httpHealthcheckEndpoint == "" {
		return nil
	}
	response, err := HTTP.Get(httpHealthcheckEndpoint)
	if err != nil {
		return err
	}
	defer CloseBody(response)
	return nil
}

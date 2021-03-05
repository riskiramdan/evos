package hosts

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

// HTTPManager interface
type HTTPManager struct{}

// ConfigHTTP config struct
type ConfigHTTP struct {
	DebugClient bool   `envconfig:"DEBUG_CLIENT" default:"true"`
	Timeout     string `envconfig:"TIMEOUT" default:"60s"`
	RetryBad    int    `envconfig:"RETRY_BAD" default:"1"`
}

var (
	envHTTP ConfigHTTP
)

// HTTPGet func
func (hm *HTTPManager) HTTPGet(url string, header http.Header) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(envHTTP.DebugClient)
	timeout, _ := time.ParseDuration(envHTTP.Timeout)
	reqagent := request.Get(url)
	header.Set("Accept-Encoding", "identity")
	reqagent.Header = header
	_, body, errs := reqagent.
		Timeout(timeout).
		Retry(envHTTP.RetryBad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}

// HTTPPost func
func (hm *HTTPManager) HTTPPost(url string, jsondata interface{}) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(envHTTP.DebugClient)
	timeout, _ := time.ParseDuration(envHTTP.Timeout)
	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Post(url)
	reqagent.Header.Set("Content-Type", "application/json")
	_, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(envHTTP.RetryBad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}

// HTTPPostWithHeader func
func (hm *HTTPManager) HTTPPostWithHeader(url string, jsondata interface{}, header http.Header) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(envHTTP.DebugClient)
	timeout, _ := time.ParseDuration(envHTTP.Timeout)
	reqagent := request.Post(url)
	reqagent.Header = header
	_, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(envHTTP.RetryBad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}

// HTTPPutWithHeader func
func (hm *HTTPManager) HTTPPutWithHeader(url string, jsondata interface{}, header http.Header) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(envHTTP.DebugClient)
	timeout, _ := time.ParseDuration(envHTTP.Timeout)
	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Put(url)
	reqagent.Header = header
	_, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(envHTTP.RetryBad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}

// HTTPDeleteWithHeader func
func (hm *HTTPManager) HTTPDeleteWithHeader(url string, jsondata interface{}, header http.Header) ([]byte, error) {
	request := gorequest.New()
	request.SetDebug(envHTTP.DebugClient)
	timeout, _ := time.ParseDuration(envHTTP.Timeout)
	if url[:5] == "https" {
		request.TLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	reqagent := request.Delete(url)
	reqagent.Header = header
	_, body, errs := reqagent.
		Send(jsondata).
		Timeout(timeout).
		Retry(envHTTP.RetryBad, time.Second, http.StatusInternalServerError).
		End()
	if errs != nil {
		return []byte(body), errs[0]
	}
	return []byte(body), nil
}

package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Require a real FTP running

func TestProxy_GET(t *testing.T) {
	req, _ := http.NewRequest("GET", "/config.yml", nil)
	rr := httptest.NewRecorder()

	target, _ := url.Parse("ftp://bob:azerty@localhost:21")
	handler := http.HandlerFunc(proxy(configuration{
		port:   8080,
		target: *target,
	}))

	handler.ServeHTTP(rr, req)

	checkStatus(t, rr, http.StatusOK)
}

func TestProxy_GET_404(t *testing.T) {
	req, _ := http.NewRequest("GET", "/doesnotexist.yml", nil)
	rr := httptest.NewRecorder()

	target, _ := url.Parse("ftp://bob:azerty@localhost:21")
	handler := http.HandlerFunc(proxy(configuration{
		port:   8080,
		target: *target,
	}))

	handler.ServeHTTP(rr, req)

	checkStatus(t, rr, http.StatusNotFound)
}

func TestProxy_POST(t *testing.T) {
	req, _ := http.NewRequest("POST", "/file.txt", strings.NewReader("content"))
	rr := httptest.NewRecorder()

	target, _ := url.Parse("ftp://bob:azerty@localhost:21")
	handler := http.HandlerFunc(proxy(configuration{
		port:   8080,
		target: *target,
	}))

	handler.ServeHTTP(rr, req)

	checkStatus(t, rr, http.StatusOK)
}

func TestProxy_OtherHttpMethod(t *testing.T) {
	req, _ := http.NewRequest("DELETE", "/file.txt", nil)
	rr := httptest.NewRecorder()

	target, _ := url.Parse("ftp://bob:azerty@localhost:21")
	handler := http.HandlerFunc(proxy(configuration{
		port:   8080,
		target: *target,
	}))

	handler.ServeHTTP(rr, req)

	checkStatus(t, rr, http.StatusMethodNotAllowed)
}

func TestProxy_InvalidLoginPassword(t *testing.T) {
	req, _ := http.NewRequest("GET", "/file.txt", nil)
	rr := httptest.NewRecorder()

	target, _ := url.Parse("ftp://invalid:azerty@localhost:21")
	handler := http.HandlerFunc(proxy(configuration{
		port:   8080,
		target: *target,
	}))

	handler.ServeHTTP(rr, req)

	checkStatus(t, rr, http.StatusUnauthorized)
}

func checkStatus(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) {
	if status := rr.Code; status != expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}
}

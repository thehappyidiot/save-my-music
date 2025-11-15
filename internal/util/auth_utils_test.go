package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TODO: Generate valid token without using actual token?
const SAMPLE_JWT = ""

func createSampleRequest() *http.Request {
	request := httptest.NewRequest("POST",
		"http://example.com",
		bytes.NewReader([]byte(
			fmt.Sprintf("credential=%s&g_csrf_token=sometoken", SAMPLE_JWT),
		)),
	)
	request.AddCookie(&http.Cookie{
		Name:  "g_csrf_token",
		Value: "sometoken",
	})
	request.Header.Add("content-type", "application/x-www-form-urlencoded")

	return request
}

func TestValidateCsrfValid(t *testing.T) {
	request := createSampleRequest()

	err := validateCsrf(request)

	if err != nil {
		t.Errorf("Expected: no error. Got: %s", err)
	}
}

func TestValidateCsrfNoCookie(t *testing.T) {
	request := createSampleRequest()

	request.Header.Del("cookie")

	err := validateCsrf(request)

	if err == nil {
		t.Errorf("Expected: error, got: no error")
	}
}

func TestValidateCsrfNoBody(t *testing.T) {
	request := createSampleRequest()
	request.Body = nil

	err := validateCsrf(request)

	if err == nil {
		t.Errorf("Expected: error, got: no error")
	}
}

func TestValidateCsrfMismatch(t *testing.T) {
	request := createSampleRequest()
	request.Body = io.NopCloser(
		strings.NewReader("g_csrf_token=someothertoken?"),
	)

	err := validateCsrf(request)

	if err == nil {
		t.Errorf("Expected: error, got: no error")
	}
}

/*
*
func TestValidateIdTokenValid(t *testing.T) {
	request := createSampleRequest()

	payload, err := validateIdToken(request, "clientid")

	if err != nil {
		t.Errorf("Expected: no error, got: %s", err)
	}

	actualName := fmt.Sprintf("%v", payload.Claims["name"])
	if actualName != "May, James" {
		t.Errorf("Expected: `May, James`, got: %s", actualName)
	}
} */

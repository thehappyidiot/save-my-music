package util

import (
	"fmt"
	"net/http"

	"google.golang.org/api/idtoken"
)

const GOOGLE_CSRF_KEY = "g_csrf_token"

func ValidateGoogleAuthRequest(req *http.Request, googleClientId string) (*idtoken.Payload, error) {
	// https://developers.google.com/identity/gsi/web/guides/verify-google-id-token

	err := req.ParseForm()
	if err != nil {
		return &idtoken.Payload{}, fmt.Errorf(
			"bad request: %s", err,
		)
	}

	err = validateCsrf(req)
	if err != nil {
		return &idtoken.Payload{}, err
	}

	return validateIdToken(req, googleClientId)
}

func validateCsrf(req *http.Request) error {
	csrfTokenCookie := req.CookiesNamed(GOOGLE_CSRF_KEY)
	if len(csrfTokenCookie) != 1 {
		return fmt.Errorf(
			"request did not contain a unique CSRF token cookie - `%s`",
			GOOGLE_CSRF_KEY,
		)
	}

	csrfTokenBody := req.FormValue(GOOGLE_CSRF_KEY)
	if csrfTokenBody == "" {
		return fmt.Errorf(
			"request did not contain a CSRF token in body - `%s`",
			GOOGLE_CSRF_KEY,
		)
	}

	if csrfTokenCookie[0].Value != csrfTokenBody {
		return fmt.Errorf(
			"failed to verify double submit CSRF token",
		)
	}

	return nil
}

func validateIdToken(req *http.Request, googleClientId string) (*idtoken.Payload, error) {
	idTokenString := req.FormValue("credential")

	payload, err := idtoken.Validate(req.Context(), idTokenString, googleClientId)

	if err != nil {
		return &idtoken.Payload{}, fmt.Errorf(
			"failed to validate idToken: %s", err,
		)
	}

	return payload, nil
}


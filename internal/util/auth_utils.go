package util

import (
	"fmt"
	"net/http"
	"strings"
)

const AUTH_HEADER_KEY = "Authorization"
const AUTH_SCHEME_BEARER = "Bearer"

func GetBearerToken(headers http.Header) (string, error) {
	authHeaderValue := headers.Get(AUTH_HEADER_KEY)

	if authHeaderValue == "" {
		return "", fmt.Errorf("authorization header not found or is empty")
	}

	authHeaderParts := strings.Split(authHeaderValue, " ")

	if len(authHeaderParts) != 2 {
		return "", fmt.Errorf(
			"invalid authorization header",
		)
	}

	if authHeaderParts[0] != AUTH_SCHEME_BEARER {
		return "", fmt.Errorf(
			"authorization header with %s scheme not present",
			AUTH_SCHEME_BEARER,
		)
	}

	token := authHeaderParts[1]

	return token, nil
}

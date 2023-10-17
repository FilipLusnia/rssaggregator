package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authVal := headers.Get("Authorization")
	if authVal == "" {
		return "", errors.New("no auth info found")
	}

	authVals := strings.Split(authVal, " ")
	if (len(authVals) != 2) || (authVals[0] != "ApiKey") {
		return "", errors.New("malformed auth header")
	}

	return authVals[1], nil
}

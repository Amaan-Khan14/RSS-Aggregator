package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetApiKey returns the API key from
// the headers of http request
// Example:
// Authorization:ApiKey {insert api key here}
func GetApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("API Key not found")
	}

	vals := strings.Split(val, " ")
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first path of auth header")
	}
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	return vals[1], nil
}

package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	key, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("authorization header does not exist")
	}

	for _, k := range key {
		if strings.Contains(k, "ApiKey") {
			return k[7:], nil
		}
	}

	return "", fmt.Errorf("authorization header does not contain ApiKey")
}

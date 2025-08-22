package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	token, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("authorization header does not exist")
	}

	for _, tok := range token {
		if strings.Contains(tok, "Bearer") {
			return tok[7:], nil
		}
	}

	return "", fmt.Errorf("authorization header does not contain Bearer")
}

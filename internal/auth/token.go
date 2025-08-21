package auth

import (
	"fmt"
	"net/http"
)

func GetBearerToken(headers http.Header) (string, error) {
	token, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("authorization header does not exist")
	}

	return token[0][7:], nil
}

package utils

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func GetHeaderToken(header http.Header) (string, error) {
	auth := header.Get("Authorization")
	if auth == "" {
		return "", errors.New("Authorization header is required")
	}
	if !strings.Contains(auth, "Bearer ") {
		return "", errors.New("Authorization header must be Bearer")
	}
	return strings.Split(auth, " ")[1], nil
}

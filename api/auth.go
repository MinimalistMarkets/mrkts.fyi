package api

import (
	"net/http"
	"os"
	"strings"
)

func verifyToken(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	splitToken := strings.Split(authHeader, "Bearer")
	if len(splitToken) != 2 {
		return false
	}

	bearerToken := strings.TrimSpace(splitToken[1])
	return bearerToken == os.Getenv("TOKEN")
}

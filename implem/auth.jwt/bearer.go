package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Depado/vuemonit/interactor"
	"github.com/dgrijalva/jwt-go"
)

// Convenience function that extracts the JWT from the incoming request's header
func (j jwtProvider) ExtractBearerToken(r *http.Request) (string, error) {
	var raw string

	if h := r.Header.Get("Authorization"); len(h) > 7 && strings.EqualFold(h[0:7], "BEARER ") {
		raw = h[7:]
	}
	if raw == "" {
		return raw, interactor.ErrBearerTokenNotFound
	}

	return raw, nil
}

func (j jwtProvider) ValidateBearerToken(r *http.Request) (*jwt.StandardClaims, error) {
	t, err := j.ExtractBearerToken(r)
	if err != nil {
		return nil, fmt.Errorf("extract bearer token: %w", err)
	}
	claims, err := j.CheckToken(t)
	if err != nil {
		return nil, fmt.Errorf("check token: %w", err)
	}

	return claims, nil
}

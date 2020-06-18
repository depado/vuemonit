package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/Depado/vuemonit/cmd"
	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
)

type jwtProvider struct {
	secret []byte
}

// NewJWTAuthProvider will create a new simple JWT authorization provider
func NewJWTAuthProvider(conf *cmd.Conf) interactor.AuthProvider {
	return &jwtProvider{secret: []byte(conf.Server.JWT.Secret)}
}

func (j jwtProvider) GenerateTokenPair(user *models.User) (string, string, error) {
	// Create the basic claims using the standard claims because we don't need
	// anything else
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		Subject:   user.ID,
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
	}

	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
	if err != nil {
		return "", "", fmt.Errorf("signing access token: %w", err)
	}

	claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Subject:   user.ID,
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
	}

	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
	if err != nil {
		return "", "", fmt.Errorf("signing refresh token: %w", err)
	}

	return access, refresh, nil
}

// CheckJWT will check whether or not the JWT is valid and return its claims if
// so
func (j jwtProvider) Check(token string) (jwt.StandardClaims, error) {
	claims := jwt.StandardClaims{}

	tkn, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return claims, fmt.Errorf("%v: %w", err, interactor.ErrJWT)
	}
	if !tkn.Valid {
		return claims, interactor.ErrInvalidJWT
	}

	return claims, nil
}

// Convenience function that extracts the JWT from the incoming request's header
func (j jwtProvider) Extract(r *http.Request) (string, error) {
	var raw string

	if h := r.Header.Get("Authorization"); len(h) > 7 && strings.EqualFold(h[0:7], "BEARER ") {
		raw = h[7:]
	}
	if raw == "" {
		return raw, interactor.ErrBearerTokenNotFound
	}

	return raw, nil
}

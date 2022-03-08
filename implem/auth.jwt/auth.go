package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/Depado/vuemonit/cmd"
	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
)

const cookieName = "tok"

type jwtProvider struct {
	secret []byte
	https  bool
	domain string
}

// NewJWTAuthProvider will create a new simple JWT authorization provider
func NewJWTAuthProvider(conf *cmd.Conf) interactor.AuthProvider {
	return &jwtProvider{
		secret: []byte(conf.Server.JWT.Secret),
		https:  conf.Server.Cookie.HTTPS,
		domain: conf.Server.Cookie.Domain,
	}
}

func (j jwtProvider) GenerateTokenPair(user *models.User) (*models.TokenPair, error) {
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
		return nil, fmt.Errorf("signing access token: %w", err)
	}

	claims = &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		Subject:   user.ID,
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
	}

	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(j.secret)
	if err != nil {
		return nil, fmt.Errorf("signing refresh token: %w", err)
	}

	return &models.TokenPair{Access: access, Refresh: refresh}, nil
}

// CheckJWT will check whether or not the JWT is valid and return its claims if
// so
func (j jwtProvider) CheckToken(token string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
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

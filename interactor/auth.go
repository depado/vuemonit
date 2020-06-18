package interactor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Depado/vuemonit/models"
	"github.com/dgrijalva/jwt-go"
)

// Login allows users to login and will return an access token on successful
// login
func (i Interactor) Login(email, password string) (string, error) {
	var u *models.User
	var err error

	if u, err = i.Store.GetUserByEmail(email); err != nil {
		return "", fmt.Errorf("get user by email: %w", err)
	}

	if !u.CheckPassword(password) {
		return "", ErrInvalidCredentials
	}

	u.LastLogin = time.Now()
	if err = i.Store.SaveUser(u); err != nil {
		return "", fmt.Errorf("save user: %w", err)
	}

	return i.Auth.GenerateJWT(u)
}

func (i Interactor) Register(email, password string) error {
	u, err := models.NewUser(email, password)
	if err != nil {
		return fmt.Errorf("new user: %w", err)
	}
	if err = i.Store.SaveUser(u); err != nil {
		return fmt.Errorf("save new user: %w", err)
	}
	return nil
}

// AuthCheck extracts the JWT token and the associated user from the
// incoming HTTP request
func (i Interactor) AuthCheck(r *http.Request) (*models.User, error) {
	var err error
	var token string
	var user *models.User
	var claims jwt.StandardClaims

	if token, err = i.Auth.Extract(r); err != nil {
		return nil, fmt.Errorf("extract token: %w", err)
	}

	if claims, err = i.Auth.Check(token); err != nil {
		return nil, fmt.Errorf("check token: %w", err)
	}

	if user, err = i.Store.GetUserByID(claims.Subject); err != nil {
		return nil, fmt.Errorf("retrieve user by id: %w", err)
	}

	return user, nil
}

package interactor

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Depado/vuemonit/models"
	"github.com/golang-jwt/jwt/v4"
)

// Login allows users to login and will return an access token on successful
// login
func (i Interactor) Login(email, password string) (*models.TokenPair, *http.Cookie, error) {
	var u *models.User
	var err error

	if u, err = i.Store.GetUserByEmail(email); err != nil {
		return nil, nil, fmt.Errorf("get user by email: %w", err)
	}

	if !u.CheckPassword(password) {
		return nil, nil, ErrInvalidCredentials
	}

	u.LastLogin = time.Now()
	if err = i.Store.SaveUser(u); err != nil {
		return nil, nil, fmt.Errorf("save user: %w", err)
	}

	tp, err := i.Auth.GenerateTokenPair(u)
	if err != nil {
		return nil, nil, fmt.Errorf("generate token pair: %w", err)
	}

	c, err := i.Auth.GenerateCookie(u, tp)
	if err != nil {
		return nil, nil, fmt.Errorf("generate cookie: %w", err)
	}

	return tp, c, nil
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
func (i Interactor) AuthCheck(w http.ResponseWriter, r *http.Request) (*models.User, error) {
	var err error
	var user *models.User
	var issue bool
	var claims *jwt.RegisteredClaims

	// Request has bearer auth, process and return
	if claims, err = i.Auth.ValidateBearerToken(r); err != nil {
		if claims, issue, err = i.Auth.ValidateCookie(r); err != nil {
			return nil, fmt.Errorf("no valid cookie or bearer found: %w", err)
		}
	}

	if user, err = i.Store.GetUserByID(claims.Subject); err != nil {
		return nil, fmt.Errorf("retrieve user by id: %w", err)
	}

	if issue {
		// Perform additional checks on user to see if the session must continue
		c, err := i.Auth.GenerateCookie(user, nil)
		if err != nil {
			return nil, fmt.Errorf("generate cookie: %w", err)
		}
		http.SetCookie(w, c)
	}

	return user, nil
}

func (i Interactor) Refresh(token string) (*models.TokenPair, error) {
	claims, err := i.Auth.CheckToken(token)
	if err != nil {
		return nil, fmt.Errorf("check token: %w", err)
	}

	user, err := i.Store.GetUserByID(claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("retrieve user by id: %w", err)
	}

	return i.Auth.GenerateTokenPair(user)
}

func (i Interactor) Logout() *http.Cookie {
	return i.Auth.DropAccessCookie()
}

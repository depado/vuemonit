package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
)

func (j jwtProvider) GenerateCookie(u *models.User, tp *models.TokenPair) (*http.Cookie, error) {
	var err error
	if tp == nil {
		if tp, err = j.GenerateTokenPair(u); err != nil {
			return nil, fmt.Errorf("generate token pair: %w", err)
		}
	}
	v, err := tp.Marhsal()
	if err != nil {
		return nil, fmt.Errorf("marshal token pair: %w", err)
	}
	return &http.Cookie{
		Name:     cookieName,
		Value:    base64.StdEncoding.EncodeToString(v),
		Path:     "/",
		HttpOnly: true,
		Secure:   j.https,
		Domain:   j.domain,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	}, nil
}

func (j jwtProvider) ValidateCookie(r *http.Request) (*jwt.StandardClaims, bool, error) {
	tp, err := j.extractCookieTokens(r)
	if err != nil {
		return nil, false, fmt.Errorf("extract cookie token: %w", err)
	}

	// Access token is still valid, return as soon as possible
	claims, err := j.CheckToken(tp.Access)
	if err == nil {
		return claims, false, nil
	}

	claims, err = j.CheckToken(tp.Refresh)
	if err != nil {
		return nil, false, fmt.Errorf("check refresh: %w", err)
	}

	return claims, true, nil
}

func (j jwtProvider) extractCookieTokens(r *http.Request) (*models.TokenPair, error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return nil, fmt.Errorf("cookie not found: %v: %w", err, interactor.ErrCookieNotFound)
	}
	raw, err := base64.StdEncoding.DecodeString(c.Value)
	if err != nil {
		return nil, fmt.Errorf("unable to decode base64: %w", err)
	}
	tp := &models.TokenPair{}
	if err := json.Unmarshal(raw, tp); err != nil {
		return nil, fmt.Errorf("%v: %w", err, interactor.ErrCookieFormat)
	}
	return tp, nil
}

func (j jwtProvider) DropAccessCookie() *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Path:     "/",
		Value:    "",
		HttpOnly: true,
		Secure:   j.https,
		Domain:   j.domain,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	}
}

package interactor

import "errors"

// ErrInvalidJWT is returned when the JWT isn't valid (expired, unsigned, etc.)
var ErrInvalidJWT = errors.New("invalid jwt")

// ErrBearerTokenNotFound is returned when the bearer token can't be found in
// the http request
var ErrBearerTokenNotFound = errors.New("token not found")

// ErrJWT is a JWT error
var ErrJWT = errors.New("jwt error")

// ErrInvalidCredentials is returned when invalid credentials are provided
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrNotFound is returned when an element isn't found
var ErrNotFound = errors.New("not found")

// ErrPermission is returned when a permission error occurs
var ErrPermission = errors.New("permission error")

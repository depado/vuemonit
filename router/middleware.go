package router

import (
	"errors"
	"net/http"

	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
	"github.com/gin-gonic/gin"
)

// AuthRequired is a middleware to simply check if a user is logged-in or not
// This middleware will check both the cookie and the Authorization header,
// allowing users to provide their authentication tokens the way they like
// In the case of a cookie based authentication, if the access token is expired
// this middleware will silently refresh the tokens as to avoid unnecessary
// calls
func (r Router) AuthRequired() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		var err error
		var user *models.User

		if user, err = r.lh.AuthCheck(c.Writer, c.Request); err != nil {
			r.log.Debug().Err(err).Msg("auth check failed")
			switch {
			case errors.Is(err, interactor.ErrBearerTokenNotFound):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "missing authorization header or wrong token format",
				})
			case errors.Is(err, interactor.ErrCookieFormat) || errors.Is(err, interactor.ErrCookieNotFound):
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "no cookie found or invalid cookie",
				})
			case errors.Is(err, interactor.ErrInvalidJWT) || errors.Is(err, interactor.ErrJWT):
				c.AbortWithStatus(http.StatusUnauthorized)
			default:
				c.AbortWithStatus(http.StatusInternalServerError)
			}
			return
		}

		c.Set("user", user)
		c.Next()
	})
}

// getUserFromContext is a simple helper function that retrieves a user from
// gin's context if any
func getUserFromContext(c *gin.Context) (*models.User, error) {
	if v, ok := c.Get("user"); ok {
		if user, ok := v.(*models.User); ok {
			return user, nil
		}
	}
	return nil, errors.New("tried to access protected route without a user")
}

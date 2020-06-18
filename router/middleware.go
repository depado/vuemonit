package router

import (
	"errors"
	"net/http"

	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
	"github.com/gin-gonic/gin"
)

// AuthRequired is a middleware to simply check if a user is logged-in or not
// Given a valid JWT in the Authorization header, this will also fetch the user
// from the database
func (r Router) AuthRequired() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		var err error
		var user *models.User

		if user, err = r.lh.AuthCheck(c.Request); err != nil {
			switch {
			case errors.Is(err, interactor.ErrBearerTokenNotFound):
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "missing authorization header or wrong token format",
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

func getUserFromContext(c *gin.Context) (*models.User, error) {
	if v, ok := c.Get("user"); ok {
		if user, ok := v.(*models.User); ok {
			return user, nil
		}
	}
	return nil, errors.New("tried to access protected route without a user")
}

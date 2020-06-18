package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r Router) Me(c *gin.Context) {
	clog := r.log.With().Str("route", "/api/v1/me").Str("method", "GET").Logger()
	u, err := getUserFromContext(c)
	if err != nil {
		clog.Err(err).Msg("tried to access auth required route without a user")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, r.lh.FormatSelf(u))
}

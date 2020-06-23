package router

import (
	"errors"
	"net/http"

	"github.com/Depado/vuemonit/interactor"
	"github.com/asdine/storm/v3"
	"github.com/gin-gonic/gin"
)

func (r Router) GetTimedResponses(c *gin.Context) {
	clog := r.log.With().Str("route", "/api/v1/service").Str("method", "GET").Logger()
	id := c.Param("id")
	if id == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	u, err := getUserFromContext(c)
	if err != nil {
		clog.Err(err).Msg("tried to access auth required route without a user")
		c.Status(http.StatusInternalServerError)
		return
	}

	tr, err := r.lh.GetTimedResponsesByServiceID(u, id)
	if err != nil {
		clog.Debug().Err(err).Send()
		if errors.Is(err, interactor.ErrNotFound) || errors.Is(err, storm.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "service not found"})
		} else if errors.Is(err, interactor.ErrPermission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "permission error"})
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, r.lh.FormatTimedResponses(tr))
}

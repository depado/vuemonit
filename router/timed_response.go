package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Depado/vuemonit/interactor"
	"github.com/Depado/vuemonit/models"
	"github.com/asdine/storm/v3"
	"github.com/gin-gonic/gin"
)

func parseRange(from, to string) (time.Time, time.Time, error) {
	var f, t time.Time
	if (from != "" && to == "") || (from == "" && to != "") {
		return f, t, fmt.Errorf("both 'from' and 'to' parameters are required or none of them should be set")
	}
	i, err := strconv.ParseInt(from, 10, 64)
	if err != nil {
		return f, t, fmt.Errorf("'from' parameter isn't a valid timestamp")
	}
	f = time.Unix(i, 0)
	i, err = strconv.ParseInt(to, 10, 64)
	if err != nil {
		return f, t, fmt.Errorf("'to' parameter isn't a valid timestamp")
	}
	t = time.Unix(i, 0)
	return f, t, nil
}

// GetTimedResponses is the endpoint in charge of retrieving timed responses
func (r Router) GetTimedResponses(c *gin.Context) {
	var hasrange bool
	var fromt, tot time.Time
	var tr []*models.TimedResponse
	var err error

	id := c.Param("id")
	if id == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	clog := r.log.With().Str("route", fmt.Sprintf("/api/v1/service/%s/tr", id)).Str("method", "GET").Logger()

	from := c.Query("from")
	to := c.Query("to")
	if from != "" && to != "" {
		if fromt, tot, err = parseRange(from, to); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		hasrange = true
	}

	u, err := getUserFromContext(c)
	if err != nil {
		clog.Err(err).Msg("tried to access auth required route without a user")
		c.Status(http.StatusInternalServerError)
		return
	}
	if hasrange {
		tr, err = r.lh.GetTimedResponseRange(u, id, fromt, tot)
	} else {
		tr, err = r.lh.GetTimedResponsesByServiceID(u, id, 0, true)
	}
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

func (r Router) TailTimedResponses(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	clog := r.log.With().Str("route", fmt.Sprintf("/api/v1/service/%s/tr/tail", id)).Str("method", "GET").Logger()

	u, err := getUserFromContext(c)
	if err != nil {
		clog.Err(err).Msg("tried to access auth required route without a user")
		c.Status(http.StatusInternalServerError)
		return
	}

	tr, err := r.lh.GetTimedResponsesByServiceID(u, id, 100, true)
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

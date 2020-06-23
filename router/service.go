package router

import (
	"errors"
	"net/http"
	"time"

	"github.com/Depado/vuemonit/interactor"
	"github.com/asdine/storm/v3"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ServiceQuery struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HealthCheck `json:"healthcheck"`
}

type HealthCheck struct {
	URL   string        `json:"url"`
	Every time.Duration `json:"every"`
}

func (s ServiceQuery) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Description, validation.Required),
		validation.Field(&s.HealthCheck.URL, validation.Required, is.URL),
	)
}

func (r Router) PostService(c *gin.Context) {
	clog := r.log.With().Str("route", "/api/v1/service").Str("method", "POST").Logger()
	var err error
	var sq ServiceQuery
	u, err := getUserFromContext(c)
	if err != nil {
		clog.Err(err).Msg("tried to access auth required route without a user")
		c.Status(http.StatusInternalServerError)
		return
	}

	if err = c.ShouldBind(&sq); err != nil {
		clog.Debug().Err(err).Msg("unable to bind")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err = sq.Validate(); err != nil {
		clog.Debug().Err(err).Msg("validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s, err := r.lh.NewService(u, sq.Name, sq.Description, sq.HealthCheck.URL)
	if err != nil {
		clog.Err(err).Msg("unable to create new service")
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, r.lh.FormatService(s))
}

func (r Router) GetService(c *gin.Context) {
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

	s, err := r.lh.GetServiceByID(u, id)
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

	c.JSON(http.StatusOK, r.lh.FormatService(s))
}

func (r Router) GetServices(c *gin.Context) {
	clog := r.log.With().Str("route", "/api/v1/services").Str("method", "GET").Logger()

	u, err := getUserFromContext(c)
	if err != nil {
		clog.Err(err).Msg("tried to access auth required route without a user")
		c.Status(http.StatusInternalServerError)
		return
	}

	svx, err := r.lh.GetServices(u)
	if err != nil {
		clog.Debug().Err(err).Send()
		if errors.Is(err, interactor.ErrNotFound) || errors.Is(err, storm.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "services not found"})
		} else {
			c.Status(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, r.lh.FormatServices(svx))
}

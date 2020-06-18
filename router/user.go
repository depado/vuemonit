package router

import (
	"errors"
	"net/http"

	"github.com/Depado/vuemonit/interactor"
	"github.com/asdine/storm/v3"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CredentialQuery struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c CredentialQuery) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Length(12, 50)),
	)
}

func (r Router) Register(c *gin.Context) {
	clog := r.log.With().Str("route", "/api/v1/register").Str("method", "POST").Logger()
	var err error
	var cq CredentialQuery

	if err = c.ShouldBind(&cq); err != nil {
		clog.Debug().Err(err).Msg("unable to bind")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err = cq.Validate(); err != nil {
		clog.Debug().Err(err).Msg("validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := r.lh.Register(cq.Email, cq.Password); err != nil {
		if errors.Is(err, storm.ErrAlreadyExists) {
			c.JSON(http.StatusForbidden, gin.H{"error": "account already exists"})
			return
		}
		clog.Err(err).Msg("error occured during register")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}

func (r Router) Login(c *gin.Context) {
	clog := r.log.With().Str("route", "/api/v1/login").Str("method", "POST").Logger()
	var err error
	var cq CredentialQuery

	if err = c.ShouldBind(&cq); err != nil {
		clog.Debug().Err(err).Msg("unable to bind")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err = cq.Validate(); err != nil {
		clog.Debug().Err(err).Msg("validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tok, err := r.lh.Login(cq.Email, cq.Password)
	if err != nil {
		if errors.Is(err, interactor.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
			return
		}
		clog.Err(err).Str("email", cq.Email).Msg("unable to login")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tok})
}

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

package router

import (
	"errors"
	"net/http"

	"github.com/Depado/vuemonit/interactor"
	"github.com/asdine/storm/v3"
	"github.com/gin-gonic/gin"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Register will allow a new user to register
func (r Router) Register(c *gin.Context) {
	clog := r.log.With().Str("route", "/api/v1/register").Str("method", "POST").Logger()
	var err error
	var cq CredentialQuery

	if !r.register {
		c.JSON(http.StatusForbidden, gin.H{"error": "registering isn't allowed on this instance"})
		return
	}

	if err = c.ShouldBind(&cq); err != nil {
		clog.Debug().Err(err).Msg("unable to bind")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err = cq.ValidateRegister(); err != nil {
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

// Login is the route used to login a user. This will generate both a short
// lived access token and long lived refresh token
// This endpoint will send back both a cookie and a JSON containing those tokens
// and thus can be used to login both on the frontend app and a standard http
// client
func (r Router) Login(c *gin.Context) {
	clog := r.log.With().Str("route", "/api/v1/login").Str("method", "POST").Logger()
	var err error
	var cq CredentialQuery

	if err = c.ShouldBind(&cq); err != nil {
		clog.Debug().Err(err).Msg("unable to bind")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err = cq.ValidateLogin(); err != nil {
		clog.Debug().Err(err).Msg("validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tp, cookie, err := r.lh.Login(cq.Email, cq.Password)
	if err != nil {
		if errors.Is(err, interactor.ErrInvalidCredentials) || errors.Is(err, interactor.ErrNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong credentials"})
			return
		}
		clog.Err(err).Str("email", cq.Email).Msg("unable to login")
		c.Status(http.StatusInternalServerError)
		return
	}
	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, tp)
}

// Refresh is the route used to refresh the access token, given a valid refresh
// token. This route should not be used by the frontend application since the
// cookie based authentication does silent refreshes when needed
func (r Router) Refresh(c *gin.Context) {
	clog := r.log.With().Str("route", "/api/v1/refresh").Str("method", "POST").Logger()
	var err error
	var rq RefreshQuery

	if err = c.ShouldBind(&rq); err != nil {
		clog.Debug().Err(err).Msg("unable to bind")
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	if err = rq.Validate(); err != nil {
		clog.Debug().Err(err).Msg("validation failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tp, err := r.lh.Refresh(rq.RefreshToken)
	if err != nil {
		clog.Debug().Err(err).Msg("unable to refresh")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, tp)
}

// Logout will return an invalidated cookie, effectively disabling the cookie
// that was already present in the browser, if any
func (r Router) Logout(c *gin.Context) {
	http.SetCookie(c.Writer, r.lh.Logout())
}

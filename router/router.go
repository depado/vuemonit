package router

import (
	"net/http"
	"path"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/Depado/vuemonit/cmd"
	"github.com/Depado/vuemonit/interactor"
	"github.com/rs/zerolog"
)

type front struct {
	serve bool
	path  string
}

type Router struct {
	e        *gin.Engine
	log      *zerolog.Logger
	front    front
	lh       interactor.LogicHandler
	register bool
}

func NewRouter(c *cmd.Conf, e *gin.Engine, log *zerolog.Logger, lh interactor.LogicHandler) *Router {
	r := &Router{e: e, log: log,
		front: front{
			serve: c.Front.Serve,
			path:  c.Front.Path,
		},
		lh:       lh,
		register: c.Register,
	}
	if c.Server.Log {
		r.e.Use(gin.Logger())
	}
	return r
}

// SetRoutes will setup the various served routes
func (r Router) SetRoutes() {
	if r.front.serve {
		r.e.Use(static.Serve("/", static.LocalFile(r.front.path, true)))
		r.e.LoadHTMLGlob(path.Join(r.front.path, "index.html"))
	}

	g := r.e.Group("/api/v1")
	{
		// Simple health check route
		g.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

		// Auth related routes
		g.POST("/register", r.Register)
		g.POST("/login", r.Login)
		g.POST("/refresh", r.Refresh)
		g.GET("/logout", r.Logout)

		// Data related routes
		g.GET("/me", r.AuthRequired(), r.Me)
		g.POST("/service", r.AuthRequired(), r.PostService)
		g.GET("/services", r.AuthRequired(), r.GetServices)
		g.GET("/service/:id", r.AuthRequired(), r.GetService)
		g.GET("/service/:id/tr", r.AuthRequired(), r.GetTimedResponses)
	}
}

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
	e     *gin.Engine
	log   *zerolog.Logger
	front front
	lh    interactor.LogicHandler
}

func NewRouter(c *cmd.Conf, e *gin.Engine, log *zerolog.Logger, lh interactor.LogicHandler) *Router {
	r := &Router{e: e, log: log,
		front: front{
			serve: c.Front.Serve,
			path:  c.Front.Path,
		},
		lh: lh,
	}
	if c.Server.Log {
		r.e.Use(gin.Logger())
	}
	return r
}

func (r Router) SetRoutes() {
	if r.front.serve {
		r.e.Use(static.Serve("/", static.LocalFile(r.front.path, true)))
		r.e.LoadHTMLGlob(path.Join(r.front.path, "index.html"))
	}

	g := r.e.Group("/api/v1")
	{
		g.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })
		g.POST("/register", r.Register)
		g.POST("/login", r.Login)
		g.GET("/me", r.AuthRequired(), r.Me)
		g.POST("/service", r.AuthRequired(), r.PostService)
		g.GET("/service/:id", r.AuthRequired(), r.GetService)
	}
}

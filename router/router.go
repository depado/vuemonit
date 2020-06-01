package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Depado/vuemonit/cmd"
)

type Router struct {
	e   *gin.Engine
	log *zap.Logger
}

func NewRouter(c *cmd.Conf, e *gin.Engine, log *zap.Logger) *Router {
	r := &Router{
		e:   e,
		log: log,
	}
	if c.Server.Log {
		r.e.Use(gin.Logger())
	}
	return r
}

func (r Router) SetRoutes() {
	g := r.e.Group("/api/v1")
	g.GET("/health", func(c *gin.Context) {
		r.log.Info("I'm there")
	})
}

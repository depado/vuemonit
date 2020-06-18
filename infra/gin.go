package infra

import (
	"github.com/Depado/ginprom"
	"github.com/Depado/vuemonit/cmd"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewServer(conf *cmd.Conf, corsc *cors.Config) *gin.Engine {
	setMode(conf.Server.Mode)

	r := gin.New()
	if !conf.Prometheus.Disabled {
		p := ginprom.New(
			ginprom.Subsystem(conf.Prometheus.Prefix),
			ginprom.Engine(r),
		)
		r.Use(p.Instrument())
	}

	r.Use(gin.Recovery())

	if corsc != nil {
		r.Use(cors.New(*corsc))
	}

	return r
}

func setMode(mode string) {
	switch mode {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	case "release":
		gin.SetMode(gin.ReleaseMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}
}

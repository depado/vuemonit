package router

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/Depado/vuemonit/cmd"
)

func Run(r *Router, c *cmd.Conf) {
	r.SetRoutes()
	r.log.Info("starting", zap.String("host", c.Server.Host), zap.Int("port", c.Server.Port))
	if err := r.e.Run(fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)); err != nil {
		r.log.Fatal("couldn't start server", zap.Error(err))
	}
}

package router

import (
	"fmt"

	"github.com/Depado/vuemonit/cmd"
)

func Run(r *Router, c *cmd.Conf) {
	r.SetRoutes()
	r.log.Info().Str("host", c.Server.Host).Int("port", c.Server.Port).Msg("starting server")
	if err := r.e.Run(fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)); err != nil {
		r.log.Err(err).Msg("unable to start server")
	}
}

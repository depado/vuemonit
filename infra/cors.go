package infra

import (
	"fmt"
	"time"

	"github.com/Depado/vuemonit/cmd"
	"github.com/gin-contrib/cors"
)

// NewCORS will return a new CORS configuration if given a valid configuration
func NewCORS(conf *cmd.Conf) (*cors.Config, error) {
	cconf := conf.Server.CORS
	if !cconf.Enable {
		return nil, nil
	}
	c := &cors.Config{
		AllowCredentials: true,
		MaxAge:           50 * time.Second,
		AllowMethods:     cconf.Methods,
		AllowHeaders:     cconf.Headers,
		ExposeHeaders:    cconf.Expose,
	}

	switch {
	case len(cconf.Origins) > 0:
		c.AllowOrigins = cconf.Origins
	case cconf.All:
		c.AllowAllOrigins = true
	default:
		return nil, fmt.Errorf("allow all origins disabled but no allowed origins provided")
	}

	return c, nil
}

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"

	_ "github.com/Depado/vuemonit/statik"
)

func (r Router) SetupEmbeddedFront() {
	// Open the statik box
	sfs, err := fs.New()
	if err != nil {
		r.log.Fatal().Err(err).Msg("unable to initialize static files")
	}

	// Setup the fileserver
	fs := http.FileServer(sfs)
	r.e.Use(func(c *gin.Context) {
		// Special case for /
		if c.Request.URL.Path == "/" {
			c.FileFromFS("/main.html", sfs) // Forced to rename index.html to main.html
			fs.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
		// In case there is a static file matching
		if _, err := sfs.Open(c.Request.URL.Path); err == nil {
			fs.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})

	// Adds a handler that will default to our main index since we want the
	// frontend to determine its own page
	r.e.NoRoute(func(c *gin.Context) {
		c.FileFromFS("/main.html", sfs)
	})
}

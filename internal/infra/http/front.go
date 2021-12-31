package http

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

// Frontend represents HTTP API for frontend related assets.
type Frontend struct{}

// NewFrontend creates a new frontend provider instance.
func NewFrontend() *Frontend {
	return &Frontend{}
}

// SetupRoutes creates a frontend API instance in order to serve frontend related files.
func (f *Frontend) SetupRoutes(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./web/dist", true)))

	// serve all unknown routes with the index file to support browser-level routing
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})
}

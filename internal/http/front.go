package http

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

type Frontend struct {
}

func NewFrontend() *Frontend {
	return &Frontend{}
}

func (f *Frontend) SetupRoutes(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("./web/dist", true)))

	// serve all unknown routes with the index file to support browser-level routing
	r.NoRoute(func(c *gin.Context) {
		c.File("./web/dist/index.html")
	})
}

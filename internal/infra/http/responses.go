package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type httpErr struct {
	Error string `json:"error"`
}

func success(c *gin.Context, h interface{}) {
	c.JSON(http.StatusOK, h)
}

func badRequestError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, httpErr{
		Error: err.Error(),
	})
}

func internalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, httpErr{
		Error: err.Error(),
	})
}

func unauthorizedError(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, httpErr{
		Error: err.Error(),
	})
}

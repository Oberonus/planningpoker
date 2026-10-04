package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	logrus.Errorf("HTTP request failed: %v", err)

	c.JSON(http.StatusInternalServerError, httpErr{
		Error: "internal server error",
	})
}

func unauthorizedError(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, httpErr{
		Error: err.Error(),
	})
}

package http

import (
	"errors"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *API) withUser(cb func(*gin.Context, string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			unauthorizedError(c, errors.New("unauthorized"))
			return
		}

		userID, err := h.authenticator.AuthenticateByToken(parts[1])
		if err != nil {
			logrus.Infof("user auth failed: %v", err)
			unauthorizedError(c, errors.New("unauthorized"))
			return
		}

		cb(c, userID)
	}
}

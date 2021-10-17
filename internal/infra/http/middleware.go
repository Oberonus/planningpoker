package http

import (
	"errors"
	"planningpoker/internal/domain/users"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *API) withUser(cb func(*gin.Context, *users.User)) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		parts := strings.Split(auth, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			unauthorizedError(c, errors.New("unauthorized"))
			return
		}

		cmd, err := users.NewAuthByIDCommand(parts[1])
		if err != nil {
			badRequestError(c, err)
			return
		}

		user, err := h.usersService.AuthenticateByID(*cmd)
		if err != nil {
			unauthorizedError(c, err)
			return
		}

		if user == nil {
			unauthorizedError(c, errors.New("user not found"))
			return
		}

		cb(c, user)
	}
}

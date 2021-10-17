package http

import (
	"planningpoker/internal/domain/users"

	"github.com/gin-gonic/gin"
)

func (h *API) Register(c *gin.Context) {
	pl := struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}{}
	if err := c.BindJSON(&pl); err != nil {
		badRequestError(c, err)
		return
	}

	cmd, err := users.NewRegisterCommand(pl.Name)
	if err != nil {
		badRequestError(c, err)
		return
	}

	user, err := h.usersService.Register(*cmd)
	if err != nil {
		badRequestError(c, err)
		return
	}

	success(c, gin.H{
		"user_id": user.ID(),
	})
}

func (h *API) CurrentUser(c *gin.Context, user *users.User) {
	success(c, gin.H{
		"name": user.Name(),
	})
}

func (h *API) ChangeUserData(c *gin.Context, user *users.User) {
	pl := struct {
		Name string `json:"name"`
	}{}
	if err := c.BindJSON(&pl); err != nil {
		badRequestError(c, err)
		return
	}

	cmd, err := users.NewUpdateCommand(user.ID(), pl.Name)
	if err != nil {
		badRequestError(c, err)
		return
	}

	_, err = h.usersService.Update(*cmd)
	if err != nil {
		internalError(c, err)
		return
	}

	success(c, gin.H{})
}

package http

import (
	"errors"
	"planningpoker/internal/domain/users"

	"github.com/gin-gonic/gin"
)

func (h *API) register(c *gin.Context) {
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

func (h *API) currentUser(c *gin.Context, userID string) {
	user, err := h.usersService.Get(userID)
	if err != nil || user == nil {
		badRequestError(c, err)
		return
	}

	if user == nil {
		badRequestError(c, errors.New("user not found"))
		return
	}

	success(c, gin.H{
		"name": user.Name(),
	})
}

func (h *API) changeUserData(c *gin.Context, userID string) {
	pl := struct {
		Name string `json:"name"`
	}{}
	if err := c.BindJSON(&pl); err != nil {
		badRequestError(c, err)
		return
	}

	cmd, err := users.NewUpdateCommand(userID, pl.Name)
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

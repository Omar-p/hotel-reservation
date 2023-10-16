package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/omar-p/hotel-reservation/db"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(store db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: store,
	}
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, _ := h.userStore.GetUsers(context.Background())
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id", "")
	userById, err := h.userStore.GetUserByID(context.Background(), id)
	if err != nil {
		return err
	}
	return c.JSON(userById)
}

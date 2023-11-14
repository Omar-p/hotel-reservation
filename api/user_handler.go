package api

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/omar-p/hotel-reservation/db"
	"github.com/omar-p/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
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
	users, _ := h.userStore.GetUsers(c.Context())
	return c.JSON(users)
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id", "")
	userById, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(userById)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var reqBody types.CreateUserRequest
	if err := c.BodyParser(&reqBody); err != nil {
		return err
	}
	if errors := reqBody.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	user, err := types.NewUserFromCreateRequest(&reqBody)
	if err != nil {
		return err
	}
	user, err = h.userStore.InsertUser(context.Background(), user)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandlePutUser(ctx *fiber.Ctx) error {
	var (
		reqBody types.UpdateUserRequest
		userID  = ctx.Params("id")
	)
	if err := ctx.BodyParser(&reqBody); err != nil {
		return err
	}
	if errors := reqBody.Validate(); len(errors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}
	if err := h.userStore.UpdateUser(ctx.Context(), bson.M{"_id": userID}, &reqBody); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")
	if err := h.userStore.DeleteUser(ctx.Context(), userID); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/omar-p/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	return c.SendString("Get users")
}

func HandleGetUser(c *fiber.Ctx) error {
	path, _ := c.ParamsInt("id", 0)
	return c.JSON(types.User{
		ID:        path,
		FirstName: "Omar",
		LastName:  "Shabaan",
	})

}

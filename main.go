package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{"msg": "Hello, World 👋!"})
	})
	app.Listen(":8080")
}
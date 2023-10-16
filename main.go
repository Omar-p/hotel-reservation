package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/omar-p/hotel-reservation/api"
)

func main() {
	listenAddr := flag.String("listenAdder", ":8080", "port that the server will listen on")
	flag.Parse()

	app := fiber.New()
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/users", api.HandleGetUsers)
	apiV1.Get("/users/:id", api.HandleGetUser)

	app.Listen(*listenAddr)
}

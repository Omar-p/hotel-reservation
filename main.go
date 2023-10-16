package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/omar-p/hotel-reservation/api"
	"github.com/omar-p/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var (
	dbUri = os.Getenv("DB_URI")
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	},
}

func main() {
	listenAddr := flag.String("listenAdder", ":8080", "port that the server will listen on")
	flag.Parse()

	ctx := context.Background()
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	// handlers initialization
	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(config)
	apiV1 := app.Group("/api/v1")

	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Get("/users/:id", userHandler.HandleGetUser)

	app.Listen(*listenAddr)
}

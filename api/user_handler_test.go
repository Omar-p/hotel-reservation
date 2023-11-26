package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/omar-p/hotel-reservation/db"
	"github.com/omar-p/hotel-reservation/types"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http/httptest"
	"testing"
)

var (
	dbName = "Hotel-Reservation-Test"
)

type testDB struct {
	db.UserStore
}

func setup(t *testing.T) (*testDB, func()) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:7.0.3"))
	if err != nil {
		panic(err)
	}

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		panic(err)
	}
	if err = mongoClient.Ping(ctx, nil); err != nil {
		log.Println("ping")
		log.Println(err)
	}

	return &testDB{UserStore: db.NewMongoUserStore(mongoClient, dbName)},
		func() {
			fmt.Println("-- Clean up the container --")

			// Teardown: Stop and remove the MongoDB container
			if err := mongodbContainer.Terminate(ctx); err != nil {
				t.Fatal("Could not stop MongoDB container:", err)
			}
		}
}

func TestPostUser(t *testing.T) {
	tdb, teardown := setup(t)
	defer teardown()

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserRequest{
		FirstName: "omar",
		LastName:  "shabaan",
		Email:     "omar@email.com",
		Password:  "123456789",
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	var u types.User
	err = json.NewDecoder(resp.Body).Decode(&u)
	if err != nil {
		t.Fatalf("error in decoding response body %v", err)
	}
	// check if the user is a valid mongo object id
	if u.ID.IsZero() {
		t.Fatalf("expected a valid object id, got %v", u.ID)
	}
	// check that password is not returned
	if u.EncryptedPassword != "" {
		t.Fatalf("expected empty password, got %s", u.EncryptedPassword)
	}

	if u.FirstName != params.FirstName {
		// use printf to print the value of the variable
		t.Fatalf("expected %s, got %s", params.FirstName, u.FirstName)
	}
	if u.LastName != params.LastName {
		t.Fatalf("expected %s, got %s", params.LastName, u.LastName)
	}
	if u.Email != params.Email {
		t.Fatalf("expected %s, got %s", params.Email, u.Email)
	}

}

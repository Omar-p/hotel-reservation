package db

import (
	"context"
	"fmt"
	"github.com/omar-p/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

var userCollection = os.Getenv("USER_COLLECTION")

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
	UpdateUser(context.Context, bson.M, *types.UpdateUserRequest) error
	DeleteUser(context.Context, string) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(c *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: c,
		coll:   c.Database(dbName).Collection(userCollection),
	}

}

func (m *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	result, err := m.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (m *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	users := make([]*types.User, 0)
	err = cur.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (m *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, updateRequest *types.UpdateUserRequest) error {
	update := bson.D{
		{"$set", updateRequest.ToBSON()},
	}
	result, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no documents matched the filter")
	}
	return nil
}

func (m *MongoUserStore) DeleteUser(ctx context.Context, s string) error {
	oid, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return err
	}
	result, err := m.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no documents matched the filter")
	}
	return nil
}

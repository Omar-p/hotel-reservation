package db

import (
	"context"
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

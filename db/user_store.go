package db

import (
	"context"
	"github.com/MarkAndzh/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

const userColl = "users"

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	PutUser(context.Context, string, types.UpdateUserParams) error
}

type MongoUserStore struct {
	client *mongo.Client
	dbname string
	coll   *mongo.Collection
	logger *slog.Logger
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	res, err := s.coll.DeleteOne(ctx, bson.D{{"_id", userId}})
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	if res.DeletedCount != 1 {
		s.logger.Error("Couldn't find user with id: ", id)
		return err
	}
	s.logger.Info("Successfully deleted user wit id: ", id)
	return nil
}

func (s *MongoUserStore) PutUser(ctx context.Context, id string, params types.UpdateUserParams) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	filter := bson.D{{"_id", oid}}
	values := params.ToBSON()
	update := bson.D{{"$set", values}}
	res, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		s.logger.Error("did not update record")
		return err
	}
	return nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err = cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}
	s.logger.Debug("got user with id: %v", id)
	return &user, nil
}

func NewMongoUserStore(logger *slog.Logger, c *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: c,
		coll:   c.Database(DBNAME).Collection(userColl),
		logger: logger,
	}
}

package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepo struct {
	UserColl *mongo.Collection
}

func NewMongoUserRepo(db *mongo.Client) *MongoUserRepo {
	return &MongoUserRepo{UserColl: db.Database(dbname).Collection(userCollection)}
}

func (r *MongoUserRepo) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.UserColl.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("Create user error: %v", err)
	}
	return nil
}

func (r *MongoUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := r.UserColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("GetUserByEmail error: %v", err)
	}
	return &user, nil
}

func (r *MongoUserRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	filter := bson.M{"login": login}
	err := r.UserColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("GetUserByLogin error: %v", err)
	}
	return &user, nil
}

func (r *MongoUserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	oid, err := convertToObjectIDs(id)
	if err != nil {
		return nil, fmt.Errorf("InvalidID: %v", err)
	}
	var user models.User
	filter := bson.M{"_id": oid[0]}
	err = r.UserColl.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID error: %v", err)
	}
	return &user, nil
}

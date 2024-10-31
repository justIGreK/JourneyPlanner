package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"

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
	return err
}

func (r *MongoUserRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := r.UserColl.FindOne(ctx, filter).Decode(&user)
	return &user, err
}

func (r *MongoUserRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	filter := bson.M{"login": login}
	err := r.UserColl.FindOne(ctx, filter).Decode(&user)
	return &user, err
}

func (r *MongoUserRepo) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	oid, err := convertToObjectIDs(id)
	if err != nil {
		return nil, errors.New("InvalidID")
	}
	var user models.User
	filter := bson.M{"_id": oid[0]}
	err = r.UserColl.FindOne(ctx, filter).Decode(&user)
	return &user, err
}

package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepo struct {
	UserColl *mongo.Collection
}

func NewUserTaskRepo(db *mongo.Client) *MongoUserRepo {
	return &MongoUserRepo{UserColl: db.Database(dbname).Collection(userCollection)}
}

func (r *MongoUserRepo) CreateUser(ctx context.Context, user models.User) error {
	_, err := r.UserColl.InsertOne(ctx, user)
	return err
}

func (r *MongoUserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	filter := bson.M{"email": email}
	err := r.UserColl.FindOne(ctx, filter).Decode(&user)
	return user, err
}

func (r *MongoUserRepo) GetUserByLogin(ctx context.Context, login string) (models.User, error) {
	var user models.User
	filter := bson.M{"login": login}
	err := r.UserColl.FindOne(ctx, filter).Decode(&user)
	return user, err
}

func (r *MongoUserRepo) GetUserByID(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	var user models.User
	filter := bson.M{"_id": id}
	err := r.UserColl.FindOne(ctx, filter).Decode(&user)
	return user, err
}

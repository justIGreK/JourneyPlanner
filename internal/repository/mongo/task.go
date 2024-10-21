package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTaskRepo struct {
	TaskColl *mongo.Collection
}

func NewMongoTaskRepo(db *mongo.Client) *MongoTaskRepo {
	return &MongoTaskRepo{TaskColl: db.Database(dbname).Collection(taskCollection)}
}

func (r *MongoTaskRepo) AddTask(ctx context.Context, task models.Task) error {
	_, err := r.TaskColl.InsertOne(ctx, task)
	return err
}

func (r *MongoTaskRepo) GetTaskList(ctx context.Context, userLogin string, groupID primitive.ObjectID) ([]models.Task, error) {
	var taskList []models.Task
	filter := bson.M{"group_id": groupID}
	cursor, err := r.TaskColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &taskList)
	if err != nil {
		logs.Error("cursorAll", err)
		return nil, err
	}
	return taskList, nil
}

func (r *MongoTaskRepo) GetTaskById(ctx context.Context, taskID, groupID primitive.ObjectID) (*models.Task, error) {
	var task models.Task
	filter := bson.M{
		"$and": []bson.M{
			{"_id": taskID},
			{"group_id": groupID},
		},
	}
	err := r.TaskColl.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *MongoTaskRepo) DeleteTask(ctx context.Context, taskID primitive.ObjectID) error {
	filter := bson.M{"_id": taskID}
	_, err := r.TaskColl.DeleteOne(ctx, filter)
	if err != nil {
		return  err
	}
	return  nil
}

func (r *MongoTaskRepo) UpdateTask(ctx context.Context, taskID primitive.ObjectID, newTask models.Task) error {
	filter := bson.M{"_id": taskID}

	update := bson.M{}
	if newTask.Title != "" {
		update["title"] = newTask.Title
	}
	if !newTask.StartTime.IsZero() {
		update["start_time"] = newTask.StartTime
	}
	if newTask.Duration != 0 {
		update["duration"] = newTask.Duration
	}
	if !newTask.EndTime.IsZero() {
		update["end_time"] = newTask.EndTime
	}
	updateQuery := bson.M{
        "$set": update,
    }

	fmt.Println(update)
	_, err := r.TaskColl.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		return err
	}
	return nil
}

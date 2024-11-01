package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTaskRepo struct {
	TaskColl *mongo.Collection
}

func NewMongoTaskRepo(db *mongo.Client) *MongoTaskRepo {
	return &MongoTaskRepo{TaskColl: db.Database(dbname).Collection(taskCollection)}
}

func (r *MongoTaskRepo) AddTask(ctx context.Context, task models.Task, groupID string) error {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	task.GroupID = oid[0]
	_, err = r.TaskColl.InsertOne(ctx, task)
	if err != nil {
		return fmt.Errorf("AddTask error: %v", err)
	}
	return nil
}

func (r *MongoTaskRepo) GetTaskList(ctx context.Context, userLogin, groupID string) ([]models.Task, error) {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return nil, fmt.Errorf("InvalidID: %v", err)
	}
	var taskList []models.Task
	filter := bson.M{"group_id": oid[0]}
	cursor, err := r.TaskColl.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("GetTaskList: %v", err)
	}
	err = cursor.All(ctx, &taskList)
	if err != nil {
		return nil, fmt.Errorf("GetTaskList all() error: %v", err)
	}
	return taskList, nil
}

func (r *MongoTaskRepo) GetTaskById(ctx context.Context, taskID, groupID string) (*models.Task, error) {
	oid, err := convertToObjectIDs(taskID, groupID)
	if err != nil {
		return nil, fmt.Errorf("InvalidID: %v", err)
	}
	var task models.Task
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"group_id": oid[1]},
		},
	}
	err = r.TaskColl.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return nil, fmt.Errorf("GetTaskByID error: %v", err)
	}
	return &task, nil
}

func (r *MongoTaskRepo) DeleteTask(ctx context.Context, taskID string) error {
	oid, err := convertToObjectIDs(taskID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{"_id": oid[0]}
	_, err = r.TaskColl.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("DeleteTask error: %v", err)
	}
	return nil
}

func (r *MongoTaskRepo) UpdateTask(ctx context.Context, taskID string, newTask models.Task) error {
	oid, err := convertToObjectIDs(taskID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{"_id": oid[0]}

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

	_, err = r.TaskColl.UpdateOne(ctx, filter, updateQuery)
	if err != nil {
		return fmt.Errorf("UpdateTask error: %v", err)
	}
	return nil
}

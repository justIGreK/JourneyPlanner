package mongorepo

import (
	"JourneyPlanner/internal/models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoPollRepo struct {
	PollColl *mongo.Collection
}

func NewMongoPollRepo(db *mongo.Client) *MongoPollRepo {
	return &MongoPollRepo{PollColl: db.Database(dbname).Collection(pollCollection)}
}

func (r *MongoPollRepo) CreatePoll(ctx context.Context, poll models.Poll, groupID string) error {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	poll.GroupID = oid[0]
	_, err = r.PollColl.InsertOne(ctx, poll)
	if err != nil {
		return fmt.Errorf("CreatePoll error: %v", err)
	}
	return nil
}

func (r *MongoPollRepo) GetPollList(ctx context.Context, groupID string) ([]models.Poll, []models.Poll, error) {
	oid, err := convertToObjectIDs(groupID)
	if err != nil {
		return nil, nil, fmt.Errorf("InvalidID: %v", err)
	}
	now := time.Now().UTC()
	getPolls := func(filter bson.M) ([]models.Poll, error) {
		cursor, err := r.PollColl.Find(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("GetPollList error: %v", err)
		}
		var polls []models.Poll
		if err := cursor.All(ctx, &polls); err != nil {
			return nil, fmt.Errorf("GetPollList error, cursor.All(): %v", err)
		}
		return polls, nil
	}
	closedFilter := bson.M{
		"group_id": oid[0],
		"$or": []bson.M{
			{"endtime": bson.M{"$lt": now}},
			{"isEarlyClosed": true},
		},
	}
	closedPolls, err := getPolls(closedFilter)
	if err != nil {
		return nil, nil, err
	}
	openFilter := bson.M{
		"group_id":      oid[0],
		"isEarlyClosed": false,
		"endtime":       bson.M{"$gt": now},
	}
	openPolls, err := getPolls(openFilter)
	if err != nil {
		return nil, nil, err
	}

	return openPolls, closedPolls, nil
}
func (r *MongoPollRepo) GetPollById(ctx context.Context, pollID string) (*models.Poll, error) {
	oid, err := convertToObjectIDs(pollID)
	if err != nil {
		return nil, fmt.Errorf("InvalidID: %v", err)
	}
	var poll models.Poll
	filter := bson.M{
		"_id": oid[0],
	}
	err = r.PollColl.FindOne(ctx, filter).Decode(&poll)
	if err != nil {
		return nil, fmt.Errorf("GetPollById error: %v", err)
	}
	return &poll, nil
}

func (r *MongoPollRepo) ClosePoll(ctx context.Context, pollID string) error {
	oid, err := convertToObjectIDs(pollID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"isEarlyClosed": false},
		},
	}
	update := bson.M{"$set": bson.M{"isEarlyClosed": true}}
	_, err = r.PollColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("ClosePoll error: %v", err)
	}
	return nil
}

func (r *MongoPollRepo) DeletePoll(ctx context.Context, pollID string) error {
	oid, err := convertToObjectIDs(pollID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	filter := bson.M{"_id": oid[0]}
	_, err = r.PollColl.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("Delete poll error: %v", err)
	}
	return nil
}

const (
	firstOption  = "firstOption"
	secondOption = "secondOption"
)

func (r *MongoPollRepo) AddVote(ctx context.Context, pollID, voteOption, userLogin string) error {
	oid, err := convertToObjectIDs(pollID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	now := time.Now().UTC()
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"endtime": bson.M{"$gt": now}},
			{"isEarlyClosed": false},
		},
	}
	var votesField string
	switch voteOption {
	case firstOption:
		votesField = "votes1"
	case secondOption:
		votesField = "votes2"
	default:
		return fmt.Errorf("unexpected error: invalid vote option %v", voteOption)
	}

	update := bson.M{"$push": bson.M{votesField: userLogin}}
	_, err = r.PollColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("ClosePoll error: %v", err)
	}
	return nil
}
func (r *MongoPollRepo) RemoveVote(ctx context.Context, pollID, userLogin string) error {
	oid, err := convertToObjectIDs(pollID)
	if err != nil {
		return fmt.Errorf("InvalidID: %v", err)
	}
	now := time.Now().UTC()
	filter := bson.M{
		"$and": []bson.M{
			{"_id": oid[0]},
			{"endtime": bson.M{"$gt": now}},
			{"isEarlyClosed": false},
		},
	}

	update := bson.M{"$pull": bson.M{
		"votes1": userLogin,
		"votes2": userLogin,
	}}
	_, err = r.PollColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("Remove vote error: %v", err)
	}
	return nil
}

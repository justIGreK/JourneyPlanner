package service

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PollRepository interface {
	CreatePoll(ctx context.Context, poll models.Poll) error
	GetPollList(ctx context.Context, groupOID primitive.ObjectID) ([]models.Poll,[]models.Poll, error)
	GetPollById(ctx context.Context, pollOID primitive.ObjectID) (*models.Poll, error)
	DeletePoll(ctx context.Context, pollOID primitive.ObjectID) error
	ClosePoll(ctx context.Context, pollOID primitive.ObjectID) error
	RemoveVote(ctx context.Context, pollOID primitive.ObjectID, userLogin string) error
	AddVote(ctx context.Context, pollOID primitive.ObjectID, voteOption, userLogin string) error
}

type PollSrv struct {
	Poll  PollRepository
	Group GroupRepository
}

func NewPollSrv(pollRepo PollRepository, groupRepo GroupRepository) *PollSrv {
	return &PollSrv{Poll: pollRepo, Group: groupRepo}
}

func (s *PollSrv) CreatePoll(ctx context.Context, pollInfo models.CreatePoll, userLogin string) error {
	groupOID, err := primitive.ObjectIDFromHex(pollInfo.GroupID)
	if err != nil {
		return err
	}

	_, err = s.Group.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return errors.New("this group is not exist or you are not member of it")
	}

	now := time.Now().UTC()
	votingEndTime := now.Add(time.Duration(pollInfo.Duration) * time.Minute)

	newPoll := models.Poll{
		GroupID:       groupOID,
		Creator:       userLogin,
		Title:         pollInfo.Title,
		FirstOption:   pollInfo.FirstOption,
		Votes1:        []string{},
		SecondOption:  pollInfo.SecondOption,
		Votes2:        []string{},
		EndTime:       votingEndTime,
		IsEarlyClosed: false,
	}
	err = s.Poll.CreatePoll(ctx, newPoll)
	return err
}

func (s *PollSrv) GetPollList(ctx context.Context, groupID, userLogin string) (*models.PollList, error) {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, err
	}
	_, err = s.Group.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return nil, errors.New("this group is not exist or you are not member of it")
	}

	openPolls, closedPolls, err := s.Poll.GetPollList(ctx, groupOID)
	if err != nil {
		return nil, err
	}

	pollList := s.preparePollList(openPolls, closedPolls)
	return &pollList, nil
}

func (s *PollSrv) preparePollList(openPolls, closedPolls []models.Poll) (models.PollList){
	var pollList models.PollList 
	for _, poll := range openPolls{
		pollList.OpenPolls = append(pollList.OpenPolls, models.PrintPollList{
			ID: poll.ID,
			Title: poll.Title,
			Creator: poll.Creator,
			FirstOption: poll.FirstOption,
			FirstVotesCount: len(poll.Votes1),
			SecondOption: poll.SecondOption,
			SecondVotesCount: len(poll.Votes2),
			EndTime: poll.EndTime.Format("2006-01-02 15:04:05"),
		})
	}
	for _, poll := range closedPolls{
		pollList.ClosedPolls = append(pollList.ClosedPolls, models.PrintPollList{
			ID: poll.ID,
			Title: poll.Title,
			Creator: poll.Creator,
			FirstOption: poll.FirstOption,
			FirstVotesCount: len(poll.Votes1),
			SecondOption: poll.SecondOption,
			SecondVotesCount: len(poll.Votes2),
			EndTime: poll.EndTime.Format("2006-01-02 15:04:05"),
		})
	}
	return pollList

}

func (s *PollSrv) DeletePollByID(ctx context.Context, pollID, groupID, userLogin string) error {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	group, err := s.Group.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return errors.New("this group is not exist or you are not member of it")
	}
	pollOID, err := primitive.ObjectIDFromHex(pollID)
	if err != nil {
		return err
	}
	poll, err := s.Poll.GetPollById(ctx, pollOID)
	if err != nil {
		return errors.New("poll is not found")
	}
	if group.LeaderLogin != userLogin && poll.Creator != userLogin {
		return errors.New("you have no permissions to do this")
	}
	err = s.Poll.DeletePoll(ctx, pollOID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PollSrv) ClosePoll(ctx context.Context, pollID, groupID, userLogin string) error {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return err
	}
	group, err := s.Group.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return errors.New("this group is not exist or you are not member of it")
	}
	pollOID, err := primitive.ObjectIDFromHex(pollID)
	if err != nil {
		return err
	}
	poll, err := s.Poll.GetPollById(ctx, pollOID)
	if err != nil {
		return errors.New("poll is not found")
	}
	now := time.Now().UTC()
	if poll.IsEarlyClosed || poll.EndTime.Before(now) {
		return errors.New("poll is already closed")
	}
	if group.LeaderLogin != userLogin && poll.Creator != userLogin {
		return errors.New("you have no permissions to do this")
	}
	err = s.Poll.ClosePoll(ctx, pollOID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PollSrv) VotePoll(ctx context.Context, userLogin string, vote models.AddVote) error {
	groupOID, err := primitive.ObjectIDFromHex(vote.GroupID)
	if err != nil {
		return err
	}
	_, err = s.Group.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return errors.New("this group is not exist or you are not member of it")
	}
	pollOID, err := primitive.ObjectIDFromHex(vote.PollID)
	if err != nil {
		return err
	}
	poll, err := s.Poll.GetPollById(ctx, pollOID)
	if err != nil {
		return errors.New("poll is not found")
	}
	now := time.Now().UTC()
	if poll.IsEarlyClosed || poll.EndTime.Before(now) {
		return errors.New("poll is already closed")
	}
	err = s.Poll.RemoveVote(ctx, pollOID, userLogin)
	if err != nil {
		return err
	}
	err = s.Poll.AddVote(ctx, pollOID, vote.Option, userLogin)
	if err != nil {
		return err
	}
	return nil
}


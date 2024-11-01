package service

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"
	"time"
)

type PollRepository interface {
	CreatePoll(ctx context.Context, poll models.Poll, groupID string) error
	GetPollList(ctx context.Context, groupID string) ([]models.Poll, []models.Poll, error)
	GetPollById(ctx context.Context, pollID string) (*models.Poll, error)
	DeletePoll(ctx context.Context, pollID string) error
	ClosePoll(ctx context.Context, pollID string) error
	RemoveVote(ctx context.Context, pollID, userLogin string) error
	AddVote(ctx context.Context, pollID, voteOption, userLogin string) error
}

type PollSrv struct {
	Poll  PollRepository
	Group GroupRepository
}

func NewPollSrv(pollRepo PollRepository, groupRepo GroupRepository) *PollSrv {
	return &PollSrv{Poll: pollRepo, Group: groupRepo}
}

func (s *PollSrv) CreatePoll(ctx context.Context, pollInfo models.CreatePoll, userLogin string) error {
	_, err := s.Group.GetGroup(ctx, pollInfo.GroupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("this group is not exist or you are not member of it")
	}

	now := time.Now().UTC()
	votingEndTime := now.Add(time.Duration(pollInfo.Duration) * time.Minute)

	newPoll := models.Poll{
		Creator:       userLogin,
		Title:         pollInfo.Title,
		FirstOption:   pollInfo.FirstOption,
		Votes1:        []string{},
		SecondOption:  pollInfo.SecondOption,
		Votes2:        []string{},
		EndTime:       votingEndTime,
		IsEarlyClosed: false,
	}
	err = s.Poll.CreatePoll(ctx, newPoll, pollInfo.GroupID)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}
	return nil
}

func (s *PollSrv) GetPollList(ctx context.Context, groupID, userLogin string) (*models.PollList, error) {
	_, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("this group is not exist or you are not member of it")
	}

	openPolls, closedPolls, err := s.Poll.GetPollList(ctx, groupID)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("System error")
	}

	pollList := s.preparePollList(openPolls, closedPolls)
	return &pollList, nil
}

func (s *PollSrv) preparePollList(openPolls, closedPolls []models.Poll) models.PollList {
	pollList := models.PollList{}
	for _, poll := range openPolls {
		printPoll := models.PrintPollList{
			ID:               poll.ID,
			Title:            poll.Title,
			Creator:          poll.Creator,
			FirstOption:      poll.FirstOption,
			FirstVotesCount:  len(poll.Votes1),
			SecondOption:     poll.SecondOption,
			SecondVotesCount: len(poll.Votes2),
			EndTime:          poll.EndTime.Format("2006-01-02 15:04:05"),
		}
		pollList.OpenPolls = append(pollList.OpenPolls, printPoll)
	}
	for _, poll := range closedPolls {
		printPoll := models.PrintPollList{
			ID:               poll.ID,
			Title:            poll.Title,
			Creator:          poll.Creator,
			FirstOption:      poll.FirstOption,
			FirstVotesCount:  len(poll.Votes1),
			SecondOption:     poll.SecondOption,
			SecondVotesCount: len(poll.Votes2),
			EndTime:          poll.EndTime.Format("2006-01-02 15:04:05"),
		}
		pollList.ClosedPolls = append(pollList.ClosedPolls, printPoll)
	}
	return pollList
}

func (s *PollSrv) DeletePollByID(ctx context.Context, pollID, groupID, userLogin string) error {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("this group is not exist or you are not member of it")
	}
	poll, err := s.Poll.GetPollById(ctx, pollID)
	if err != nil {
		logs.Error(err)
		return errors.New("poll is not found")
	}
	if group.LeaderLogin != userLogin && poll.Creator != userLogin {
		return errors.New("you have no permissions to do this")
	}
	err = s.Poll.DeletePoll(ctx, pollID)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}
	return nil
}

func (s *PollSrv) ClosePoll(ctx context.Context, pollID, groupID, userLogin string) error {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("this group is not exist or you are not member of it")
	}
	poll, err := s.Poll.GetPollById(ctx, pollID)
	if err != nil {
		logs.Error(err)
		return errors.New("poll is not found")
	}
	now := time.Now().UTC()
	if poll.IsEarlyClosed || poll.EndTime.Before(now) {
		return errors.New("poll is already closed")
	}
	if group.LeaderLogin != userLogin && poll.Creator != userLogin {
		return errors.New("you have no permissions to do this")
	}
	err = s.Poll.ClosePoll(ctx, pollID)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}
	return nil
}

func (s *PollSrv) VotePoll(ctx context.Context, userLogin string, vote models.AddVote) error {
	_, err := s.Group.GetGroup(ctx, vote.GroupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("this group is not exist or you are not member of it")
	}
	poll, err := s.Poll.GetPollById(ctx, vote.PollID)
	if err != nil {
		logs.Error(err)
		return errors.New("poll is not found")
	}
	now := time.Now().UTC()
	if poll.IsEarlyClosed || poll.EndTime.Before(now) {
		return errors.New("poll is already closed")
	}
	err = s.Poll.RemoveVote(ctx, vote.PollID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}
	err = s.Poll.AddVote(ctx, vote.PollID, vote.Option, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("System error")
	}
	return nil
}

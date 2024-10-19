package service

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"
	"math/rand"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var logs *zap.SugaredLogger

func SetLogger(l *zap.Logger) {
	logs = l.Sugar()
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, group models.Group) error
	GetGroupList(ctx context.Context, userLogin string) ([]models.Group, error)
	GetGroupById(ctx context.Context, groupID primitive.ObjectID, userLogin string) (*models.Group, error)
	ChangeGroupLeader(ctx context.Context, groupID primitive.ObjectID, userLogin string) error
	DeleteGroup(ctx context.Context, groupID primitive.ObjectID) error
	LeaveGroup(ctx context.Context, groupID primitive.ObjectID, userLogin string) error
}

type GroupSrv struct {
	GroupRepository
	UserRepository
}

func NewGroupSrv(groupRepo GroupRepository, userRepo UserRepository) *GroupSrv {
	return &GroupSrv{GroupRepository: groupRepo, UserRepository: userRepo}
}

func (s *GroupSrv) CreateGroup(ctx context.Context, groupName, userLogin string, invites []string) error {

	group := models.Group{
		Name:        groupName,
		LeaderLogin: userLogin,
		Members:     []string{userLogin},
		Tasks:       []models.Task{},
		Polls:       []models.Poll{},
		IsActive:    true,
	}
	err := s.GroupRepository.CreateGroup(ctx, group)
	return err
}

func (s *GroupSrv) GetGroupList(ctx context.Context, userLogin string) ([]models.GroupList, error) {
	groups, err := s.GroupRepository.GetGroupList(ctx, userLogin)
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, errors.New("your grouplist is empty")
	}

	var groupsList []models.GroupList
	for _, group := range groups {
		groupsList = append(groupsList, models.GroupList{
			ID:           group.ID,
			Name:         group.Name,
			MembersCount: len(group.Members),
		})
	}
	return groupsList, nil

}

func (s *GroupSrv) GetGroup(ctx context.Context, groupID, userLogin string) (*models.Group, error) {
	oid, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return nil, errors.New("InvalidID")
	}
	group, err := s.GroupRepository.GetGroupById(ctx, oid, userLogin)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (s *GroupSrv) LeaveGroup(ctx context.Context, groupID, userLogin string) error {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New("InvalidID")
	}
	group, err := s.GroupRepository.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return err
	}
	if len(group.Members) <= 1 {
		if err := s.GroupRepository.DeleteGroup(ctx, groupOID); err != nil {
			return err
		}
		return nil
	} else {
		if group.LeaderLogin == userLogin {
			newLeader := s.getRandomLeader(group.Members, userLogin)
			err := s.GroupRepository.ChangeGroupLeader(ctx, groupOID, newLeader)
			if err != nil {
				return err
			}
		}

		err := s.GroupRepository.LeaveGroup(ctx, groupOID, userLogin)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *GroupSrv) getRandomLeader(members []string, userLogin string) string {
	removeUser := func(slice []string, value string) []string {
		newSlice := []string{}
		for _, v := range slice {
			if v != value {
				newSlice = append(newSlice, v)
			}
		}
		return newSlice
	}
	members = removeUser(members, userLogin)
	randomNumber := rand.Intn(len(members))
	return members[randomNumber]

}
func (s *GroupSrv) GiveLeaderRole(ctx context.Context, groupID, userLogin, memberLogin string) error {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New("InvalidID")
	}
	group, err := s.GroupRepository.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return err
	}
	var isRealMember bool
	for _, member := range group.Members {
		if member == memberLogin{
			isRealMember = true 
			break
		}
	}
	if !isRealMember{
		return errors.New("no such member")
	}
	if group.LeaderLogin != userLogin{
		return errors.New("you have no permissions to do this")
	}
	err = s.GroupRepository.ChangeGroupLeader(ctx, groupOID, memberLogin)
	if err != nil {
		return err
	}
	return nil
}

func (s *GroupSrv) DeleteGroup(ctx context.Context, groupID, userLogin string) error {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New("InvalidID")
	}
	group, err := s.GroupRepository.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return err
	}
	if group.LeaderLogin != userLogin{
		return errors.New("you have no permissions to do this")
	}
	err = s.GroupRepository.DeleteGroup(ctx, groupOID)
	if err != nil {
		return err
	}
	return nil
}

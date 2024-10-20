package service

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	JoinGroup(ctx context.Context, groupID primitive.ObjectID, userLogin string) error
	LeaveGroup(ctx context.Context, groupID primitive.ObjectID, userLogin string) error
	CheckGroupForExist(ctx context.Context, groupID primitive.ObjectID) (bool, *models.Group, error) 
}
type InviteRepository interface {
	AddInvitation(ctx context.Context, inviteCreator, invitedUser, groupName, token string) error
	GetInvites(ctx context.Context, userLogin string) ([]models.Invitation, error)
	DeleteInviteByID(ctx context.Context, inviteID primitive.ObjectID, userLogin string) (int64, error)
	DeleteInviteByToken(ctx context.Context, token string)  error
}

type GroupSrv struct {
	GroupRepository
	UserRepository
	InviteRepository
}

func NewGroupSrv(groupRepo GroupRepository, userRepo UserRepository, inviteRepo InviteRepository) *GroupSrv {
	return &GroupSrv{GroupRepository: groupRepo, UserRepository: userRepo, InviteRepository: inviteRepo}
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
	if group.LeaderLogin != userLogin {
		return errors.New("you have no permissions to do this")
	}
	var isRealMember bool
	for _, member := range group.Members {
		if member == memberLogin {
			isRealMember = true
			break
		}
	}
	if !isRealMember {
		return errors.New("no such member")
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
	if group.LeaderLogin != userLogin {
		return errors.New("you have no permissions to do this")
	}
	err = s.GroupRepository.DeleteGroup(ctx, groupOID)
	if err != nil {
		return err
	}
	return nil
}

func (s *GroupSrv) InviteUser(ctx context.Context, groupID, userLogin string, invitedUser string) error {
	groupOID, err := primitive.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New("InvalidID")
	}
	// ============================================================================== start checking
	group, err := s.GroupRepository.GetGroupById(ctx, groupOID, userLogin)
	if err != nil {
		return errors.New("group is not found, or you are not a member of it")
	}

	_, err = s.GetUserByLogin(ctx, invitedUser)
	if err != nil {
		return errors.New("user not found")
	}
	isOkay := true
	for _, member := range group.Members {
		if member == invitedUser {
			isOkay = false
			break
		}
	}
	if !isOkay {
		return errors.New("user is alredy member of this group")
	}
	// ============================================================================== end of check

	inviteToken, err := s.GetInviteToken(invitedUser, groupID)
	if err != nil {
		return err
	}

	err = s.InviteRepository.AddInvitation(ctx, userLogin, invitedUser, group.Name, inviteToken)
	if err != nil {
		return err
	}
	return nil
}

func (s *GroupSrv) GetInviteToken(invitedUser, groupID string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	expirationTime := time.Now().UTC().Add(24 * time.Hour)

	claims := &models.InvitationToken{
		UserLogin: invitedUser,
		GroupID:   groupID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func (s *GroupSrv) JoinGroup(ctx context.Context, token string) error {
	inviteDetails, err := s.ValidateInvitationToken(token)
	if err != nil {
		return err
	}
	groupOID, err := primitive.ObjectIDFromHex(inviteDetails.GroupID)
	if err != nil {
		return errors.New("invalid group id")
	}
	_, err = s.GetUserByLogin(ctx, inviteDetails.UserLogin)
	if err !=nil{
		return errors.New("user was not found")
	}
	isOkay, group, err := s.CheckGroupForExist(ctx, groupOID)
	if !isOkay{
		return errors.New("this group is no longer exist")
	}
	for _, member := range group.Members{
		if member == inviteDetails.UserLogin{
			isOkay = false
		}
	}
	if !isOkay{
		return errors.New("You are already member of this group")
	}
	err = s.GroupRepository.JoinGroup(ctx, groupOID, inviteDetails.UserLogin)
	if err != nil{
		return err
	}

	err = s.DeleteInviteByToken(ctx, token)
	if err != nil{
		logs.Warn(err)
	}
	return nil
}

func (s *GroupSrv) ValidateInvitationToken(tokenString string) (*models.InvitationToken, error) {
	secretKey := os.Getenv("SECRET_KEY")
	claims := &models.InvitationToken{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("expired token")
	}

	return claims, nil
}

func (s *GroupSrv) GetInviteList(ctx context.Context, userLogin string) ([]models.InvitationList, error) {
	invites, err := s.GetInvites(ctx, userLogin)
	if err != nil {
		return nil, err
	}
	inviteList := s.inviteFormat(invites)
	return inviteList, nil
}

func (s *GroupSrv) inviteFormat(invites []models.Invitation) []models.InvitationList {
	var inviteList []models.InvitationList
	for _, invite := range invites {
		inviteList = append(inviteList, models.InvitationList{
			Invite_ID: invite.Invite_ID,
			InvitationText: fmt.Sprintf("User %v invited you to the group %v", invite.Sender, invite.GroupName),
			InvitationLink: fmt.Sprintf("http://localhost:8080/join-group?token=%s", invite.Token),
		})
	}
	return inviteList
}

func (s *GroupSrv) DeclineInvite(ctx context.Context, userLogin, inviteID string)error{
	inviteOID, err := primitive.ObjectIDFromHex(inviteID)
	if err != nil {
		return err
	}
	modDocs, err := s.InviteRepository.DeleteInviteByID(ctx, inviteOID, userLogin)
	if err != nil{
		return err
	}
	if modDocs == 0{
		return errors.New("invite wasn't found")
	}
	return nil
}
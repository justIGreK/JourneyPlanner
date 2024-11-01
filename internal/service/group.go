package service

import (
	"JourneyPlanner/internal/models"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

var logs *zap.SugaredLogger

func SetLogger(l *zap.Logger) {
	logs = l.Sugar()
}

type GroupRepository interface {
	CreateGroup(ctx context.Context, group models.Group) (string, error)
	GetGroupList(ctx context.Context, userLogin string) ([]models.Group, error)
	GetGroup(ctx context.Context, groupID string, userLogin ...string) (*models.Group, error)
	ChangeGroupLeader(ctx context.Context, groupID, userLogin string) error
	DeleteGroup(ctx context.Context, groupID string) error
	JoinGroup(ctx context.Context, groupID, userLogin string) error
	LeaveGroup(ctx context.Context, groupID, userLogin string) error
}
type InviteRepository interface {
	AddInvitation(ctx context.Context, invite models.Invitation) error
	GetInvites(ctx context.Context, userLogin string) ([]models.Invitation, error)
	DeleteInviteByID(ctx context.Context, inviteID, userLogin string) (int64, error)
	DeleteInviteByToken(ctx context.Context, token string) error
	IsAlreadyInvited(ctx context.Context, groupID, userLogin string) (bool, error)
}
type BlackListRepository interface {
	CreateBlacklist(ctx context.Context, groupID string) error
	BanUser(ctx context.Context, groupID, userLogin string) error
	UnbanUser(ctx context.Context, groupID, userLogin string) error
	GetBlacklist(ctx context.Context, groupID string) (*models.BlackList, error)
}

type WebSockerConn interface {
	KickUser(userLogin, groupID string)
}

type GroupSrv struct {
	Group                GroupRepository
	User                 UserRepository
	Invite               InviteRepository
	BlackList            BlackListRepository
	NotifyUserDisconnect func(userLogin string, groupID string)
}

func NewGroupSrv(groupRepo GroupRepository, userRepo UserRepository,
	inviteRepo InviteRepository, blackList BlackListRepository) *GroupSrv {
	return &GroupSrv{Group: groupRepo, User: userRepo,
		Invite: inviteRepo, BlackList: blackList}
}

func (s *GroupSrv) CreateGroup(ctx context.Context, groupName, userLogin string) error {
	group := models.Group{
		Name:        groupName,
		LeaderLogin: userLogin,
		Members:     []string{userLogin},
		Tasks:       []models.Task{},
		Polls:       []models.Poll{},
		IsActive:    true,
	}
	groupOID, err := s.Group.CreateGroup(ctx, group)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("createGroup error: %v", err)
	}
	err = s.BlackList.CreateBlacklist(ctx, groupOID)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("create blacklist error: %v", err)
	}
	return nil
}

func (s *GroupSrv) GetGroupList(ctx context.Context, userLogin string) ([]models.GroupList, error) {
	groups, err := s.Group.GetGroupList(ctx, userLogin)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("failed to get groups")
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

func (s *GroupSrv) GetGroupByID(ctx context.Context, groupID, userLogin string) (*models.Group, error) {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("failed to get group")
	}
	return group, nil
}

func (s *GroupSrv) BanMember(ctx context.Context, groupID, memberLogin, userLogin string) error {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to ban user")
	}
	if group.LeaderLogin != userLogin {
		return errors.New("you have no permissions to do this")
	}
	isOkay := false
	for _, member := range group.Members {
		if member == memberLogin {
			isOkay = true
			break
		}
	}
	if !isOkay {
		return errors.New("Member is not found")
	}
	err = s.Group.LeaveGroup(ctx, groupID, memberLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to ban user")
	}
	err = s.BlackList.BanUser(ctx, groupID, memberLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to ban user")
	}
	s.NotifyUserDisconnect(userLogin, groupID)
	return nil
}
func (s *GroupSrv) UnbanMember(ctx context.Context, groupID, memberLogin, userLogin string) error {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to unban user")
	}
	if group.LeaderLogin != userLogin {
		return errors.New("you have no permissions to do this")
	}
	blacklist, err := s.BlackList.GetBlacklist(ctx, groupID)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to unban user")
	}
	isOkay := false
	for _, member := range blacklist.Blacklist {
		if member == memberLogin {
			isOkay = true
			break
		}
	}
	if !isOkay {
		return errors.New("this user is not banned in this group")
	}
	err = s.BlackList.UnbanUser(ctx, groupID, memberLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to unban user")
	}
	return nil
}

func (s *GroupSrv) GetBlacklist(ctx context.Context, groupID, userLogin string) (*models.BlackList, error) {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("failed to get black list")
	}
	if group.LeaderLogin != userLogin {
		return nil, errors.New("you have no permissions to do this")
	}

	blacklist, err := s.BlackList.GetBlacklist(ctx, groupID)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("failed to get black list")
	}
	return blacklist, nil
}

func (s *GroupSrv) LeaveGroup(ctx context.Context, groupID, userLogin string) error {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to leave group")
	}
	if len(group.Members) <= 1 {
		if err := s.Group.DeleteGroup(ctx, groupID); err != nil {
			logs.Error(err)
			return errors.New("failed to leave group")
		}
		return nil
	} else {
		if group.LeaderLogin == userLogin {
			newLeader := s.getRandomLeader(group.Members, userLogin)
			err := s.Group.ChangeGroupLeader(ctx, groupID, newLeader)
			if err != nil {
				logs.Error(err)
				return errors.New("failed to change group leader")
			}
		}

		err := s.Group.LeaveGroup(ctx, groupID, userLogin)
		if err != nil {
			logs.Error(err)
			return errors.New("failed to leave group, please try later")
		}
	}
	s.NotifyUserDisconnect(userLogin, groupID)
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
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(members))))
	if err != nil {
		logs.Errorf("error during getting random member: %v", err)
		return ""
	}
	return members[int32(n.Int64())]
}
func (s *GroupSrv) GiveLeaderRole(ctx context.Context, groupID, userLogin, memberLogin string) error {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to get group")
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
	err = s.Group.ChangeGroupLeader(ctx, groupID, memberLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to change leader")
	}
	return nil
}

func (s *GroupSrv) DeleteGroup(ctx context.Context, groupID, userLogin string) error {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("Failed to find group")
	}
	if group.LeaderLogin != userLogin {
		return errors.New("you have no permissions to do this")
	}
	err = s.Group.DeleteGroup(ctx, groupID)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to delete group")
	}
	return nil
}

func (s *GroupSrv) InviteUser(ctx context.Context, groupID, userLogin, invitedUser string) error {
	group, err := s.Group.GetGroup(ctx, groupID, userLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("group is not found, or you are not a member of it")
	}

	_, err = s.User.GetUserByLogin(ctx, invitedUser)
	if err != nil {
		logs.Error(err)
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
	blacklist, err := s.BlackList.GetBlacklist(ctx, groupID)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to get blacklist of group")
	}
	for _, member := range blacklist.Blacklist {
		if member == invitedUser {
			isOkay = false
			break
		}
	}
	if !isOkay {
		return errors.New("this user have been banned from this group")
	}
	isOkay, err = s.Invite.IsAlreadyInvited(ctx, groupID, invitedUser)
	if !isOkay {
		logs.Error(err)
		return errors.New("this user is already invited to this group")
	}

	inviteToken, err := s.GetInviteToken(invitedUser, groupID)
	if err != nil {
		logs.Error(err)
		return errors.New("failed to create invite, please try later")
	}

	invite := models.Invitation{
		Sender:    userLogin,
		Receiver:  invitedUser,
		GroupID:   group.ID,
		GroupName: group.Name,
		Token:     inviteToken,
		IsUsed:    false,
	}
	err = s.Invite.AddInvitation(ctx, invite)
	if err != nil {
		logs.Error(err)
		return errors.New("Failed to send invitation")
	}
	return nil
}

func (s *GroupSrv) GetInviteToken(invitedUser, groupID string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	expirationTime := time.Now().UTC().Add(HoursInDay * time.Hour)

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
		return "", fmt.Errorf("Failed to create token")
	}
	return tokenString, nil
}
func (s *GroupSrv) JoinGroup(ctx context.Context, token string) error {
	inviteDetails, err := s.ValidateInvitationToken(token)
	if err != nil {
		logs.Error(err)
		return errors.New("Failed to confirm your invitation, please try again later")
	}
	_, err = s.User.GetUserByLogin(ctx, inviteDetails.UserLogin)
	if err != nil {
		logs.Error(err)
		return errors.New("user was not found")
	}
	group, err := s.Group.GetGroup(ctx, inviteDetails.GroupID)
	if group != nil {
		logs.Error(err)
		return errors.New("this group is no longer exist")
	}
	isOkay := true
	for _, member := range group.Members {
		if member == inviteDetails.UserLogin {
			isOkay = false
			break
		}
	}
	if !isOkay {
		return errors.New("You are already member of this group")
	}
	blacklist, err := s.BlackList.GetBlacklist(ctx, inviteDetails.GroupID)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("cant get blacklist of group, %v", err)
	}
	for _, member := range blacklist.Blacklist {
		if member == inviteDetails.UserLogin {
			isOkay = false
			break
		}
	}
	if !isOkay {
		return errors.New("You have been banned from this group")
	}

	err = s.Group.JoinGroup(ctx, inviteDetails.GroupID, inviteDetails.UserLogin)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("JoinGroup error: %v", err)
	}

	err = s.Invite.DeleteInviteByToken(ctx, token)
	if err != nil {
		logs.Error(err)
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
		return nil, fmt.Errorf("ValidateInviteToken error: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("expired token")
	}

	return claims, nil
}

func (s *GroupSrv) GetInviteList(ctx context.Context, userLogin string) ([]models.InvitationList, error) {
	invites, err := s.Invite.GetInvites(ctx, userLogin)
	if err != nil {
		logs.Error(err)
		return nil, fmt.Errorf("GetInviteList error: %v", err)
	}
	inviteList := s.inviteFormat(invites)
	return inviteList, nil
}

func (s *GroupSrv) inviteFormat(invites []models.Invitation) []models.InvitationList {
	var inviteList []models.InvitationList
	for _, invite := range invites {
		inviteList = append(inviteList, models.InvitationList{
			Invite_ID:      invite.Invite_ID,
			InvitationText: fmt.Sprintf("User %v invited you to the group %v", invite.Sender, invite.GroupName),
			InvitationLink: fmt.Sprintf("http://localhost:8080/join-group?token=%s", invite.Token),
		})
	}
	return inviteList
}

func (s *GroupSrv) DeclineInvite(ctx context.Context, userLogin, inviteID string) error {
	modDocs, err := s.Invite.DeleteInviteByID(ctx, inviteID, userLogin)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("DeclineInvite error: %v", err)
	}
	if modDocs == 0 {
		return errors.New("invite wasn't found")
	}
	return nil
}

package service

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
}

type UserSrv struct {
	User UserRepository
}

func NewUserSrv(userRepo UserRepository) *UserSrv {
	return &UserSrv{User: userRepo}
}

func (s *UserSrv) RegisterUser(ctx context.Context, user models.SignUp) error {
	if !s.isValidEmail(user.Email) {
		return errors.New("invalid email")
	}
	err := s.duplicateCheck(ctx, user)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Error("error of generating password", err)
		return errors.New("unfortunately we were unable to process your request, please try again later")
	}
	newUser := models.User{
		Login:        user.Login,
		Email:        user.Email,
		Password:     user.Password,
		PasswordHash: string(hashedPassword),
	}
	err = s.User.CreateUser(ctx, newUser)
	if err != nil {
		logs.Error(err)
		return fmt.Errorf("error during creating user:%v", err)
	}
	return nil
}

func (s *UserSrv) LoginUser(ctx context.Context, option, password string) (string, error) {
	var user *models.User
	var err error
	if s.isValidEmail(option) {
		user, err = s.User.GetUserByEmail(ctx, option)
	} else {
		user, err = s.User.GetUserByLogin(ctx, option)
	}
	if err != nil {
		logs.Error(err)
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		logs.Error(err)
		return "", errors.New("invalid credentials")
	}

	token, err := s.GeneratePasetoToken(user.Login)
	if err != nil {
		logs.Error(err)
		return "", fmt.Errorf("error during generating token: %v", err)
	}

	return token, nil
}

const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

func (s *UserSrv) isValidEmail(email string) bool {
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func (s *UserSrv) duplicateCheck(ctx context.Context, user models.SignUp) error {
	if _, err := s.User.GetUserByLogin(ctx, user.Login); err == nil {
		logs.Error(err)
		return errors.New("this login is already registered")
	}
	if _, err := s.User.GetUserByEmail(ctx, user.Email); err == nil {
		logs.Error(err)
		return errors.New("this email is already registered")
	}
	return nil
}

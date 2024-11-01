package chat

import (
	"JourneyPlanner/internal/models"
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

var logs *zap.SugaredLogger

func SetLogger(l *zap.Logger) {
	logs = l.Sugar()
}

type ChatRepository interface {
	InsertMessage(ctx context.Context, msg models.Message) error
	FindMessagesByChatID(ctx context.Context, groupID string) ([]models.Message, error)
}

type ChatSrv struct {
	repo ChatRepository
}

func NewChatService(repo ChatRepository) *ChatSrv {
	return &ChatSrv{repo: repo}
}

func (s *ChatSrv) SaveMessage(ctx context.Context, msg models.Message) error {
	msg.Time = time.Now().UTC()
	err := s.repo.InsertMessage(ctx, msg)
	if err != nil {
		logs.Error(err)
		return errors.New("failed save message")
	}
	return nil
}

func (s *ChatSrv) GetChatHistory(ctx context.Context, groupID string) ([]models.Message, error) {
	messages, err := s.repo.FindMessagesByChatID(ctx, groupID)
	if err != nil {
		logs.Error(err)
		return nil, errors.New("failed to get history")
	}
	return messages, nil
}

package chat

import (
	"JourneyPlanner/internal/models"
	"context"
	"time"
)

type ChatRepository interface{
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
	return s.repo.InsertMessage(ctx, msg)
}

func (s *ChatSrv) GetChatHistory(ctx context.Context, groupID string) ([]models.Message, error) {
	return s.repo.FindMessagesByChatID(ctx, groupID)
}

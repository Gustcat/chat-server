package repository

import (
	"context"
	"github.com/Gustcat/chat-server/internal/model"
)

type ChatRepository interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type MessageRepository interface {
	SendMessage(ctx context.Context, message *model.Message) error
}

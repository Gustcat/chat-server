package repository

import (
	"context"
	repomodel "github.com/Gustcat/chat-server/internal/repository/message/model"
)

type ChatRepository interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type MessageRepository interface {
	SendMessage(ctx context.Context, message *repomodel.Message) error
}

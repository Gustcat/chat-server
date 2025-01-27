package service

import (
	"context"
	"github.com/Gustcat/chat-server/internal/model"
)

type ChatService interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type MessageService interface {
	SendMessage(ctx context.Context, message *model.Message) error
}

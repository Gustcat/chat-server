package converter

import (
	"github.com/Gustcat/chat-server/internal/model"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
)

func ToMessageFromDesc(message *desc.Message) *model.Message {
	return &model.Message{
		Username:  message.From,
		Text:      message.Text,
		CreatedAt: message.Timestamp.AsTime(),
		ChatID:    message.ChatId,
	}
}

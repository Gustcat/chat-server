package converter

import (
	"github.com/Gustcat/chat-server/internal/model"
	repomodel "github.com/Gustcat/chat-server/internal/repository/message/model"
)

func ToMessageFromService(message *model.Message) *repomodel.Message {
	return &repomodel.Message{
		ID:        message.ID,
		Username:  message.Username,
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
		ChatID:    message.ChatID,
	}
}

package message

import (
	"github.com/Gustcat/chat-server/internal/repository"
	"github.com/Gustcat/chat-server/internal/service"
)

type serv struct {
	messageRepository repository.MessageRepository
}

func NewMessageService(messageRepository repository.MessageRepository) service.MessageService {
	return &serv{messageRepository: messageRepository}
}

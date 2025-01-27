package chat

import (
	"github.com/Gustcat/chat-server/internal/repository"
	"github.com/Gustcat/chat-server/internal/service"
)

type serv struct {
	chatRepository repository.ChatRepository
}

func NewChatService(chatRepository repository.ChatRepository) service.ChatService {
	return &serv{
		chatRepository: chatRepository,
	}
}

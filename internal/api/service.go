package api

import (
	"github.com/Gustcat/chat-server/internal/service"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
)

type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService    service.ChatService
	messageService service.MessageService
}

func NewImplementation(
	chatService service.ChatService,
	messageService service.MessageService) *Implementation {
	return &Implementation{
		chatService:    chatService,
		messageService: messageService,
	}
}

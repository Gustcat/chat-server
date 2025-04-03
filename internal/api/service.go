package api

import (
	"github.com/Gustcat/chat-server/internal/service"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
	"sync"
)

type Chat struct {
	streams map[string]desc.ChatV1_ConnectServer
	m       sync.RWMutex
}

type Implementation struct {
	desc.UnimplementedChatV1Server
	chatService    service.ChatService
	messageService service.MessageService

	chats  map[int64]*Chat
	mxChat sync.RWMutex

	channels  map[int64]chan *desc.Message
	mxChannel sync.RWMutex
}

func NewImplementation(
	chatService service.ChatService,
	messageService service.MessageService) *Implementation {
	return &Implementation{
		chatService:    chatService,
		messageService: messageService,
		chats:          make(map[int64]*Chat),
		channels:       make(map[int64]chan *desc.Message),
	}
}

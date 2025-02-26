package message

import (
	db "github.com/Gustcat/chat-server/internal/client"
	"github.com/Gustcat/chat-server/internal/repository"
	"github.com/Gustcat/chat-server/internal/service"
)

type serv struct {
	messageRepository repository.MessageRepository
	txManager         db.TxManager
}

func NewMessageService(messageRepository repository.MessageRepository, txManager db.TxManager) service.MessageService {
	return &serv{
		messageRepository: messageRepository,
		txManager:         txManager,
	}
}

func NewMockService(deps ...interface{}) service.MessageService {
	srv := serv{}

	for _, v := range deps {
		switch s := v.(type) {
		case repository.MessageRepository:
			srv.messageRepository = s
		}
	}

	return &srv
}

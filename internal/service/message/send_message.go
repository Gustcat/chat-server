package message

import (
	"context"
	"github.com/Gustcat/chat-server/internal/model"
	"github.com/Gustcat/chat-server/internal/repository/message/converter"
)

func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {
	err := s.messageRepository.SendMessage(ctx, converter.ToMessageFromService(message))
	if err != nil {
		return err
	}

	return nil
}

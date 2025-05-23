package message

import (
	"context"
	"github.com/Gustcat/chat-server/internal/model"
)

func (s *serv) SendMessage(ctx context.Context, message *model.Message) error {
	err := s.messageRepository.SendMessage(ctx, message)
	if err != nil {
		return err
	}

	return nil
}

package api

import (
	"context"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.chatService.Create(ctx, req.GetUsernames())
	if err != nil {
		return nil, err
	}

	i.channels[id] = make(chan *desc.Message, 100)

	return &desc.CreateResponse{Id: id}, nil
}

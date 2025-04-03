package api

import (
	"context"
	"github.com/Gustcat/chat-server/internal/converter"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	i.mxChannel.RLock()
	chatChan, ok := i.channels[req.GetMessage().GetChatId()]
	i.mxChannel.RUnlock()

	if !ok {
		return nil, status.Error(codes.NotFound, "chat not found")
	}

	err := i.messageService.SendMessage(ctx, converter.ToMessageFromDesc(req.GetMessage()))
	if err != nil {
		return nil, err
	}

	chatChan <- req.GetMessage()

	return &emptypb.Empty{}, nil
}

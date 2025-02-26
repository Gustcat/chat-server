package tests

import (
	"context"
	"fmt"
	chat "github.com/Gustcat/chat-server/internal/api"
	"github.com/Gustcat/chat-server/internal/model"
	"github.com/Gustcat/chat-server/internal/service"
	servicemocks "github.com/Gustcat/chat-server/internal/service/mocks"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type messageServiceMockFunc func(mc *minimock.Controller) service.MessageService

	type args struct {
		ctx context.Context
		req *desc.SendMessageRequest
	}

	var (
		ctx = context.Background()

		username  = gofakeit.Name()
		text      = gofakeit.Sentence(5)
		createdAt = gofakeit.Date()
		chatId    = gofakeit.Int64()

		req = &desc.SendMessageRequest{
			Message: &desc.Message{
				From:      username,
				Text:      text,
				Timestamp: timestamppb.New(createdAt),
				ChatId:    chatId,
			},
		}

		message = &model.Message{
			Username:  username,
			Text:      text,
			CreatedAt: createdAt,
			ChatID:    chatId,
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name               string
		args               args
		err                error
		expected           *emptypb.Empty
		messageServiceMock messageServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				req: req,
				ctx: ctx,
			},
			err:      nil,
			expected: &emptypb.Empty{},
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				mock := servicemocks.NewMessageServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				req: req,
				ctx: ctx,
			},
			err:      serviceErr,
			expected: nil,
			messageServiceMock: func(mc *minimock.Controller) service.MessageService {
				mock := servicemocks.NewMessageServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)
			messageServiceMock := tt.messageServiceMock(mc)
			api := chat.NewImplementation(nil, messageServiceMock)

			resp, err := api.SendMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, resp)
		})
	}
}

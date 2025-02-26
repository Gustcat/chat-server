package tests

import (
	"context"
	"fmt"
	chat "github.com/Gustcat/chat-server/internal/api"
	"github.com/Gustcat/chat-server/internal/service"
	servicemocks "github.com/Gustcat/chat-server/internal/service/mocks"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.DeleteRequest
	}

	var (
		ctx = context.Background()

		id = gofakeit.Int64()

		req = &desc.DeleteRequest{
			Id: id,
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		err             error
		expected        *emptypb.Empty
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				req: req,
				ctx: ctx,
			},
			err:      nil,
			expected: &emptypb.Empty{},
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := servicemocks.NewChatServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
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
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := servicemocks.NewChatServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(serviceErr)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mc := minimock.NewController(t)
			chatServiceMock := tt.chatServiceMock(mc)
			api := chat.NewImplementation(chatServiceMock, nil)

			resp, err := api.Delete(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, resp)
		})
	}
}

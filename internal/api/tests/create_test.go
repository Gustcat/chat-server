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
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type chatServiceMockFunc func(mc *minimock.Controller) service.ChatService

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var usernames []string

	for i := 0; i < 5; i++ {
		usernames = append(usernames, gofakeit.Name())
	}

	var (
		ctx = context.Background()

		id = gofakeit.Int64()

		req = &desc.CreateRequest{
			Usernames: usernames,
		}

		res = &desc.CreateResponse{
			Id: id,
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name            string
		args            args
		expected        *desc.CreateResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				req: req,
				ctx: ctx,
			},
			expected: res,
			err:      nil,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := servicemocks.NewChatServiceMock(mc)
				mock.CreateMock.Expect(ctx, usernames).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				req: req,
				ctx: ctx,
			},
			expected: nil,
			err:      serviceErr,
			chatServiceMock: func(mc *minimock.Controller) service.ChatService {
				mock := servicemocks.NewChatServiceMock(mc)
				mock.CreateMock.Expect(ctx, usernames).Return(0, serviceErr)
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

			resp, err := api.Create(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, resp)
		})
	}
}

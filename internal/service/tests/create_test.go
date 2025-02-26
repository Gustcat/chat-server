package tests

import (
	"context"
	"fmt"
	"github.com/Gustcat/chat-server/internal/repository"
	repomocks "github.com/Gustcat/chat-server/internal/repository/mocks"
	chatservice "github.com/Gustcat/chat-server/internal/service/chat"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type chatRepositoryMockFunc func(mc *minimock.Controller) repository.ChatRepository

	type args struct {
		ctx       context.Context
		usernames []string
	}

	var usernames []string

	for i := 0; i < 5; i++ {
		usernames = append(usernames, gofakeit.Name())
	}

	var (
		ctx = context.Background()

		id = gofakeit.Int64()

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name               string
		args               args
		expected           int64
		err                error
		chatRepositoryMock chatRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				usernames: usernames,
				ctx:       ctx,
			},
			expected: id,
			err:      nil,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repomocks.NewChatRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, usernames).Return(id, nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				usernames: usernames,
				ctx:       ctx,
			},
			expected: 0,
			err:      serviceErr,
			chatRepositoryMock: func(mc *minimock.Controller) repository.ChatRepository {
				mock := repomocks.NewChatRepositoryMock(mc)
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
			chatRepositoryMock := tt.chatRepositoryMock(mc)
			api := chatservice.NewMockService(chatRepositoryMock)

			resp, err := api.Create(tt.args.ctx, tt.args.usernames)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, resp)
		})
	}
}

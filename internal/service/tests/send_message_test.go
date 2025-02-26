package tests

import (
	"context"
	"fmt"
	servicemodel "github.com/Gustcat/chat-server/internal/model"
	"github.com/Gustcat/chat-server/internal/repository"
	repomocks "github.com/Gustcat/chat-server/internal/repository/mocks"
	messageservice "github.com/Gustcat/chat-server/internal/service/message"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type messageRepositoryMockFunc func(mc *minimock.Controller) repository.MessageRepository

	type args struct {
		ctx     context.Context
		message *servicemodel.Message
	}

	var (
		ctx = context.Background()

		username  = gofakeit.Name()
		text      = gofakeit.Sentence(5)
		createdAt = gofakeit.Date()
		chatId    = gofakeit.Int64()

		message = &servicemodel.Message{
			Username:  username,
			Text:      text,
			CreatedAt: createdAt,
			ChatID:    chatId,
		}

		serviceErr = fmt.Errorf("service error")
	)

	tests := []struct {
		name                  string
		args                  args
		err                   error
		messageRepositoryMock messageRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				message: message,
				ctx:     ctx,
			},
			err: nil,
			messageRepositoryMock: func(mc *minimock.Controller) repository.MessageRepository {
				mock := repomocks.NewMessageRepositoryMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(nil)
				return mock
			},
		},
		{
			name: "service error case",
			args: args{
				message: message,
				ctx:     ctx,
			},
			err: serviceErr,
			messageRepositoryMock: func(mc *minimock.Controller) repository.MessageRepository {
				mock := repomocks.NewMessageRepositoryMock(mc)
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
			messageRepositoryMock := tt.messageRepositoryMock(mc)
			api := messageservice.NewMockService(messageRepositoryMock)

			err := api.SendMessage(tt.args.ctx, tt.args.message)
			require.Equal(t, tt.err, err)
		})
	}
}

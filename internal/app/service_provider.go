package app

import (
	"context"
	"github.com/Gustcat/chat-server/internal/api"
	"github.com/Gustcat/chat-server/internal/closer"
	"github.com/Gustcat/chat-server/internal/config"
	"github.com/Gustcat/chat-server/internal/repository"
	"github.com/Gustcat/chat-server/internal/repository/chat"
	"github.com/Gustcat/chat-server/internal/repository/message"
	"github.com/Gustcat/chat-server/internal/service"
	servicechat "github.com/Gustcat/chat-server/internal/service/chat"
	servicemessage "github.com/Gustcat/chat-server/internal/service/message"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	pgPool            *pgxpool.Pool
	chatRepository    repository.ChatRepository
	messageRepository repository.MessageRepository

	chatService    service.ChatService
	messageService service.MessageService

	impl *api.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}
		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) PgPool(ctx context.Context) *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to db: %s", err.Error())
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.pgPool = pool
	}

	return s.pgPool
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chat.NewChatRepository(s.PgPool(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) MessageRepository(ctx context.Context) repository.MessageRepository {
	if s.messageRepository == nil {
		s.messageRepository = message.NewMessageRepository(s.PgPool(ctx))
	}
	return s.messageRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = servicechat.NewChatService(s.ChatRepository(ctx))
	}
	return s.chatService
}

func (s *serviceProvider) MessageService(ctx context.Context) service.MessageService {
	if s.messageService == nil {
		s.messageService = servicemessage.NewMessageService(s.MessageRepository(ctx))
	}
	return s.messageService
}

func (s *serviceProvider) Impl(ctx context.Context) *api.Implementation {
	if s.impl == nil {
		s.impl = api.NewImplementation(s.ChatService(ctx), s.MessageService(ctx))
	}
	return s.impl
}

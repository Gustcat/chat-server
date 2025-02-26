package app

import (
	"context"
	"github.com/Gustcat/chat-server/internal/api"
	db "github.com/Gustcat/chat-server/internal/client"
	"github.com/Gustcat/chat-server/internal/client/db/pg"
	"github.com/Gustcat/chat-server/internal/client/db/transaction"
	"github.com/Gustcat/chat-server/internal/closer"
	"github.com/Gustcat/chat-server/internal/config"
	"github.com/Gustcat/chat-server/internal/repository"
	"github.com/Gustcat/chat-server/internal/repository/chat"
	"github.com/Gustcat/chat-server/internal/repository/message"
	"github.com/Gustcat/chat-server/internal/service"
	servicechat "github.com/Gustcat/chat-server/internal/service/chat"
	servicemessage "github.com/Gustcat/chat-server/internal/service/message"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient          db.Client
	txManager         db.TxManager
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

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to db: %s", err.Error())
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) ChatRepository(ctx context.Context) repository.ChatRepository {
	if s.chatRepository == nil {
		s.chatRepository = chat.NewChatRepository(s.DBClient(ctx))
	}

	return s.chatRepository
}

func (s *serviceProvider) MessageRepository(ctx context.Context) repository.MessageRepository {
	if s.messageRepository == nil {
		s.messageRepository = message.NewMessageRepository(s.DBClient(ctx))
	}
	return s.messageRepository
}

func (s *serviceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = servicechat.NewChatService(s.ChatRepository(ctx), s.TxManager(ctx))
	}
	return s.chatService
}

func (s *serviceProvider) MessageService(ctx context.Context) service.MessageService {
	if s.messageService == nil {
		s.messageService = servicemessage.NewMessageService(s.MessageRepository(ctx), s.TxManager(ctx))
	}
	return s.messageService
}

func (s *serviceProvider) Impl(ctx context.Context) *api.Implementation {
	if s.impl == nil {
		s.impl = api.NewImplementation(s.ChatService(ctx), s.MessageService(ctx))
	}
	return s.impl
}

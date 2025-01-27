package main

import (
	"context"
	"github.com/Gustcat/chat-server/internal/api"
	"github.com/Gustcat/chat-server/internal/config"
	"github.com/Gustcat/chat-server/internal/repository/chat"
	"github.com/Gustcat/chat-server/internal/repository/message"
	servicechat "github.com/Gustcat/chat-server/internal/service/chat"
	servicemessage "github.com/Gustcat/chat-server/internal/service/message"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	ctx := context.Background()

	err := config.Load(".env")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	chatRepo := chat.NewChatRepository(pool)
	chatServ := servicechat.NewChatService(chatRepo)

	messageRepo := message.NewMessageRepository(pool)
	messageServ := servicemessage.NewMessageService(messageRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, api.NewImplementation(chatServ, messageServ))

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

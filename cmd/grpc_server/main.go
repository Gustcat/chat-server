package main

import (
	"context"
	"flag"
	"github.com/Gustcat/chat-server/internal/config"
	desc "github.com/Gustcat/chat-server/pkg/chat_v1"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"time"
)

type server struct {
	desc.UnimplementedChatV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	builderInsertChat := sq.Insert("chat").
		Columns("created_at").
		Values(sq.Expr("DEFAULT")).
		Suffix("RETURNING id")

	query, args, err := builderInsertChat.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to build query: %v", err)
	}

	var chatID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&chatID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to insert chat: %v", err)
	}

	for _, username := range req.GetUsernames() {
		builderInsertUsername := sq.Insert("chat_users").
			PlaceholderFormat(sq.Dollar).
			Columns("chat_id", "username").
			Values(chatID, username)

		query, args, err := builderInsertUsername.ToSql()
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to build query: %v", err)
		}

		_, err = s.pool.Exec(ctx, query, args...)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "failed to insert username %s: %v", username, err)
		}
	}

	return &desc.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	builderDelete := sq.Delete("chat").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()})

	query, args, err := builderDelete.ToSql()

	ct, err := s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "не удалось выполнить SQL-запрос: %v", err)
	}

	if ct.RowsAffected() == 0 {
		return nil, status.Errorf(codes.NotFound, "запись с id %d не найдена", req.GetId())
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	message := req.GetMessage()
	timestamp := message.GetTimestamp()
	fakeCreatedAt := time.Unix(timestamp.GetSeconds(), int64(timestamp.GetNanos()))

	builderSendMessage := sq.Insert("message").
		PlaceholderFormat(sq.Dollar).
		Columns("text", "username", "created_at", "chat_id").
		Values(message.GetText(), message.GetFrom(), fakeCreatedAt, 1)

	query, args, err := builderSendMessage.ToSql()

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "не удалось выполнить SQL-запрос: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func main() {
	flag.Parse()
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

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterChatV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

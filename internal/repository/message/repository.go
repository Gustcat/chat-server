package message

import (
	"context"
	db "github.com/Gustcat/chat-server/internal/client"
	"github.com/Gustcat/chat-server/internal/model"
	"github.com/Gustcat/chat-server/internal/repository"
	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	tableName = "message"

	idColumn        = "id"
	textColumn      = "text"
	usernameColumn  = "username"
	chatIdColumn    = "chat_id"
	createdAtColumn = "created_at"
)

type repo struct {
	db db.Client
}

func NewMessageRepository(db db.Client) repository.MessageRepository {
	return &repo{db: db}
}

func (r *repo) SendMessage(ctx context.Context, message *model.Message) error {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usernameColumn, textColumn, chatIdColumn, createdAtColumn).
		Values(message.Username, message.Text, message.ChatID, message.CreatedAt)

	query, args, err := builder.ToSql()

	q := db.Query{
		Name:     "message_repositoty.SendMeassage",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "не удалось выполнить SQL-запрос: %v", err)
	}

	return nil
}

package message

import (
	"context"
	"github.com/Gustcat/chat-server/internal/repository"
	repomodel "github.com/Gustcat/chat-server/internal/repository/message/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
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
	db *pgxpool.Pool
}

func NewMessageRepository(db *pgxpool.Pool) repository.MessageRepository {
	return &repo{db: db}
}

func (r *repo) SendMessage(ctx context.Context, message *repomodel.Message) error {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(usernameColumn, textColumn, chatIdColumn, createdAtColumn).
		Values(message.Username, message.Text, message.ChatID, message.CreatedAt)

	query, args, err := builder.ToSql()

	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "не удалось выполнить SQL-запрос: %v", err)
	}

	return nil
}

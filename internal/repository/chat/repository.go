package chat

import (
	"context"
	"fmt"
	db "github.com/Gustcat/chat-server/internal/client"
	"github.com/Gustcat/chat-server/internal/repository"
	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	tableName       = "chat"
	idColumn        = "id"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"

	subTableName       = "chat_users"
	chatIdColumn       = "chat_id"
	usernameColumn     = "username"
	subCreatedAtColumn = "created_at"
)

type repo struct {
	db db.Client
}

func NewChatRepository(db db.Client) repository.ChatRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, usernames []string) (int64, error) {
	builderChat := sq.Expr(fmt.Sprintf("INSERT INTO %s DEFAULT VALUES RETURNING %s", tableName, idColumn))

	query, args, err := builderChat.ToSql()
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "chat_repositoty.Create",
		QueryRaw: query,
	}

	var chatID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&chatID)
	if err != nil {
		return 0, status.Errorf(codes.InvalidArgument, "failed to insert chat: %v", err)
	}

	for _, username := range usernames {
		builderInsertUsername := sq.Insert(subTableName).
			PlaceholderFormat(sq.Dollar).
			Columns(chatIdColumn, usernameColumn).
			Values(chatID, username)

		query, args, err := builderInsertUsername.ToSql()
		if err != nil {
			return 0, status.Errorf(codes.InvalidArgument, "failed to build query: %v", err)
		}

		q = db.Query{
			Name:     "chat_repositoty.LinkUsers",
			QueryRaw: query,
		}

		_, err = r.db.DB().ExecContext(ctx, q, args...)
		if err != nil {
			return 0, status.Errorf(codes.InvalidArgument, "failed to insert username %s: %v", username, err)
		}
	}

	return chatID, nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()

	q := db.Query{
		Name:     "chat_repositoty.Delete",
		QueryRaw: query,
	}

	ct, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "не удалось выполнить SQL-запрос: %v", err)
	}

	if ct.RowsAffected() == 0 {
		return status.Errorf(codes.NotFound, "запись с id %d не найдена", id)
	}

	return nil
}

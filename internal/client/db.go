package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Client клиент для работы с БД
type Client interface {
	DB() DB
	Close() error
}

// Query обертка над запросом, хранящая имя запроса и сам запрос
// Имя запроса используется для логированияч и потенциально может использоваться, например, для трейсинга
type Query struct {
	Name     string
	QueryRaw string
}
type Handler func(ctx context.Context) error

// TxManager менеджер транзакций, который выполняет указанный пользователем обработчик в транзакции
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Transactor интерфейс для работы с транзакциями
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// SQLExecer комбинирует NamedExecer и QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// QueryExecer интерфейс для работы с обычными запросами
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// NamedExecer интерфейс для работы с именованными запросами с помощью тегов в структурах
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// Pinger интерфейс для проверки соединения с БД
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB интерфейс для работы с БД
type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}

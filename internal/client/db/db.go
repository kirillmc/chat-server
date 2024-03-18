package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Client - клиент для работы с БД
type Client interface {
	DB() DB
	Close() error
}

// Query - обертка над запросом, хранящая имя запроса и сам запрос
// Имя запроса ипользуется для логгирования и потенциально может использоваться еще где-то,
// например для трейсинга
type Query struct {
	Name     string
	QueryRaw string
}

// SQLExecer -  комбинирует NamedExecer и QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer - интерфейс для работы с именованнымми запросами с помощью тегов в структурах
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecer - интерфейс для работы с обычными запросами - по сути обертка над pgx
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger - интерфейс для проверки соединения с БД - обертка над Ping()
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB - композиция интерфейсов - интерфейс для работы с БД
type DB interface {
	SQLExecer
	Pinger
	Close()
}

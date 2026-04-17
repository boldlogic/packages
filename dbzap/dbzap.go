package dbzap

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
	"go.uber.org/zap"
)

var (
	ErrDBOpen = errors.New("не удалось открыть подключение к БД")
	ErrDBPing = errors.New("не удалось проверить подключение к БД")
)

var sqlOpen = sql.Open

// Pool объединяет подключение к базе данных и логгер,
// используемый при работе с ним.
type Pool struct {
	Db     *sql.DB
	Logger *zap.Logger
}

// New создаёт новое подключение к SQL Server по переданному DSN и
// возвращает готовый пул соединений.
//
// Функция проверяет доступность базы через `PingContext`. Если подключение
// установить не удалось, возвращается ошибка.
func New(ctx context.Context, dsn string, logger *zap.Logger) (*Pool, error) {
	db, err := openDB(ctx, dsn)
	if err != nil {
		if logger != nil {
			logger.Error("не удалось подключиться к БД", zap.Error(err))
		}
		return nil, err
	}
	return &Pool{Db: db, Logger: logger}, nil
}

// openDB открывает соединение с SQL Server и проверяет его доступность
// через `PingContext`.
func openDB(ctx context.Context, dsn string) (*sql.DB, error) {
	conn, err := sqlOpen("sqlserver", dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrDBOpen, err)
	}
	if err := conn.PingContext(ctx); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("%w: %w", ErrDBPing, err)
	}
	return conn, nil
}

// Close закрывает подключение к БД.
func (p *Pool) Close() {
	p.Db.Close()
}

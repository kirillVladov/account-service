package tx_manager

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextTxKey struct{}

// QueryExecutor is the minimal contract repositories need
// to execute queries either via pool or inside a transaction.
type QueryExecutor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) *TxManager {
	return &TxManager{db: db}
}

// WithinTransaction executes callback within a single DB transaction.
// Created transaction is stored in context for repository methods.
func (m *TxManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		if rec := recover(); rec != nil {
			_ = tx.Rollback(ctx)
			panic(rec)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	ctx = context.WithValue(ctx, contextTxKey{}, tx)
	if err = fn(ctx); err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func TxFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(contextTxKey{}).(pgx.Tx)
	return tx, ok
}

func ExecutorFromContext(ctx context.Context, db *pgxpool.Pool) QueryExecutor {
	if tx, ok := TxFromContext(ctx); ok {
		return tx
	}

	return db
}

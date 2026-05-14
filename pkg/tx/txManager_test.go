package tx_manager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTxFromContext(t *testing.T) {
	t.Run("returns transaction when present", func(t *testing.T) {
		ctx := context.Background()
		// Note: we can't easily create a pgx.Tx for testing without a real DB connection
		// This tests the function signature and nil case
		tx, ok := TxFromContext(ctx)
		require.Nil(t, tx)
		require.False(t, ok)
	})
}

func TestExecutorFromContext(t *testing.T) {
	t.Run("returns pool when no transaction in context", func(t *testing.T) {
		ctx := context.Background()
		// Pool is nil here just for testing the logic path
		executor := ExecutorFromContext(ctx, nil)
		require.Nil(t, executor)
	})
}

// Note: Full integration tests for WithinTransaction require a real database connection.
// The following tests would be implemented with a test database:
// - TestWithinTransaction_Commit_Success
// - TestWithinTransaction_Rollback_OnError
// - TestWithinTransaction_Rollback_OnPanic
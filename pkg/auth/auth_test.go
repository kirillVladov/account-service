package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithUserID_GetUserIDFromCtx(t *testing.T) {
	t.Run("stores and retrieves user id", func(t *testing.T) {
		ctx := context.Background()
		userID := uint64(123)

		ctx = WithUserID(ctx, userID)
		retrievedID, ok := GetUserIDFromCtx(ctx)

		require.True(t, ok)
		require.Equal(t, userID, retrievedID)
	})

	t.Run("returns false for empty context", func(t *testing.T) {
		ctx := context.Background()
		userID, ok := GetUserIDFromCtx(ctx)

		require.False(t, ok)
		require.Equal(t, uint64(0), userID)
	})

	t.Run("returns false for zero user id", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithUserID(ctx, 0)
		userID, ok := GetUserIDFromCtx(ctx)

		require.False(t, ok)
		require.Equal(t, uint64(0), userID)
	})

	t.Run("overwrites existing user id", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithUserID(ctx, 123)
		ctx = WithUserID(ctx, 456)

		userID, ok := GetUserIDFromCtx(ctx)

		require.True(t, ok)
		require.Equal(t, uint64(456), userID)
	})
}

package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithUserID_GetUserIDFromCtx(t *testing.T) {
	t.Run("stores and retrieves user id", func(t *testing.T) {
		ctx := context.Background()
		token := "test-token"

		ctx = WithToken(ctx, token)
		retrievedToken, ok := GetTokenFromCtx(ctx)

		require.True(t, ok)
		require.Equal(t, token, retrievedToken)
	})

	t.Run("returns false for empty context", func(t *testing.T) {
		ctx := context.Background()
		token, ok := GetTokenFromCtx(ctx)

		require.False(t, ok)
		require.Equal(t, "", token)
	})

	t.Run("returns false for zero user id", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithToken(ctx, "")
		token, ok := GetTokenFromCtx(ctx)

		require.False(t, ok)
		require.Equal(t, "", token)
	})

	t.Run("overwrites existing user id", func(t *testing.T) {
		ctx := context.Background()
		ctx = WithToken(ctx, "test-1")
		ctx = WithToken(ctx, "test-2")

		token, ok := GetTokenFromCtx(ctx)

		require.True(t, ok)
		require.Equal(t, "test-2", token)
	})
}

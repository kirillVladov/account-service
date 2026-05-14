package auth

import "context"

type userIDKey struct{}

// WithUserID stores user id in context.
func WithUserID(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// GetUserIDFromCtx loads user id from context.
// Second return value indicates whether the id is present and non-zero.
func GetUserIDFromCtx(ctx context.Context) (uint64, bool) {
	v := ctx.Value(userIDKey{})
	userID, ok := v.(uint64)
	return userID, ok && userID != 0
}

package auth

import "context"

type TokenKey struct{}

// WithToken stores token in context.
func WithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, TokenKey{}, token)
}

// GetTokenFromCtx loads token from context.
func GetTokenFromCtx(ctx context.Context) (string, bool) {
	v := ctx.Value(TokenKey{})
	token, ok := v.(string)

	return token, ok && token != ""
}

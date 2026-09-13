package confirm_email

import (
	"context"
)

type AccountRepository struct {
}

type Action struct {
}

func New() *Action {
	return &Action{}
}

func (a *Action) Confirm(ctx context.Context, confirmationToken string) error {
	return nil
}

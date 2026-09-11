package send_email_confirmation_email

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID, organizationID int64) (dto.Account, error)
}

type NotificationsGateway interface {
	SendEmailConfirmation(ctx context.Context, organizationID int64, confirmationLink, email string) error
}

type Action struct {
	accountRepository    AccountRepository
	notificationsGateway NotificationsGateway
}

func New(
	accountRepository AccountRepository,
	notificationsGateway NotificationsGateway,
) *Action {
	return &Action{
		accountRepository:    accountRepository,
		notificationsGateway: notificationsGateway,
	}
}

func (a *Action) Send(ctx context.Context, accountID uuid.UUID, organizationID int64) error {
	account, err := a.accountRepository.GetByID(ctx, accountID, organizationID)
	if err != nil {
		return fmt.Errorf("get account by id: %w", err)
	}

	if account.IsConfirmed {
		return nil
	}

	err = a.notificationsGateway.SendEmailConfirmation(ctx, account.OrganizationID, "", account.Email)
	if err != nil {
		return fmt.Errorf("send email confirmation: %w", err)
	}

	return nil
}

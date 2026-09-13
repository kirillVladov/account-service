package queue_email_confirmation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/pkg/token_manager"
)

type AccountRepository interface {
	GetByID(ctx context.Context, id uuid.UUID, organizationID int64) (dto.Account, error)
}

type TokensRepository interface {
	CreateConfirmationToken(ctx context.Context, userID uuid.UUID, organizationID int64, tokenHash string, expiresAt time.Time) error
}

type Producer interface {
	ProduceAccountConfirmationEvent(ctx context.Context, accountID uuid.UUID, organizationID int64) error
}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Action struct {
	accountRepo AccountRepository
	tokensRepo  TokensRepository
	producer    Producer
	txManager   TxManager
	tokenTTL    time.Duration
}

func New(
	accountRepo AccountRepository,
	tokensRepo TokensRepository,
	producer Producer,
	txManager TxManager,
	tokenTTL time.Duration,
) *Action {
	return &Action{
		accountRepo: accountRepo,
		tokensRepo:  tokensRepo,
		producer:    producer,
		txManager:   txManager,
		tokenTTL:    tokenTTL,
	}
}

func (a *Action) Queue(ctx context.Context, accountID uuid.UUID, organizationID int64) error {
	return a.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		account, err := a.accountRepo.GetByID(ctx, accountID, organizationID)
		if err != nil {
			return fmt.Errorf("get account by id: %w", err)
		}

		if account.IsConfirmed {
			return nil
		}

		token, err := generateToken()
		if err != nil {
			return fmt.Errorf("generate token: %w", err)
		}

		tokenHash := token_manager.Hash(token)
		expiresAt := time.Now().Add(a.tokenTTL)

		if err = a.tokensRepo.CreateConfirmationToken(ctx, account.ID, account.OrganizationID, tokenHash, expiresAt); err != nil {
			return fmt.Errorf("create confirmation token: %w", err)
		}

		if err = a.producer.ProduceAccountConfirmationEvent(ctx, account.ID, account.OrganizationID); err != nil {
			return fmt.Errorf("produce account confirmation event: %w", err)
		}

		return nil
	})
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

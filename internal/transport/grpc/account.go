package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/internal/application/dto/errs"
	pb "github.com/kirillVladov/account-service/internal/gen/grpc"
	"github.com/kirillVladov/account-service/pkg/token_manager"
)

type TokenManager interface {
	ValidateAccess(raw string) (*token_manager.Claims, error)
}

type RefreshTokenAction interface {
	Refresh(ctx context.Context, oldToken, oldRefreshToken string) (string, string, error)
}

type CreateAccountAction interface {
	Do(ctx context.Context, account dto.AccountCreateRequest) (dto.Account, string, string, error)
}

type GetUserAction interface {
	Do(ctx context.Context, id uuid.UUID, organizationID int64) (dto.Account, error)
}

type LoginUserAction interface {
	Do(ctx context.Context, email, password string, organizationID int64) (dto.Account, string, string, error)
}

type ConfirmEmailAction interface {
	Confirm(ctx context.Context, confirmationToken string) error
}

type QueueEmailConfirmationAction interface {
	Queue(ctx context.Context, accountID uuid.UUID, organizationID int64) error
}

type AccountHandlers struct {
	pb.UnimplementedAccountServiceServer

	create             CreateAccountAction
	confirmEmail       ConfirmEmailAction
	get                GetUserAction
	tokenManager       TokenManager
	refreshTokenAction RefreshTokenAction
	login              LoginUserAction
	queueEmailConf     QueueEmailConfirmationAction
}

func NewAccountHandlers(
	create CreateAccountAction,
	get GetUserAction,
	tokenManager TokenManager,
	refreshTokenAction RefreshTokenAction,
	login LoginUserAction,
	confirmEmail ConfirmEmailAction,
	queueEmailConf QueueEmailConfirmationAction,
) *AccountHandlers {
	return &AccountHandlers{
		create:             create,
		get:                get,
		tokenManager:       tokenManager,
		refreshTokenAction: refreshTokenAction,
		login:              login,
		confirmEmail:       confirmEmail,
		queueEmailConf:     queueEmailConf,
	}
}

func (h *AccountHandlers) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountReply, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	if req.GetOrganizationId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "organization_id is invalid")
	}

	if req.GetIdempotencyKey() == "" {
		return nil, status.Error(codes.InvalidArgument, "idempotency_key is null")
	}

	idempotencyKey, err := uuid.Parse(req.GetIdempotencyKey())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "idempotency_key is invalid")
	}

	request := dto.AccountCreateRequest{
		Email:          req.GetEmail(),
		Password:       req.GetPassword(),
		OrganizationID: req.GetOrganizationId(),
		IdempotencyKey: idempotencyKey,
	}

	account, token, refreshToken, err := h.create.Do(ctx, request)
	if err != nil {
		if errors.Is(err, errs.ErrAccountAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "account already exists")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("create account: %v", err))
	}

	return &pb.CreateAccountReply{Account: pbAccountFromDTO(account), RefreshToken: refreshToken, Token: token}, nil
}

func (h *AccountHandlers) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenReply, error) {
	claims, err := h.tokenManager.ValidateAccess(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not valid token")
	}

	return &pb.VerifyTokenReply{AccountId: claims.UserID, OrganizationId: claims.OrganizationID}, nil
}

func (h *AccountHandlers) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenReply, error) {
	token, refreshToken, err := h.refreshTokenAction.Refresh(ctx, req.GetToken(), req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, errs.ErrAccountBlocked) {
			return nil, status.Error(codes.PermissionDenied, "account blocked")
		}

		if errors.Is(err, errs.ErrForbidden) {
			return nil, status.Error(codes.PermissionDenied, "forbidden")
		}

		return nil, status.Error(codes.Internal, "not valid token")
	}

	return &pb.RefreshTokenReply{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AccountHandlers) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountReply, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	if req.GetOrganizationId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "organization_id is invalid")
	}

	account, err := h.get.Do(ctx, id, req.GetOrganizationId())
	if err != nil {
		if errors.Is(err, errs.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, pb.GetAccountRequest_ACCOUNT_NOT_FOUND.Enum().String())
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("get user: %v", err))
	}

	return &pb.GetAccountReply{Account: pbAccountFromDTO(account)}, nil
}

func (h *AccountHandlers) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	if req.GetOrganizationId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "organization_id is invalid")
	}

	account, token, refreshToken, err := h.login.Do(ctx, req.GetEmail(), req.GetPassword(), req.GetOrganizationId())
	if err != nil {
		if errors.Is(err, errs.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}

		if errors.Is(err, errs.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, "invalid credentials")
		}

		if errors.Is(err, errs.ErrAccountBlocked) {
			return nil, status.Error(codes.PermissionDenied, "account blocked")
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("login: %v", err))
	}

	return &pb.LoginReply{Account: pbAccountFromDTO(account), Token: token, RefreshToken: refreshToken}, nil
}

func (h *AccountHandlers) ConfirmEmail(ctx context.Context, req *pb.ConfirmEmailRequest) (*pb.ConfirmEmailReply, error) {
	if err := h.confirmEmail.Confirm(ctx, req.GetConfirmationToken()); err != nil {
		switch {
		case errors.Is(err, errs.ErrAccountNotFound),
			errors.Is(err, errs.ErrTokenNotValid),
			errors.Is(err, errs.ErrAccountBlocked):
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("confirm email: %v", err))
		case errors.Is(err, errs.ErrAccountAlreadyConfirmed):
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("confirm email: %v", err))
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("confirm email: %v", err))
	}

	return &pb.ConfirmEmailReply{}, nil
}

func (h *AccountHandlers) ResendEmailConfirmation(ctx context.Context, req *pb.ResendEmailConfirmationRequest) (*pb.ResendEmailConfirmationReply, error) {
	accountID, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account_id")
	}

	if req.GetOrganizationId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "organization_id is invalid")
	}

	if err = h.queueEmailConf.Queue(ctx, accountID, req.GetOrganizationId()); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("resend email confirmation: %v", err))
	}

	return &pb.ResendEmailConfirmationReply{}, nil
}

func pbAccountFromDTO(a dto.Account) *pb.Account {
	return &pb.Account{
		Id:             a.ID.String(),
		Email:          a.Email,
		OrganizationId: a.OrganizationID,
		Confirmed:      a.IsConfirmed,
		Blocked:        a.IsBlocked,
	}
}

package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	createUserAction "github.com/kirillVladov/account-service/internal/application/action/create_user"
	getByTelegramId "github.com/kirillVladov/account-service/internal/application/action/get_by_telegram_id"
	getUserAction "github.com/kirillVladov/account-service/internal/application/action/get_user"
	"github.com/kirillVladov/account-service/internal/application/dto"
	"github.com/kirillVladov/account-service/internal/application/dto/errs"
	pb "github.com/kirillVladov/account-service/internal/gen/grpc"
	"github.com/kirillVladov/account-service/pkg/token_manager"
)

type TokenManager interface {
	ValidateAccess(raw string) (*token_manager.Claims, error)
}

type RefreshTokenAction interface {
	Refresh(ctx context.Context, token string) (string, string, error)
}

type AccountHandlers struct {
	pb.UnimplementedAccountServiceServer

	create             *createUserAction.CreateUserAction
	get                *getUserAction.GetUserAction
	getByTelegramID    *getByTelegramId.Action
	tokenManager       TokenManager
	refreshTokenAction RefreshTokenAction
}

func NewAccountHandlers(
	create *createUserAction.CreateUserAction,
	get *getUserAction.GetUserAction,
	getByTelegramID *getByTelegramId.Action,
	tokenManager TokenManager,
) *AccountHandlers {
	return &AccountHandlers{
		create:          create,
		get:             get,
		getByTelegramID: getByTelegramID,
		tokenManager:    tokenManager,
	}
}

func (h *AccountHandlers) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountReply, error) {
	if h == nil || h.create == nil {
		return nil, status.Error(codes.Internal, "create user handler not initialized")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	account := dto.Account{
		Email:      req.GetEmail(),
		Name:       req.GetName(),
		TelegramID: req.GetTelegramId(),
		Phone:      req.GetPhone(),
	}

	account, err := h.create.Do(ctx, account)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("create account: %v", err))
	}

	return &pb.CreateAccountReply{Account: pbAccountFromDTO(account)}, nil
}

func (h *AccountHandlers) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenReply, error) {
	claims, err := h.tokenManager.ValidateAccess(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not valid token")
	}

	return &pb.VerifyTokenReply{AccountId: claims.UserID}, nil
}

func (h *AccountHandlers) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenReply, error) {
	token, refreshToken, err := h.refreshTokenAction.Refresh(ctx, req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "not valid token")
	}

	return &pb.RefreshTokenReply{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AccountHandlers) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountReply, error) {
	if h == nil || h.get == nil {
		return nil, status.Error(codes.Internal, "get user handler not initialized")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	var (
		account dto.Account
		err     error
	)

	switch ident := req.GetIdentifier().(type) {
	case *pb.GetAccountRequest_Id:
		id, err := uuid.Parse(ident.Id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid id")
		}

		account, err = h.get.Do(ctx, id)
		if err != nil {
			if errors.Is(err, errs.ErrAccountNotFound) {
				return nil, status.Error(codes.NotFound, pb.GetAccountRequest_ACCOUNT_NOT_FOUND.Enum().String())
			}

			return nil, status.Error(codes.Internal, fmt.Sprintf("get user: %v", err))
		}
	case *pb.GetAccountRequest_TelegramId:
		account, err = h.getByTelegramID.Get(ctx, ident.TelegramId)
		if err != nil {
			if errors.Is(err, errs.ErrAccountNotFound) {
				return nil, status.Error(codes.NotFound, pb.GetAccountRequest_ACCOUNT_NOT_FOUND.Enum().String())
			}

			return nil, status.Error(codes.Internal, fmt.Sprintf("get user by tg id: %v", err))
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "identifier is required")
	}

	return &pb.GetAccountReply{Account: pbAccountFromDTO(account)}, nil
}

func pbAccountFromDTO(a dto.Account) *pb.Account {
	return &pb.Account{
		Id:         a.ID.String(),
		Email:      a.Email,
		Name:       a.Name,
		TelegramId: a.TelegramID,
		Phone:      a.Phone,
	}
}

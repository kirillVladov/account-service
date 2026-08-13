package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	Refresh(ctx context.Context, oldToken, oldRefreshToken string) (string, string, error)
}

type CreateAccountAction interface {
	Do(ctx context.Context, account dto.AccountCreateRequest) (dto.Account, string, string, error)
}

type AccountHandlers struct {
	pb.UnimplementedAccountServiceServer

	create             CreateAccountAction
	get                *getUserAction.GetUserAction
	tokenManager       TokenManager
	refreshTokenAction RefreshTokenAction
}

func NewAccountHandlers(
	create CreateAccountAction,
	get *getUserAction.GetUserAction,
	tokenManager TokenManager,
	refreshTokenAction RefreshTokenAction,
) *AccountHandlers {
	return &AccountHandlers{
		create:             create,
		get:                get,
		tokenManager:       tokenManager,
		refreshTokenAction: refreshTokenAction,
	}
}

func (h *AccountHandlers) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountReply, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	request := dto.AccountCreateRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	account, token, refreshToken, err := h.create.Do(ctx, request)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("create account: %v", err))
	}

	return &pb.CreateAccountReply{Account: pbAccountFromDTO(account), RefreshToken: refreshToken, Token: token}, nil
}

func (h *AccountHandlers) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenReply, error) {
	claims, err := h.tokenManager.ValidateAccess(req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "not valid token")
	}

	return &pb.VerifyTokenReply{AccountId: claims.UserID}, nil
}

func (h *AccountHandlers) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenReply, error) {
	token, refreshToken, err := h.refreshTokenAction.Refresh(ctx, req.GetToken(), req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "not valid token")
	}

	return &pb.RefreshTokenReply{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (h *AccountHandlers) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountReply, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	account, err := h.get.Do(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, pb.GetAccountRequest_ACCOUNT_NOT_FOUND.Enum().String())
		}

		return nil, status.Error(codes.Internal, fmt.Sprintf("get user: %v", err))
	}

	return &pb.GetAccountReply{Account: pbAccountFromDTO(account)}, nil
}

func pbAccountFromDTO(a dto.Account) *pb.Account {
	return &pb.Account{
		Id:    a.ID.String(),
		Email: a.Email,
	}
}

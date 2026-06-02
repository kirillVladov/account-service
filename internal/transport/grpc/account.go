package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	createUserAction "github.com/kirillVladov/account-service/internal/application/action/create_user"
	getByTelegramId "github.com/kirillVladov/account-service/internal/application/action/get_by_telegram_id"
	getUserAction "github.com/kirillVladov/account-service/internal/application/action/get_user"
	"github.com/kirillVladov/account-service/internal/application/dto"
	pb "github.com/kirillVladov/account-service/internal/gen/grpc"
)

type AccountHandlers struct {
	pb.UnimplementedAccountServiceServer

	create          *createUserAction.CreateUserAction
	get             *getUserAction.GetUserAction
	getByTelegramID *getByTelegramId.Action
}

func NewAccountHandlers(
	create *createUserAction.CreateUserAction,
	get *getUserAction.GetUserAction,
	getByTelegramID *getByTelegramId.Action,
) *AccountHandlers {
	return &AccountHandlers{create: create, get: get, getByTelegramID: getByTelegramID}
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
			return nil, status.Error(codes.Internal, fmt.Sprintf("get user: %v", err))
		}
	case *pb.GetAccountRequest_TelegramId:
		account, err = h.getByTelegramID.Get(ctx, ident.TelegramId)
		if err != nil {
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

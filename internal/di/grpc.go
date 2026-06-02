package di

import "github.com/kirillVladov/account-service/internal/transport/grpc"

func (di *DI) AccountHandler() *grpc.AccountHandlers {
	return grpc.NewAccountHandlers(
		di.CreateUserAction(),
		di.GetUserAction(),
		di.GetByTelegramIDAction(),
	)
}

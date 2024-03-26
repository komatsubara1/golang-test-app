//go:build wireinject
// +build wireinject

package main

import (
	master_repository "app/infrastructure/repository/master"
	user_repository "app/infrastructure/repository/user"
	"app/presenter/http/handler"
	"app/use_case"
	"github.com/google/wire"
)

func InitUserHandler() *handler.UserHandler {
	wire.Build(
		user_repository.NewUserRepository,
		user_repository.NewUserAuthRepository,
		use_case.NewUserUseCase,
		handler.NewUserHandler,
	)
	return nil
}

func InitItemHandler() *handler.ItemHandler {
	wire.Build(
		user_repository.NewUserRepository,
		user_repository.NewUserItemRepository,
		user_repository.NewUserPresentRepository,
		master_repository.NewItemMasterRepository,
		use_case.NewItemUseCase,
		handler.NewItemHandler,
	)
	return nil
}

func InitDebugHandler() *handler.DebugHandler {
	wire.Build(
		use_case.NewDebugUseCase,
		handler.NewDebugHandler,
	)
	return nil
}

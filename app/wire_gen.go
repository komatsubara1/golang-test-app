// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"app/infrastructure/repository/master"
	"app/infrastructure/repository/user"
	"app/presenter/http/handler"
	"app/use_case"
)

// Injectors from wire.go:

func InitUserHandler() *handler.UserHandler {
	userRepository := user.NewUserRepository()
	userAuthRepository := user.NewUserAuthRepository()
	userUseCase := use_case.NewUserUseCase(userRepository, userAuthRepository)
	userHandler := handler.NewUserHandler(userUseCase)
	return userHandler
}

func InitItemHandler() *handler.ItemHandler {
	userRepository := user.NewUserRepository()
	userItemRepository := user.NewUserItemRepository()
	userPresentRepository := user.NewUserPresentRepository()
	itemMasterRepository := master.NewItemMasterRepository()
	itemUseCase := use_case.NewItemUseCase(userRepository, userItemRepository, userPresentRepository, itemMasterRepository)
	itemHandler := handler.NewItemHandler(itemUseCase)
	return itemHandler
}

func InitDebugHandler() *handler.DebugHandler {
	debugUseCase := use_case.NewDebugUseCase()
	debugHandler := handler.NewDebugHandler(debugUseCase)
	return debugHandler
}

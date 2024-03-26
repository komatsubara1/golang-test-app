package main

import (
	"app/presenter/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

func init() {
	loadEnv()
}

// envファイル読み込み
func loadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		slog.Error("failed load .env.", err)
		return
	}
}

// @title test-app
// @version 1.0
// @license.name komatsubara.s
// @description Golangサンプルプロジェクト
func main() {
	r := gin.Default()

	SetupRoute(r)

	_ = r.Run()
}

// SetupRoute ルーティング設定
func SetupRoute(r *gin.Engine) *gin.Engine {
	// user
	uh := InitUserHandler()
	user := r.Group("/user")
	user.Use(middleware.SetUpMiddleware(), middleware.TimeShiftMiddleware())
	{
		user.POST("/create", uh.Create)
	}
	user.Use(middleware.TimeShiftMiddleware())
	{
		user.POST("/login", uh.Login)
	}

	userAuth := r.Group("/user")
	userAuth.Use(middleware.SetUpMiddleware(), middleware.JwtAuthenticationMiddleware(), middleware.TimeShiftMiddleware())
	{
		userAuth.POST("/get", uh.Get)
	}

	// item
	ih := InitItemHandler()
	item := r.Group("/item")
	item.Use(middleware.SetUpMiddleware(), middleware.JwtAuthenticationMiddleware(), middleware.TimeShiftMiddleware())
	{
		item.POST("/get", ih.Get)
		item.POST("/get_all", ih.GetAll)
		item.POST("/gain", ih.Gain)
		item.POST("/use", ih.Use)
		item.POST("/sell", ih.Sell)
	}

	// debug
	dh := InitDebugHandler()
	debug := r.Group("/debug")
	debug.Use(middleware.SetUpMiddleware(), middleware.JwtAuthenticationMiddleware(), middleware.TimeShiftMiddleware())
	{
		debug.POST("/set_time_shift", dh.SetTimeShift)
		debug.POST("/send_present", dh.SendPresent)
		/*
			debug.POST("/gain_item", dh.GainItem)
		*/
	}

	return r
}

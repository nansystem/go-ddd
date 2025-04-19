package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/nansystem/go-ddd/internal/config"
	"github.com/nansystem/go-ddd/internal/infrastructure/mysql"
	"github.com/nansystem/go-ddd/internal/presentation"
	"github.com/nansystem/go-ddd/internal/usecase"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("設定の読み込みに失敗しました: %v", err)
	}

	db, err := mysql.NewConnection(cfg.DBConfig)
	if err != nil {
		log.Fatalf("MySQLへの接続に失敗しました: %v", err)
	}
	defer db.Close()

	userRepository := mysql.NewUserRepository(db)
	userService := usecase.NewUserService(userRepository)

	e := presentation.NewRouter()
	setupRoutes(e, userService)

	err = e.Start(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRoutes(e *echo.Echo, userService *usecase.UserService) {
	// FIXME グループ追加のたびにmain.goが膨らんでしまわないようにする
	userHandler := presentation.NewUserHandler(userService)
	userHandler.SetupUserRoutes(e.Group("/users"))
}

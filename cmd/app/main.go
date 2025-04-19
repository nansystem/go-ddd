package main

import (
	"log"

	"github.com/nansystem/go-ddd/internal/infrastructure/mysql"
	"github.com/nansystem/go-ddd/internal/presentation"
	"github.com/nansystem/go-ddd/internal/usecase"
)

func main() {
	// FIXME 環境変数から読み込む
	// FIXME main.goにロジックを追加しないようにする
	db, err := mysql.NewConnection(mysql.DBConfig{
		User:     "ddduser",
		Password: "dddpass",
		Host:     "localhost",
		Port:     "13306",
		DBName:   "go_ddd",
	})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	userRepository := mysql.NewUserRepository(db)
	userService := usecase.NewUserService(userRepository)

	e := presentation.NewRouter()

	// FIXME グループ追加のたびにmain.goが膨らんでしまわないようにする
	userHandler := presentation.NewUserHandler(userService)
	userHandler.SetupUserRoutes(e.Group("/users"))

	err = e.Start(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

package main

import (
	"log"

	"github.com/nansystem/go-ddd/internal/presentation"
)

func main() {
	e := presentation.NewRouter()

	// FIXME グループ追加のたびにmain.goが膨らんでしまわないようにする
	userHandler := presentation.NewUserHandler()
	userHandler.SetupUserRoutes(e.Group("/users"))

	err := e.Start(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

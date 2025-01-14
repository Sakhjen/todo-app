package main

import (
	"log"

	"github.com/Sakhjen/todo-app"
	"github.com/Sakhjen/todo-app/pkg/handler"
)

func main() {
	handlers := new(handler.Handler)
	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running server: %s", err.Error())

	}

}

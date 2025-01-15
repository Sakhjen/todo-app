package main

import (
	"log"

	"github.com/Sakhjen/todo-app"
	"github.com/Sakhjen/todo-app/pkg/handler"
	"github.com/Sakhjen/todo-app/pkg/repository"
	"github.com/Sakhjen/todo-app/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("Error occured while running server: %s", err.Error())

	}

}

package service

import (
	"github.com/Sakhjen/todo-app"
	"github.com/Sakhjen/todo-app/pkg/repository"
)

type todoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *todoListService {
	return &todoListService{
		repo: repo,
	}
}

func (s *todoListService) Create(userId int, list todo.TodoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *todoListService) GetAll(userId int) ([]todo.TodoList, error) {
	return s.repo.GetAll(userId)
}

func (s *todoListService) GetById(userId int, id int) (todo.TodoList, error) {
	return s.repo.GetById(userId, id)
}

func (s *todoListService) Delete(userId, id int) error {
	return s.repo.Delete(userId, id)
}

func (s *todoListService) Update(userId, id int, list todo.ListUpdateInput) error {
	return s.repo.Update(userId, id, list)
}

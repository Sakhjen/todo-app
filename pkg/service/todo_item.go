package service

import (
	"github.com/Sakhjen/todo-app"
	"github.com/Sakhjen/todo-app/pkg/repository"
)

type todoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList) *todoItemService {
	return &todoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *todoItemService) Create(userId, listId int, item todo.TodoItem) (int, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(listId, item)
}

func (s *todoItemService) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	return s.repo.GetAll(userId, listId)
}

func (s *todoItemService) GetById(userId, itemId int) (todo.TodoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *todoItemService) Update(userId, id int, input todo.ItemUpdateInput) error {
	return s.repo.Update(userId, id, input)
}

func (s *todoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

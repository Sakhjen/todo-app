package repository

import (
	"fmt"
	"strings"

	"github.com/Sakhjen/todo-app"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, todoList todo.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tx.QueryRow(createListQuery, todoList.Title, todoList.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserListsQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tx.Exec(createUserListsQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId int) ([]todo.TodoList, error) {
	var lists []todo.TodoList

	query := fmt.Sprintf("SELECT l.* FROM %s l JOIN %s ul ON l.id = ul.list_id WHERE ul.user_id = $1",
		todoListsTable, usersListsTable)

	err := r.db.Select(&lists, query, userId)

	return lists, err

}

func (r *TodoListPostgres) GetById(userId int, id int) (todo.TodoList, error) {
	var list todo.TodoList

	query := fmt.Sprintf("SELECT l.* FROM %s l JOIN %s ul ON l.id = ul.list_id WHERE ul.user_id = $1 AND l.id = $2",
		todoListsTable, usersListsTable)

	err := r.db.Get(&list, query, userId, id)

	return list, err
}

func (r *TodoListPostgres) Delete(userId int, id int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE ul.list_id = tl.id AND ul.user_id = $1 AND tl.id = $2", todoListsTable, usersListsTable)
	_, err := r.db.Exec(query, userId, id)
	return err
}

func (r *TodoListPostgres) Update(userId, id int, list todo.ListUpdateInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if list.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, list.Title)
		argId++
	}

	if list.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, list.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE ul.list_id = tl.id AND ul.user_id = $%d AND tl.id = $%d",
		todoListsTable, setQuery, usersListsTable, argId, argId+1)

	args = append(args, userId, id)
	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}

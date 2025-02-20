package repository

import (
	"fmt"
	"strings"

	"github.com/Sakhjen/todo-app"
	"github.com/jmoiron/sqlx"
)

type todoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *todoItemPostgres {
	return &todoItemPostgres{
		db: db,
	}
}

func (r *todoItemPostgres) Create(listId int, item todo.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s(title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemQuery := fmt.Sprintf("INSERT INTO %s(list_id, item_id) VALUES ($1, $2)", listsItemsTable)
	_, err = tx.Exec(createListItemQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *todoItemPostgres) GetAll(userId, listId int) ([]todo.TodoItem, error) {
	var items []todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description from %s ti 
						  JOIN %s li ON ti.id = li.item_id 
						  JOIN %s ul ON ul.list_id = li.list_id 
						  WHERE li.list_id = $1 AND ul.user_id = $2`,
		todoItemsTable,
		listsItemsTable,
		usersListsTable,
	)

	if err := r.db.Select(&items, query, listId, userId); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *todoItemPostgres) GetById(userId, itemId int) (todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done
						  FROM %s ti 
						  JOIN %s li ON ti.id = li.item_id
						  JOIN %s ul ON ul.list_id = li.list_id
						  WHERE ti.id = $1 AND ul.user_id = $2`,
		todoItemsTable,
		listsItemsTable,
		usersListsTable,
	)

	if err := r.db.Get(&item, query, itemId, userId); err != nil {
		return item, err
	}

	return item, nil
}

func (r *todoItemPostgres) Update(userId, id int, item todo.ItemUpdateInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if item.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *item.Title)
		argId++
	}

	if item.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *item.Description)
		argId++
	}

	if item.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *item.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s ti 
						  SET %s 
						  FROM %s li, %s ul
						  WHERE li.list_id = ul.list_id AND li.item_id = ti.id AND ul.user_id = $%d AND ti.id = $%d`,
		todoItemsTable, setQuery, listsItemsTable, usersListsTable, argId, argId+1)

	args = append(args, userId, id)
	fmt.Printf("updateQuery: %s", query)
	fmt.Printf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r todoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf(`DELETE FROM %s ti 
						  USING %s li, %s ul 
						  WHERE li.item_id = ti.id AND ul.list_id = li.list_id AND ul.user_id = $1 AND ti.id = $2`,
		todoItemsTable,
		listsItemsTable,
		usersListsTable,
	)

	_, err := r.db.Exec(query, userId, itemId)
	return err
}

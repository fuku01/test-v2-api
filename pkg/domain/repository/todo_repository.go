package repository

import "github.com/fuku01/test-v2-api/pkg/domain/model"

type TodoRepository interface {
	ListTodos() ([]*model.Todo, error)
	CreateTodo(input *model.CreateTodoInput) (*model.Todo, error)
}

package usecase

import (
	"fmt"

	"github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/domain/repository"
)

type TodoUsecase interface {
	ListTodos() ([]*model.Todo, error)
	CreateTodo(input *model.CreateTodoInput) (*model.Todo, error)
}

type todoUsecase struct {
	tr repository.TodoRepository
}

func NewTodoUsecase(tr repository.TodoRepository) TodoUsecase {
	return &todoUsecase{
		tr: tr,
	}
}

func (u *todoUsecase) ListTodos() ([]*model.Todo, error) {
	todos, err := u.tr.ListTodos()
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (u *todoUsecase) CreateTodo(input *model.CreateTodoInput) (*model.Todo, error) {
	fmt.Println("CreateTodo usecase")
	todo, err := u.tr.CreateTodo(input)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

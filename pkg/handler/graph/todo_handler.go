package graph

import (
	"fmt"
	"strconv"

	domain_model "github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/graph/generated/model"
	"github.com/fuku01/test-v2-api/pkg/usecase"
	"github.com/samber/lo"
)

type TodoHandler interface {
	ListTodos() ([]*model.Todo, error)
}

type todoHandler struct {
	tu usecase.TodoUsecase
}

func NewTodoHandler(tu usecase.TodoUsecase) TodoHandler {
	return &todoHandler{
		tu: tu,
	}
}

func (h *todoHandler) ListTodos() ([]*model.Todo, error) {
	todos, err := h.tu.ListTodos()

	for _, todo := range todos {
		fmt.Printf("%+v\n", todo)
	}

	if err != nil {
		return nil, err
	}

	convTodos := lo.Map(todos, func(todo *domain_model.Todo, _ int) *model.Todo {
		return convTodo(todo)
	})

	return convTodos, nil
}

func convTodo(todo *domain_model.Todo) *model.Todo {
	if todo == nil {
		return nil
	}

	return &model.Todo{
		ID:        strconv.FormatUint(uint64(todo.ID), 10),
		Content:   todo.Content,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
}

package infrastructure

import (
	"fmt"

	"github.com/fuku01/test-v2-api/pkg/domain/model"
	"github.com/fuku01/test-v2-api/pkg/domain/repository"
	"gorm.io/gorm"
)

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &todoRepository{
		db: db,
	}
}

func (r *todoRepository) ListTodos() ([]*model.Todo, error) {

	todos := []*model.Todo{}
	err := r.db.Find(&todos).Error
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *todoRepository) CreateTodo(input *model.CreateTodoInput) (*model.Todo, error) {
	fmt.Println("CreateTodo repository")
	todo := &model.Todo{
		Content: input.Content,
	}

	err := r.db.Create(todo).Error
	if err != nil {
		return nil, err
	}

	return todo, nil
}

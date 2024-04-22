package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

type Todo struct {
	gorm.Model
	Content string
}

type CreateTodoInput struct {
	Content string
}

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
	Title   string
	Content string
}

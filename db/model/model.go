package gorm_model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string
}

type Message struct {
	gorm.Model
	Content string
}

// gormのDeletedAt型はgorm.Modelに含まれているが、ドメインモデルには含めない

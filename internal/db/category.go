package db

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string
}

func NewCategory(name string) *Category {
	return &Category{
		Name: name,
	}
}

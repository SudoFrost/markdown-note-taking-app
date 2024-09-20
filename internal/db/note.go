package db

import (
	"gorm.io/gorm"
)

type Note struct {
	gorm.Model
	Title      string
	Content    string
	CategoryID uint
	Category   *Category
}

func NewNote(title string, content string) *Note {
	return &Note{
		Title:   title,
		Content: content,
	}
}

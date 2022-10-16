package models

import (
	"gorm.io/gorm"
)

type Channel struct {
	gorm.Model
	TgId          int64 `gorm:"uniqueIndex"`
	TopPostsCount int   `gorm:"default:10"`
	Posts         []Post
}

type Post struct {
	gorm.Model
	Url         string `gorm:"uniqueIndex"`
	Description string
}

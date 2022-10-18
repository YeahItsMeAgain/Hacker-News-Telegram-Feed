package models

import (
	"gorm.io/gorm"
)

var (
	FeedTypes = []string{"topstories", "newstories", "beststories"}
)

type Channel struct {
	gorm.Model
	TgId       int64 `gorm:"uniqueIndex"`
	Title      string
	FeedType   string `gorm:"default:topstories"`
	PostsCount int    `gorm:"default:10"`
	Posts      []Post `gorm:"many2many:channels_posts;"`
}

type Post struct {
	gorm.Model
	PostId      int `gorm:"uniqueIndex"`
	Url         string
	Description string
	Channels    []Channel `gorm:"many2many:channels_posts;"`
}

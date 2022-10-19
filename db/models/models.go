package models

import (
	"gorm.io/gorm"
)

var (
	FeedTypes = []string{"topstories", "newstories"}
)

type Channel struct {
	gorm.Model
	TgId                int64 `gorm:"uniqueIndex"`
	Title               string
	FeedType            string     `gorm:"default:topstories"`
	PostsCount          int        `gorm:"default:10"`
	Posts               []*Post    `gorm:"many2many:channels_posts;"`
	WhitelistedKeywords []*Keyword `gorm:"many2many:whitelisted_users;"`
	BlacklistedKeywords []*Keyword `gorm:"many2many:blacklisted_users;"`
}

type Post struct {
	gorm.Model
	PostId   int        `gorm:"uniqueIndex" json:"id"`
	Url      string     `json:"url"`
	Title    string     `json:"title"`
	Channels []*Channel `gorm:"many2many:channels_posts;" json:",omitempty"`
}

type Keyword struct {
	gorm.Model
	Keyword                  string     `gorm:"uniqueIndex"`
	ChannelsWhichWhitelisted []*Channel `gorm:"many2many:whitelisted_users;"`
	ChannelsWhichBlacklisted []*Channel `gorm:"many2many:blacklisted_users;"`
}

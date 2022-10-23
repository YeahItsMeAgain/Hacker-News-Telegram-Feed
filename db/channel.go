package db

import (
	"fmt"

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
	MinimumScore        int        `gorm:"default:1"`
	Posts               []*Post    `gorm:"many2many:channels_posts;"`
	WhitelistedKeywords []*Keyword `gorm:"many2many:whitelisted_users;"`
	BlacklistedKeywords []*Keyword `gorm:"many2many:blacklisted_users;"`
}

func (channel *Channel) FirstOrCreate() {
	DB.FirstOrCreate(channel, "tg_id = ?", channel.TgId)
}

func (channel *Channel) GetAssociatedKeywords(association string) string {
	var dbKeywords []Keyword
	DB.Model(&channel).Association(association).Find(&dbKeywords)
	keywords := ""
	for _, keyword := range dbKeywords {
		keywords += fmt.Sprintf("%s\n", keyword.Keyword)
	}
	return keywords
}

func (channel *Channel) GetAssociatedKeywordsByKeyword(association string, keyword Keyword) []*Keyword {
	var keywords []*Keyword
	DB.Model(&channel).Where("keyword = ?", keyword.Keyword).Association(association).Find(&keywords)
	return keywords
}

package db

import "gorm.io/gorm"

type Keyword struct {
	gorm.Model
	Keyword                  string     `gorm:"uniqueIndex"`
	ChannelsWhichWhitelisted []*Channel `gorm:"many2many:whitelisted_users;"`
	ChannelsWhichBlacklisted []*Channel `gorm:"many2many:blacklisted_users;"`
}

func (keyword *Keyword) FirstOrCreate() {
	DB.FirstOrCreate(keyword, "keyword = ?", keyword.Keyword)
}

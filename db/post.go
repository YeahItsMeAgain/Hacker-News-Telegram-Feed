package db

import (
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	PostId   int        `gorm:"uniqueIndex" json:"id"`
	Url      string     `json:"url"`
	Title    string     `json:"title"`
	Score    int        `gorm:"default:1" json:"score"`
	Channels []*Channel `gorm:"many2many:channels_posts;" json:",omitempty"`
}

func (post *Post) Upsert() {
	DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "post_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{"score": post.Score}),
	}).Create(post)
}

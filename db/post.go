package db

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	PostId   int        `gorm:"uniqueIndex" json:"id"`
	Url      string     `json:"url"`
	Title    string     `json:"title"`
	Channels []*Channel `gorm:"many2many:channels_posts;" json:",omitempty"`
}

func (post *Post) FirstOrCreate() {
	DB.FirstOrCreate(post, "post_id = ?", post.PostId)
}

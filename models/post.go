package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title        string `gorm:"not null"`
	Content      string `gorm:"not null"`
	UserID       uint64
	Comments     []Comment `gorm:"foreignKey:PostID;"`
	CommentCount int
	// 映射查询User表会把用户的信息查不来，只取ID就好
	//User        User `gorm:"foreignKey:UserID;"`
}

package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	UserID  uint64
	PostID  uint64
	Post    Post `gorm:"foreignKey:PostID;"`
	// 映射查询User表会把用户的信息查不来，只取ID就好
	//User    User `gorm:"foreignKey:UserID;"`
}

package models

import "time"

type Comment struct {
	ID        int       `gorm:"primaryKey;comment:'评论唯一标识'" json:"id"`
	PostID    int       `gorm:"comment:'文章ID'" json:"post_id"`
	UserID    int       `gorm:"comment:'评论者ID'" json:"user_id"`
	Content   string    `gorm:"type:varchar(255);comment:'评论内容'" json:"content"`
	ParentID  int       `gorm:"comment:'父评论ID(用于回复）'" json:"parent_id"`
	CreatedAt time.Time `gorm:"comment:'评论时间'" json:"created_at"`
}

package models

import "time"

type Post struct {
	ID         int       `gorm:"primaryKey;comment:'文章ID'" json:"id"`
	UserID     int       `gorm:"column:used_id;comment:'文章所属者'" json:"user_id"`
	Title      string    `gorm:"type:varchar(255);comment:'文章标题'" json:"title"`
	AddTime    time.Time `gorm:"comment:'文章添加时间'" json:"add_time"`
	Content    string    `gorm:"type:text;comment:'文章'" json:"content"`
	Likes      int       `gorm:"comment:'点赞次数'" json:"likes"`
	Browse     int       `gorm:"comment:'浏览次数'" json:"browse"`
	Collection int       `gorm:"comment:'收藏次数'" json:"collection"`
}

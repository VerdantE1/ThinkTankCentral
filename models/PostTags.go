package models

type PostTag struct {
	ID     int `gorm:"primaryKey;comment:'文章与标签的关系标识'" json:"id"`
	PostID int `gorm:"comment:'文章ID'" json:"post_id"`
	TagID  int `gorm:"comment:'标签ID'" json:"tag_id"`
}

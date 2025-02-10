package models

type Tag struct {
	ID          int    `gorm:"comment:'分类唯一标识'" json:"id"`
	Name        string `gorm:"type:varchar(255);comment:'分类名称'" json:"name"`
	Description string `gorm:"type:text;comment:'分类描述'" json:"description"`
}

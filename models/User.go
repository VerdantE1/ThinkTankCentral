package models

import "time"

type User struct {
	ID        int       `gorm:"primaryKey;comment:'用户唯一标识'" json:"id"`
	Username  string    `gorm:"type:varchar(50);comment:'用户名'" json:"username"`
	Password  string    `gorm:"type:varchar(50);comment:'建议加密存储'" json:"password"`
	Email     string    `gorm:"type:varchar(100);comment:'用户邮箱'" json:"email"`
	Avatar    string    `gorm:"type:varchar(255);comment:'用户头像URL'" json:"avatar"`
	Bio       string    `gorm:"type:text;comment:'用户简介'" json:"bio"`
	Role      string    `gorm:"type:enum('admin','user');comment:'职位'" json:"role"`
	CreatedAt string    `gorm:"type:varchar(255);comment:'注册时间'" json:"created_at"`
	UpdatedAt time.Time `gorm:"comment:'最后更新时间'" json:"updated_at"`
}

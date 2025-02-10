package models

type Setting struct {
	ID    int    `gorm:"primaryKey;comment:'设置唯一标识'" json:"id"`
	Key   string `gorm:"type:varchar(255);comment:'配置项名称'" json:"key"`
	Value string `gorm:"type:text;comment:'配置项值'" json:"value"`
}

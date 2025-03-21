package other

import (
	"ThinkTankCentral/model/request"
	"gorm.io/gorm"
)

type MySQLOption struct {
	request.PageInfo
	Order   string
	Where   *gorm.DB
	Preload []string
}

package config

import (
	"ThinkTankCentral/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDatabase() {
	//数据源名称dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		GetEnv("DB_USER"),
		GetEnv("DB_PASSWORD"),
		GetEnv("DB_HOST"),
		GetEnv("DB_PORT"),
		GetEnv("DB_NAME"),
	)

	//连接数据库
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// 自动迁移
	err = DB.AutoMigrate(
		&models.Comment{},
		&models.PostTag{},
		&models.Post{},
		&models.Setting{},
		&models.Tag{},
		&models.User{},
	)
	if err != nil {
		panic("failed to migrate database")
	}
}

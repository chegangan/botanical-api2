package models

import (
	"fmt"
	"log"

	"botanical-api2/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // 导入 MySQL 驱动
)

var db *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error

	// 使用 setting 包中的数据库配置
	dbType := setting.Database.Type
	user := setting.Database.User
	password := setting.Database.Password
	host := setting.Database.Host
	name := setting.Database.Name

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user, password, host, name)

	db, err = gorm.Open(dbType, dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 配置数据库连接池
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// 启用日志
	if setting.RunMode == "debug" {
		db.LogMode(true)
	}

	// 迁移数据库模式
	db.AutoMigrate(&User{})
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return db
}

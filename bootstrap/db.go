package bootstrap

import (
	"GoBlog/app/models/article"
	"GoBlog/app/models/user"
	"GoBlog/pkg/model"
	"gorm.io/gorm"
	"time"
)

// SetupDB 初始化数据库和 ORM
func SetupDB() {
	db := model.ConnectDB()

	sqlDB, _ := db.DB()

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置最大空闲连接
	sqlDB.SetMaxIdleConns(25)
	// 设置每个连接的过期时间
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	// 创建和维护数据表结构
	migration(db)
}

// migration 创建和维护数据表结构
func migration(db *gorm.DB) {
	// 自动迁移
	db.AutoMigrate(
		&user.User{},
		&article.Article{},
	)
}

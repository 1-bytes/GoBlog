package database

import (
	"GoBlog/pkg/config"
	"GoBlog/pkg/logger"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB

// Initialize 初始化数据库
func Initialize() {
	initDB()
	createTables()
}

// initDB 初始化数据库连接.
func initDB() {
	var err error
	cfg := mysql.Config{
		Addr:                 config.GetString("database.mysql.host"),
		User:                 config.GetString("database.mysql.username"),
		Passwd:               config.GetString("database.mysql.password"),
		DBName:               config.GetString("database.mysql.database"),
		Net:                  "tcp",
		AllowNativePasswords: true,
	}

	// 设置数据库连接池
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	logger.LogError(err)

	// 设置最大连接数
	DB.SetMaxOpenConns(25)

	// 设置最大空闲连接数
	DB.SetMaxIdleConns(25)

	// 设置每个链接的过期时间
	DB.SetConnMaxLifetime(5 * time.Minute)

	err = DB.Ping()
	logger.LogError(err)
}

// createTables 创建数据库表.
func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := DB.Exec(createArticlesSQL)
	logger.LogError(err)
}

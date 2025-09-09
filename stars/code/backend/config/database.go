package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	// 从环境变量获取数据库配置
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "rootpassword")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "starpool")

	// 创建数据库连接字符串
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci", dbUser, dbPass, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	// 测试连接
	if err = DB.Ping(); err != nil {
		log.Fatal("数据库Ping失败: ", err)
	}

	fmt.Println("数据库连接成功!")
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

package config

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBの環境変数を読み込む
func loadDatabaseURL() (string, error) {
	mysqlDBName := os.Getenv("MYSQL_DATABASE")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlPort := os.Getenv("MYSQL_PORT")

	if mysqlDBName == "" || mysqlHost == "" || mysqlUser == "" || mysqlPassword == "" || mysqlPort == "" {
		return "", fmt.Errorf("環境変数が不足しています. MYSQL_DATABASE: %s, MYSQL_HOST: %s, MYSQL_USER: %s, MYSQL_PASSWORD: %s, MYSQL_PORT: %s", mysqlDBName, mysqlHost, mysqlUser, mysqlPassword, mysqlPort)
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDBName), nil
}

// データベース接続を行う
func NewDatabase() (*gorm.DB, error) {
	dbURL, err := loadDatabaseURL()
	if err != nil {
		return nil, fmt.Errorf("failed to load database URL: %w", err)
	}
	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	return db, nil
}

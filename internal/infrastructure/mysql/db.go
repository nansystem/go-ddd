package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// MySQLドライバを初期化のために必要
	_ "github.com/go-sql-driver/mysql"
)

// DBConfig はデータベース接続の設定を保持します
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

// NewConnection は新しいデータベース接続を作成します
func NewConnection(config DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FTokyo",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("データベース接続エラー: %w", err)
	}

	// 接続確認とコネクションプールの設定
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 接続テスト
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("データベース接続テストエラー: %w", err)
	}

	log.Println("データベースに接続しました")
	return db, nil
}

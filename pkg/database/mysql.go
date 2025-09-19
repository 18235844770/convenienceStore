package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"convenienceStore/pkg/config"
)

// NewMySQL 根据配置初始化 MySQL 连接，并应用连接池参数。
func NewMySQL(cfg config.DatabaseConfig, logger *log.Logger) (*sql.DB, error) {
	if cfg.Host == "" {
		return nil, errors.New("database host is not configured")
	}
	if cfg.User == "" {
		return nil, errors.New("database user is not configured")
	}
	if cfg.Name == "" {
		return nil, errors.New("database name is not configured")
	}

	port := cfg.Port
	if port == 0 {
		port = 3306
	}

	charset := cfg.Charset
	if charset == "" {
		charset = "utf8mb4"
	}

	parseTime := cfg.ParseTime
	if !cfg.ParseTime {
		parseTime = true
	}

	loc := cfg.Loc
	if loc == "" {
		loc = "Local"
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User,
		url.QueryEscape(cfg.Password),
		cfg.Host,
		port,
		cfg.Name,
		charset,
		parseTime,
		url.QueryEscape(loc),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime != "" {
		dur, err := time.ParseDuration(cfg.ConnMaxLifetime)
		if err != nil {
			return nil, fmt.Errorf("parse database conn_max_lifetime: %w", err)
		}
		db.SetConnMaxLifetime(dur)
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(timeoutCtx); err != nil {
		db.Close()
		return nil, err
	}

	if logger != nil {
		logger.Printf("mysql connected host=%s port=%d db=%s", cfg.Host, port, cfg.Name)
	}

	return db, nil
}

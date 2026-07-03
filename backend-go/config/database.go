package config

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Port       string
	JWTSecret  string
}

// DatabaseManager 管理数据库连接的单例管理器
type DatabaseManager struct {
	db   *sql.DB
	once sync.Once
}

var manager *DatabaseManager

func init() {
	manager = &DatabaseManager{}
}

// GetManager 获取数据库管理器实例
func GetManager() *DatabaseManager {
	return manager
}

// Init 初始化数据库连接
func (m *DatabaseManager) Init(cfg *Config) error {
	var err error
	m.once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		)

		m.db, err = sql.Open("mysql", dsn)
		if err != nil {
			err = fmt.Errorf("打开数据库连接失败：%w", err)
			return
		}

		m.db.SetMaxOpenConns(100)
		m.db.SetMaxIdleConns(10)
		m.db.SetConnMaxLifetime(time.Hour)
		m.db.SetConnMaxIdleTime(10 * time.Minute)

		if pingErr := m.db.Ping(); pingErr != nil {
			err = fmt.Errorf("数据库连接失败：%w", pingErr)
			return
		}

		fmt.Println("数据库连接成功")
	})
	return err
}

// GetDB 获取数据库连接
func (m *DatabaseManager) GetDB() *sql.DB {
	return m.db
}

// 向后兼容的函数，供现有代码使用
// 推荐使用 GetManager().GetDB() 替代

var DB *sql.DB

func LoadConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "root"),
		DBName:     getEnv("DB_NAME", "ad_bi_system"),
		Port:       getEnv("PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "ad-bi-jwt-secret-key"),
	}
}

func InitDatabase(cfg *Config) error {
	if err := GetManager().Init(cfg); err != nil {
		return err
	}
	// 更新全局变量以保持向后兼容
	DB = GetManager().GetDB()
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDB 获取数据库连接（向后兼容）
// 新代码应使用 config.GetManager().GetDB()
func GetDB() *sql.DB {
	return GetManager().GetDB()
}

// SetDB 设置数据库连接（仅供测试使用）
func SetDB(db *sql.DB) {
	manager.db = db
	DB = db
}
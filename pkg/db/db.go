package db

import (
	"fmt"
	"time"

	"github.com/thekrauss/beto-shared/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Driver   string // "postgres" / "mysql"
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string // Postgres
	LogLevel string // "silent", "error", "warn", "info"
}

func OpenDatabase(cfg Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	dsn := ""

	switch cfg.Driver {
	case "postgres":
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
		)
		dialector = postgres.Open(dsn)

	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
		)
		dialector = mysql.Open(dsn)

	default:
		return nil, errors.New(errors.CodeDBError, "unsupported database driver")
	}

	//configurable log level
	var gormLog logger.Interface
	switch cfg.LogLevel {
	case "info":
		gormLog = logger.Default.LogMode(logger.Info)
	case "warn":
		gormLog = logger.Default.LogMode(logger.Warn)
	case "error":
		gormLog = logger.Default.LogMode(logger.Error)
	default:
		gormLog = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLog,
	})
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBError, "failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBError, "failed to get sql.DB")
	}

	// Pooling settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func (c Config) ToURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

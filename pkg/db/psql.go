package db

import (
	"fmt"
	"github.com/lime008/gormzerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"helloworld/pkg/logger"
	"time"
)

type PsqlOptions struct {
	RunMode               string
	Host                  string
	Port                  int
	Username              string
	Password              string
	DBName                string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
}

func NewPsqlInstance(opts *PsqlOptions) (*gorm.DB, error) {
	//dsn := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=%t&loc=%s`,
	//	opts.Username,
	//	opts.Password,
	//	opts.Host,
	//	opts.Port,
	//	opts.DBName,
	//	true,
	//	"Local")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		opts.Host,
		opts.Username,
		opts.Password,
		opts.DBName,
		opts.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormzerolog.New(logger.GetLogger(), gormzerolog.Config{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: false,
		})})
	if opts.RunMode == "dev" {
		// 开发模式下输出 sql 语句的执行情况
		db.Use(&TracePlugin{})
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}

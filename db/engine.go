package db

import (
	"github.com/nova2018/easygin/file"
	"github.com/nova2018/easygin/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func initEngine(name string, cfg databaseConfig) {
	engine, err := newDb(cfg)
	if err != nil {
		panic(err)
	}
	setConnection(name, engine)
	SetPrefix(name, cfg.Prefix)
}

func newDb(cfg databaseConfig) (*gorm.DB, error) {
	// bootstrap
	var dialector gorm.Dialector
	switch cfg.Driver {
	case "mysql":
		dialector = mysql.Open(cfg.Connection)
	default:
		panic("driver [" + cfg.Driver + "] is not supported!")
	}

	lgGroup := newLoggerGroup()

	if cfg.SlowThreshold == 0 {
		cfg.SlowThreshold = time.Second
	}

	if cfg.ConsoleLog {
		lgGroup.AttachLogger(gormLogger.Default)
	}

	if cfg.SqlLog != "" {
		f, err := os.OpenFile(cfg.SqlLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logWriter := log.New(f, "", log.LstdFlags)
		lg := gormLogger.New(
			logWriter,
			gormLogger.Config{
				SlowThreshold:             cfg.SlowThreshold,
				Colorful:                  false,
				IgnoreRecordNotFoundError: false,
				LogLevel:                  gormLogger.LogLevel(cfg.LogLevel),
			},
		)
		lgGroup.AttachLogger(lg)
	}

	if cfg.LogConfig != nil {
		logWriter := log.New(file.NewWriter(cfg.LogConfig), "", log.LstdFlags)
		lg := gormLogger.New(
			logWriter,
			gormLogger.Config{
				SlowThreshold:             cfg.SlowThreshold,
				Colorful:                  false,
				IgnoreRecordNotFoundError: false,
				LogLevel:                  gormLogger.LogLevel(cfg.LogLevel),
			},
		)
		lgGroup.AttachLogger(lg)
	}

	{
		appLogger := newAppLogger(logger.Dynamic())
		appLogger.SlowThreshold = cfg.SlowThreshold
		lgGroup.AttachLogger(appLogger.LogMode(gormLogger.LogLevel(cfg.LogLevel)))
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.Prefix,
			SingularTable: true,
		},
		Logger: lgGroup,
	})

	if err != nil {
		panic(err)
	}

	// conn pool
	conn, err := db.DB()
	if err != nil {
		panic(err)
	}
	if cfg.Pool != nil {
		if cfg.Pool.ConnMaxLifetime > 0 {
			conn.SetConnMaxLifetime(cfg.Pool.ConnMaxLifetime)
		}
		if cfg.Pool.MaxIdleConns > 0 {
			conn.SetMaxIdleConns(cfg.Pool.MaxIdleConns)
		}
		if cfg.Pool.MaxOpenConns > 0 {
			conn.SetMaxOpenConns(cfg.Pool.MaxOpenConns)
		}
	}

	return db, nil
}

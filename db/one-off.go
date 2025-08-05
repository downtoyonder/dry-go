package db

import (
	"time"

	"github.com/downtoyonder/dry-go/config"
	"github.com/downtoyonder/dry-go/utils"
	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 一次性临时数据库
var OneOffDB = _OneOffDB{}

type _OneOffDB struct{}

func (o _OneOffDB) MySQL(dsn string) *gorm.DB {
	if dsn == "" {
		panic("dsn is empty")
	}

	c := config.NewViperFromMap(map[string]interface{}{
		"driver": "mysql",
		"dsn":    dsn,
		"debug":  true,
	})

	return NewDB(c, logger.Default)
}

func NewDB(c *viper.Viper, l logger.Interface) *gorm.DB {
	const (
		MYSQL    = "mysql"
		POSTGRES = "postgres"
		SQLITE   = "sqlite"
	)

	var (
		db  *gorm.DB
		err error

		dsn         = c.GetString("dsn")
		driver      = c.GetString("driver")
		enableDebug = c.GetBool("debug")

		gormCfg = &gorm.Config{
			Logger:                 l,
			PrepareStmt:            c.GetBool("gorm_prepare_stmt"),
			SkipDefaultTransaction: c.GetBool("gorm_skip_default_tx"),
		}
	)

	// GORM doc: https://gorm.io/docs/connecting_to_the_database.html
	switch driver {
	case MYSQL:
		db, err = gorm.Open(mysql.Open(dsn), gormCfg)
	case POSTGRES:
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), gormCfg)
	case SQLITE:
		db, err = gorm.Open(sqlite.Open(dsn), gormCfg)
	default:
		panic("unknown db driver")
	}

	utils.PanicErr(err)

	// Connection Pool config
	sqlDB, err := db.DB()
	utils.PanicErr(err)

	sqlDB.SetMaxIdleConns(30)
	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)

	if enableDebug {
		db = db.Debug()
	}

	return db
}

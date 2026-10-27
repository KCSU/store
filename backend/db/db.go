package db

import (
	"github.com/kcsu/store/config"
	"github.com/kcsu/store/logging"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect to the database
func Init(c *config.Config) (db *gorm.DB, err error) {
	dsn := c.DbConnection
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logging.Gorm{Level: logger.Warn},
	})
}

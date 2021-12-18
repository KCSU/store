package db

import (
	"github.com/kcsu/store/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(c *config.Config) (db *gorm.DB, err error) {
	dsn := c.DbConnection
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),
	})
}

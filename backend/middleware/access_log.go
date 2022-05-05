package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Access interface {
	Get(page int, size int) ([]model.AccessLog, error)
	Log(c echo.Context, verb string, metadata map[string]string) error
}

type DBAccess struct {
	db   *gorm.DB
	auth auth.Auth
}

// XXX: move to db package? Or its own?

func NewAccess(db *gorm.DB, auth auth.Auth) Access {
	return &DBAccess{db: db, auth: auth}
}

// Write an action to the access log
//
// Requires authentication middleware
func (d *DBAccess) Log(c echo.Context, verb string, metadata map[string]string) error {
	claims := d.auth.GetClaims(c)
	message := fmt.Sprint(claims.Name, " ", verb)
	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	log := &model.AccessLog{
		Email:    claims.Email,
		Message:  message,
		Metadata: metadataJson,
	}
	return d.db.Create(log).Error
}

// Get paginated access log
func (d *DBAccess) Get(page int, size int) ([]model.AccessLog, error) {
	offset := size * (page - 1)
	var logs []model.AccessLog
	// HACK: we get one more record to know if there are more records
	err := d.db.Limit(size + 1).Offset(offset).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

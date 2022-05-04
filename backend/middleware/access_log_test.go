package middleware_test

import (
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/kcsu/store/auth"
	"github.com/kcsu/store/middleware"
	mocks "github.com/kcsu/store/mocks/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestAccessLog(t *testing.T) {
	// Setup
	sdb, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer sdb.Close()
	pdb := postgres.New(postgres.Config{
		Conn: sdb,
	})
	db, err := gorm.Open(pdb, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	assert.NoError(t, err)
	a := mocks.NewAuth(t)
	access := middleware.NewAccess(db, a)
	// HTTP
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	a.On("GetClaims", c).Return(&auth.JwtClaims{
		Name:  "Chrisjen Avasarala",
		Email: "cj123@cam.ac.uk",
	})
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "access_logs"`).
		WithArgs(
			sqlmock.AnyArg(),
			"cj123@cam.ac.uk",
			`Chrisjen Avasarala created formal "Test"`,
			`{"name":"Test"}`,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))
	mock.ExpectCommit()
	// Run
	err = access.Log(c, `created formal "Test"`, map[string]string{
		"name": "Test",
	})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

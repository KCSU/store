package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/kcsu/store/db"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GroupSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	store GroupStore
}

func (s *GroupSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	pdb := postgres.New(postgres.Config{
		Conn: db,
	})
	s.db, err = gorm.Open(pdb, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(s.T(), err)
	s.store = NewGroupStore(s.db)
}

func (s *GroupSuite) TestGetGroups() {
	s.mock.ExpectQuery(`SELECT \* FROM "groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(56),
		)
	fs, err := s.store.Get()
	s.Require().NoError(err)
	s.Len(fs, 1)
	s.EqualValues(56, fs[0].ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestGroupSuite(t *testing.T) {
	suite.Run(t, new(GroupSuite))
}

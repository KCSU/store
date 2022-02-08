package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TODO: maybe use mocket for sql mocking?

type FormalSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	store FormalStore
}

func (s *FormalSuite) SetupSuite() {
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
	s.store = NewFormalStore(s.db)
}

func (s *FormalSuite) TestGetFormals() {
	s.mock.ExpectQuery(`SELECT \* FROM "formals" WHERE date_time > NOW()`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(3),
		)
	fs, err := s.store.Get()
	s.Require().NoError(err)
	s.Len(fs, 1)
	s.EqualValues(3, fs[0].ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestGetFormalsWithGroups() {
	s.mock.ExpectQuery(`SELECT \* FROM "formals" WHERE date_time > NOW()`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(3),
		)
	// Should also preload groups
	s.mock.ExpectQuery(`SELECT \* FROM "formal_groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"group_id", "formal_id"}).AddRow(42, 3),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(42),
		)
	fs, err := s.store.GetWithGroups()
	s.Require().NoError(err)
	s.Len(fs, 1)
	s.EqualValues(3, fs[0].ID)
	s.Len(fs[0].Groups, 1)
	s.EqualValues(42, fs[0].Groups[0].ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestAllFormals() {
	s.mock.ExpectQuery(`SELECT \* FROM "formals"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(56),
		)
	fs, err := s.store.All()
	s.Require().NoError(err)
	s.Len(fs, 1)
	s.EqualValues(56, fs[0].ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestFindFormal() {
	formal := model.Formal{
		Model: model.Model{
			ID: 4,
		},
		Name:   "Test",
		Groups: []model.Group{},
	}
	s.mock.ExpectQuery(`SELECT \* FROM "formals"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(formal.ID, formal.Name),
		)
	// Should also preload groups
	s.mock.ExpectQuery(`SELECT \* FROM "formal_groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"group_id"}),
		)
	f, err := s.store.Find(4)
	s.Require().NoError(err)
	s.Equal(formal, f)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestTicketsRemaining() {
	f := model.Formal{}
	f.ID = 42
	f.Tickets = 23
	f.GuestTickets = 27
	mockCount := 10
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).AddRow(mockCount),
		)
	tr := s.store.TicketsRemaining(&f, false)
	s.Equal(f.Tickets-uint(mockCount), tr)
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).AddRow(mockCount),
		)
	tr = s.store.TicketsRemaining(&f, true)
	s.Equal(f.GuestTickets-uint(mockCount), tr)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestFormalSuite(t *testing.T) {
	suite.Run(t, new(FormalSuite))
}

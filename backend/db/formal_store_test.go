package db_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
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

func (s *FormalSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)

	pdb := postgres.New(postgres.Config{
		Conn: db,
	})
	s.db, err = gorm.Open(pdb, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	s.Require().NoError(err)
	s.store = NewFormalStore(s.db)
}

func (s *FormalSuite) TearDownTest() {
	db, err := s.db.DB()
	s.Require().NoError(err)
	db.Close()
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

func (s *FormalSuite) TestGetGroups() {
	ids := []int{2, 4, 7}
	targetGroups := []model.Group{
		{Model: model.Model{ID: 2}, Name: "A"},
		{Model: model.Model{ID: 4}, Name: "B"},
		{Model: model.Model{ID: 7}, Name: "C"},
	}
	s.mock.ExpectQuery(`SELECT \* FROM "groups"`).WithArgs(2, 4, 7).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(2, "A").AddRow(4, "B").AddRow(7, "C"),
		)
	groups, err := s.store.GetGroups(ids)
	s.NoError(err)
	s.Equal(targetGroups, groups)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestCreateFormal() {
	formal := model.Formal{
		Name: "Test",
		Groups: []model.Group{
			{Model: model.Model{ID: 5}},
		},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "formals"`).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(12),
	)
	s.mock.ExpectExec(`INSERT INTO "formal_groups"`).WithArgs(12, 5).
		WillReturnResult(
			sqlmock.NewResult(12, 1),
		)
	s.mock.ExpectCommit()
	err := s.store.Create(&formal)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestUpdateFormal() {
	// We need to check all fields are updated, even zeroes
	f := model.Formal{
		Model: model.Model{
			ID: 34,
		},
		Name:         "Test",
		Menu:         "A Menu",
		Price:        13,
		GuestPrice:   0,
		GuestLimit:   0,
		Tickets:      150,
		GuestTickets: 50,
		SaleStart:    time.Date(2021, 5, 6, 10, 0, 0, 0, time.UTC),
		SaleEnd:      time.Date(2021, 5, 7, 11, 0, 0, 0, time.UTC),
		DateTime:     time.Date(2021, 6, 1, 17, 0, 0, 0, time.UTC),
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "formals"`).WithArgs(
		sqlmock.AnyArg(), nil,
		f.Name, f.Menu, f.Price,
		f.GuestPrice, f.GuestLimit, f.Tickets, f.GuestTickets,
		f.SaleStart, f.SaleEnd, f.DateTime, f.ID,
	).WillReturnResult(
		sqlmock.NewResult(int64(f.ID), 1),
	)
	s.mock.ExpectCommit()
	err := s.store.Update(&f)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestDeleteFormal() {
	f := model.Formal{
		Model: model.Model{
			ID: 34,
		},
		Name:         "Test",
		Menu:         "A Menu",
		Price:        13,
		GuestPrice:   0,
		GuestLimit:   0,
		Tickets:      150,
		GuestTickets: 50,
		SaleStart:    time.Date(2021, 5, 6, 10, 0, 0, 0, time.UTC),
		SaleEnd:      time.Date(2021, 5, 7, 11, 0, 0, 0, time.UTC),
		DateTime:     time.Date(2021, 6, 1, 17, 0, 0, 0, time.UTC),
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "formals" SET "deleted_at"`).WithArgs(
		sqlmock.AnyArg(), f.ID,
	).WillReturnResult(
		sqlmock.NewResult(int64(f.ID), 1),
	)
	s.mock.ExpectCommit()
	err := s.store.Delete(&f)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestUpdateFormalGroups() {
	groups := []model.Group{
		{
			Model: model.Model{ID: 3},
			Name:  "Group 1",
		},
		{
			Model: model.Model{ID: 56},
			Name:  "Group 2",
		},
	}
	formalId := 32
	formal := model.Formal{
		Model: model.Model{ID: uint(formalId)},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "formals"`).
		WithArgs(sqlmock.AnyArg(), formalId).
		WillReturnResult(
			sqlmock.NewResult(int64(formalId), 1),
		)
	s.mock.ExpectExec(`INSERT INTO "formal_groups"`).
		WithArgs(formalId, 3, formalId, 56).
		WillReturnResult(
			sqlmock.NewResult(56, 1),
		)
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "formal_groups"`).
		WithArgs(formalId, 3, 56).
		WillReturnResult(
			sqlmock.NewResult(0, 0),
		)
	s.mock.ExpectCommit()
	err := s.store.UpdateGroups(formal, groups)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestFormalSuite(t *testing.T) {
	suite.Run(t, new(FormalSuite))
}

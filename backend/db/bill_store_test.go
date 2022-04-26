package db_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	. "github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type BillSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	store BillStore
}

func (s *BillSuite) SetupTest() {
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
	s.store = NewBillStore(s.db)
}

func (s *BillSuite) TearDownTest() {
	db, err := s.db.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *BillSuite) TestGetBills() {
	id := uuid.New()
	bill := model.Bill{
		Name:  "Test Bill",
		Start: time.Now(),
		End:   time.Now().Add(24 * time.Hour),
	}
	bill.ID = id
	s.mock.ExpectQuery(`SELECT \* FROM "bills"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "start", "end"}).AddRow(
				id, bill.Name, bill.Start, bill.End,
			),
		)
	bs, err := s.store.Get()
	s.Require().NoError(err)
	s.Len(bs, 1)
	s.Equal(bill, bs[0])
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *BillSuite) TestFindBill() {
	id := uuid.New()
	bill := model.Bill{
		Model: model.Model{ID: id},
		Name:  "Test Bill",
		Start: time.Now().Add(-12 * time.Hour),
		End:   time.Now().Add(24 * time.Hour),
	}
	s.mock.ExpectQuery(`SELECT \* FROM "bills"`).
		WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "start", "end"}).
				AddRow(
					id, bill.Name, bill.Start, bill.End,
				),
		)
	b, err := s.store.Find(id)
	s.Require().NoError(err)
	s.Equal(bill, b)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *BillSuite) TestFindBillWithFormals() {
	id := uuid.New()
	bill := model.Bill{
		Model: model.Model{ID: id},
		Name:  "Test Bill",
		Start: time.Now().Add(-12 * time.Hour),
		End:   time.Now().Add(24 * time.Hour),
		Formals: []model.Formal{{
			Model:  model.Model{ID: uuid.New()},
			Name:   "Test Formal",
			BillID: &id,
		}},
	}
	s.mock.ExpectQuery(`SELECT \* FROM "bills"`).
		WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "start", "end"}).
				AddRow(
					id, bill.Name, bill.Start, bill.End,
				),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "formals"`).
		WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "bill_id"}).
				AddRow(
					bill.Formals[0].ID, bill.Formals[0].Name, bill.ID,
				),
		)
	b, err := s.store.FindWithFormals(id)
	s.Require().NoError(err)
	s.Equal(bill, b)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *BillSuite) TestUpdateBill() {
	id := uuid.New()
	bill := model.Bill{
		Model: model.Model{ID: id},
		Name:  "Test Bill",
		Start: time.Now().Add(-12 * time.Hour),
		End:   time.Now().Add(24 * time.Hour),
		Formals: []model.Formal{{
			Model:  model.Model{ID: uuid.New()},
			Name:   "Test Formal",
			BillID: &id,
		}},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "bills"`).
		WithArgs(
			sqlmock.AnyArg(), nil,
			bill.Name, bill.Start, bill.End, id,
		).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.Update(&bill)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *BillSuite) TestAddFormalToBill() {
	id := uuid.New()
	bill := model.Bill{
		Model:   model.Model{ID: id},
		Name:    "Test Bill",
		Start:   time.Now().Add(-12 * time.Hour),
		End:     time.Now().Add(24 * time.Hour),
		Formals: []model.Formal{},
	}
	formalId := uuid.New()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "formals"`).
		WithArgs(
			bill.ID, sqlmock.AnyArg(), formalId,
		).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.AddFormal(&bill, formalId)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *BillSuite) TestRemoveFormalFromBill() {
	id := uuid.New()
	bill := model.Bill{
		Model:   model.Model{ID: id},
		Name:    "Test Bill",
		Start:   time.Now().Add(-12 * time.Hour),
		End:     time.Now().Add(24 * time.Hour),
		Formals: []model.Formal{},
	}
	formalId := uuid.New()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "formals"`).
		WithArgs(
			nil, sqlmock.AnyArg(), formalId,
		).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.RemoveFormal(&bill, formalId)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestBillSuite(t *testing.T) {
	suite.Run(t, new(BillSuite))
}

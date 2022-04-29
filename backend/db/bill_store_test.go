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

func (s *BillSuite) TestCreateBill() {
	id := uuid.New()
	bill := model.Bill{
		Name:  "Test Bill",
		Start: time.Now(),
		End:   time.Now().Add(24 * time.Hour),
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "bills"`).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
			bill.Name, bill.Start, bill.End,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id),
		)
	s.mock.ExpectCommit()
	err := s.store.Create(&bill)
	s.NoError(err)
	s.EqualValues(id, bill.ID)
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

func (s *BillSuite) TestDeleteBill() {
	id := uuid.New()
	bill := model.Bill{}
	bill.ID = id
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "formals" SET "bill_id"`).
		WithArgs(nil, id).WillReturnResult(sqlmock.NewResult(0, 10))
	s.mock.ExpectExec(`UPDATE "bills" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.Delete(&bill)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *BillSuite) TestAddFormalsToBill() {
	id := uuid.New()
	bill := model.Bill{
		Model:   model.Model{ID: id},
		Name:    "Test Bill",
		Start:   time.Now().Add(-12 * time.Hour),
		End:     time.Now().Add(24 * time.Hour),
		Formals: []model.Formal{},
	}
	formalIds := []uuid.UUID{uuid.New(), uuid.New()}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "formals"`).
		WithArgs(
			bill.ID, sqlmock.AnyArg(), formalIds[0], formalIds[1],
		).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.AddFormals(&bill, formalIds)
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

func (s *BillSuite) TestGetCostBreakdown() {
	id := uuid.New()
	bill := model.Bill{
		Model: model.Model{ID: id},
		Name:  "Test Bill",
		Start: time.Now().Add(-12 * time.Hour),
		End:   time.Now().Add(24 * time.Hour),
	}
	fid1, fid2 := uuid.New(), uuid.New()
	dt := time.Now().Add(7 * 24 * time.Hour)
	s.mock.ExpectQuery(`SELECT .+ FROM "formals"`).WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"formal_id", "name", "price", "guest_price", "date_time",
				"standard", "guest",
			}).AddRow(
				fid1, "Test Formal", 10, 12,
				dt, 12, 11,
			).AddRow(
				fid2, "Test Formal 2", 23, 34.5,
				dt, 36, 81,
			),
		)
	breakdowns := []model.FormalCostBreakdown{
		{
			FormalID:   fid1,
			Name:       "Test Formal",
			Price:      10,
			GuestPrice: 12,
			DateTime:   dt,
			Standard:   12,
			Guest:      11,
		},
		{
			FormalID:   fid2,
			Name:       "Test Formal 2",
			Price:      23,
			GuestPrice: 34.5,
			DateTime:   dt,
			Standard:   36,
			Guest:      81,
		},
	}
	bs, err := s.store.GetCostBreakdown(&bill)
	s.NoError(err)
	s.Equal(breakdowns, bs)
}

func (s *BillSuite) TestGetCostBreakdownByUser() {
	id := uuid.New()
	bill := model.Bill{
		Model: model.Model{ID: id},
		Name:  "Test Bill",
		Start: time.Now().Add(-12 * time.Hour),
		End:   time.Now().Add(24 * time.Hour),
	}
	breakdowns := []model.UserCostBreakdown{
		{
			Email: "abc123@cam.ac.uk",
			Cost:  11.45,
		},
		{
			Email: "def456@cam.ac.uk",
			Cost:  12.3,
		},
	}
	s.mock.ExpectQuery(`SELECT .+ FROM "formals"`).WithArgs(id).
		WillReturnRows(
			sqlmock.NewRows([]string{
				"email", "cost",
			}).AddRow(
				"abc123@cam.ac.uk", 11.45,
			).AddRow(
				"def456@cam.ac.uk", 12.3,
			),
		)
	bs, err := s.store.GetCostBreakdownByUser(&bill)
	s.NoError(err)
	s.Equal(breakdowns, bs)
}

func TestBillSuite(t *testing.T) {
	suite.Run(t, new(BillSuite))
}

package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TicketSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	store TicketStore
}

func (s *TicketSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)
	// defer db.Close()

	pdb := postgres.New(postgres.Config{
		Conn: db,
	})
	s.db, err = gorm.Open(pdb, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	s.Require().NoError(err)
	s.store = NewTicketStore(s.db)
}

func (s *TicketSuite) TearDownTest() {
	db, err := s.db.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *TicketSuite) TestGetByUserId() {
	formal := model.Formal{
		Name:    "Test Formal",
		Tickets: 55,
	}
	formal.ID = 5
	ticket := model.Ticket{
		FormalID: int(formal.ID),
		IsGuest:  false,
		IsQueue:  true,
		Formal:   &formal,
	}
	ticket.ID = 1
	userId := 65
	s.mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WithArgs(userId).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "formal_id", "is_guest", "is_queue"}).
				AddRow(ticket.ID, ticket.FormalID, ticket.IsGuest, ticket.IsQueue),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "formals"`).
		WithArgs(formal.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "tickets"}).
				AddRow(formal.ID, formal.Name, formal.Tickets),
		)
	t, err := s.store.Get(userId)
	s.NoError(err)
	s.Equal([]model.Ticket{ticket}, t)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestFindTicket() {
	ticket := model.Ticket{
		FormalID: 21,
		IsGuest:  false,
		IsQueue:  true,
	}
	ticket.ID = 34
	s.mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WithArgs(ticket.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "formal_id", "is_guest", "is_queue"}).
				AddRow(ticket.ID, ticket.FormalID, ticket.IsGuest, ticket.IsQueue),
		)
	t, err := s.store.Find(int(ticket.ID))
	s.NoError(err)
	s.Equal(ticket, t)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestFindTicketWithFormal() {
	ticket := model.Ticket{
		FormalID: 21,
		IsGuest:  false,
		IsQueue:  true,
		Formal: &model.Formal{
			Model: model.Model{ID: 21},
			Name:  "Test Formal",
		},
	}
	ticket.ID = 34
	s.mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WithArgs(ticket.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "formal_id", "is_guest", "is_queue"}).
				AddRow(ticket.ID, ticket.FormalID, ticket.IsGuest, ticket.IsQueue),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "formals"`).
		WithArgs(ticket.FormalID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(ticket.Formal.ID, ticket.Formal.Name),
		)
	t, err := s.store.FindWithFormal(int(ticket.ID))
	s.NoError(err)
	s.Equal(ticket, t)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestUpdateTicket() {
	id := 312
	mealOption := "Vegetarian"
	// s.mock.Expe
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "tickets"`).
		WithArgs(mealOption, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))
	s.mock.ExpectCommit()
	err := s.store.Update(id, &dto.TicketRequestDto{MealOption: mealOption})
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())

}

func (s *TicketSuite) TestBatchCreate() {
	tickets := []model.Ticket{
		{
			FormalID: 1,
			IsGuest:  true,
		},
		{
			FormalID: 1,
			IsGuest:  false,
		},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "formal_id", "is_guest"}).
				AddRow(1, tickets[0].FormalID, tickets[0].IsGuest).
				AddRow(2, tickets[1].FormalID, tickets[1].IsGuest),
		)
	s.mock.ExpectCommit()
	err := s.store.BatchCreate(tickets)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestCreateTicket() {
	ticket := model.Ticket{
		FormalID: 76,
		IsGuest:  true,
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "formal_id", "is_guest"}).
				AddRow(1, ticket.FormalID, ticket.IsGuest),
		)
	s.mock.ExpectCommit()
	err := s.store.Create(&ticket)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestCountGuestByFormal() {
	userId := 420
	formalId := 12
	mockCount := 3
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WithArgs(formalId, userId).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).AddRow(mockCount),
		)
	count, err := s.store.CountGuestByFormal(formalId, userId)
	s.NoError(err)
	s.EqualValues(mockCount, count)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestExistsByFormal() {
	userId := 420
	formalId := 12
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WithArgs(formalId, userId).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).AddRow(1),
		)
	exists, err := s.store.ExistsByFormal(formalId, userId)
	s.NoError(err)
	s.True(exists)
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WithArgs(formalId, userId).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).AddRow(0),
		)
	exists, err = s.store.ExistsByFormal(formalId, userId)
	s.NoError(err)
	s.False(exists)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestDeleteByFormal() {
	userId := 234
	formalId := 567
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "tickets" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), formalId, userId).
		WillReturnResult(sqlmock.NewResult(32, 3))
	s.mock.ExpectCommit()
	err := s.store.DeleteByFormal(formalId, userId)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestDeleteTicket() {
	ticketId := 69
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "tickets" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), ticketId).
		WillReturnResult(sqlmock.NewResult(int64(ticketId), 1))
	s.mock.ExpectCommit()
	err := s.store.Delete(ticketId)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestTicketSuite(t *testing.T) {
	suite.Run(t, new(TicketSuite))
}

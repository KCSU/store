package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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
		Name:             "Test Formal",
		FirstSaleTickets: 55,
	}
	formal.ID = uuid.New()
	ticket := model.Ticket{
		FormalID: formal.ID,
		IsGuest:  false,
		IsQueue:  true,
		Formal:   &formal,
	}
	ticket.ID = uuid.New()
	userId := uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WithArgs(userId).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "formal_id", "is_guest", "is_queue"}).
				AddRow(ticket.ID, ticket.FormalID, ticket.IsGuest, ticket.IsQueue),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "formals"`).
		WithArgs(formal.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "first_sale_tickets"}).
				AddRow(formal.ID, formal.Name, formal.FirstSaleTickets),
		)
	// TODO: Test getting bill?
	t, err := s.store.Get(userId)
	s.NoError(err)
	s.Equal([]model.Ticket{ticket}, t)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestFindTicket() {
	fid := uuid.New()
	uid := uuid.New()
	ticket := model.Ticket{
		FormalID: fid,
		UserID:   uid,
		IsGuest:  false,
		IsQueue:  true,
		Formal: &model.Formal{
			Model: model.Model{ID: fid},
			Name:  "Test Formal",
		},
		User: &model.User{
			Model: model.Model{ID: uid},
			Name:  "Test User",
			Email: "tus123@cam.ac.uk",
		},
	}
	ticket.ID = uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WithArgs(ticket.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "formal_id", "user_id", "is_guest", "is_queue"}).
				AddRow(ticket.ID, ticket.FormalID, ticket.UserID, ticket.IsGuest, ticket.IsQueue),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "formals"`).
		WithArgs(ticket.FormalID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(ticket.Formal.ID, ticket.Formal.Name),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(ticket.UserID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(ticket.User.ID, ticket.User.Name, ticket.User.Email),
		)
	t, err := s.store.Find(ticket.ID)
	s.NoError(err)
	s.Equal(ticket, t)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestFindTicketWithFormal() {
	fid := uuid.New()
	ticket := model.Ticket{
		FormalID: fid,
		IsGuest:  false,
		IsQueue:  true,
		Formal: &model.Formal{
			Model: model.Model{ID: fid},
			Name:  "Test Formal",
		},
	}
	ticket.ID = uuid.New()
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
	t, err := s.store.FindWithFormal(ticket.ID)
	s.NoError(err)
	s.Equal(ticket, t)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestUpdateTicket() {
	id := uuid.New()
	mealOption := "Vegetarian"
	// s.mock.Expe
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "tickets"`).
		WithArgs(mealOption, sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.Update(id, &dto.TicketRequestDto{MealOption: mealOption})
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestBatchCreate() {
	fid := uuid.New()
	tickets := []model.Ticket{
		{
			FormalID: fid,
			IsGuest:  true,
		},
		{
			FormalID: fid,
			IsGuest:  false,
		},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()).
				AddRow(uuid.New()),
		)
	s.mock.ExpectCommit()
	err := s.store.BatchCreate(tickets)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestCreateTicket() {
	ticket := model.Ticket{
		FormalID: uuid.New(),
		IsGuest:  true,
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(uuid.New()),
		)
	s.mock.ExpectCommit()
	err := s.store.Create(&ticket)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *TicketSuite) TestCountGuestByFormal() {
	userId := uuid.New()
	formalId := uuid.New()
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
	userId := uuid.New()
	formalId := uuid.New()
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
	userId := uuid.New()
	formalId := uuid.New()
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
	ticketId := uuid.New()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "tickets" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), ticketId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.Delete(ticketId)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestTicketSuite(t *testing.T) {
	suite.Run(t, new(TicketSuite))
}

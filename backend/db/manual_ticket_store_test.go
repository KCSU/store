package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ManualTicketSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	store ManualTicketStore
}

func (s *ManualTicketSuite) SetupTest() {
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
	s.store = NewManualTicketStore(s.db)
}

func (s *ManualTicketSuite) TearDownTest() {
	db, err := s.db.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *ManualTicketSuite) TestFindManualTicket() {
	ticket := model.ManualTicket{
		FormalID:      13,
		MealOption:    "Vegan",
		Type:          "guest",
		Name:          "John Doe",
		Justification: "Dancer",
		Email:         "jd123@cam.ac.uk",
	}
	ticket.ID = 11
	s.mock.ExpectQuery(`SELECT \* FROM "manual_tickets"`).
		WithArgs(11).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"id", "formal_id", "meal_option", "type", "name", "justification", "email"},
			).AddRow(
				ticket.ID, ticket.FormalID, ticket.MealOption,
				ticket.Type, ticket.Name, ticket.Justification, ticket.Email,
			),
		)
	t, err := s.store.Find(11)
	s.Require().NoError(err)
	s.Equal(ticket, t)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *ManualTicketSuite) TestCreateManualTicket() {
	ticket := model.ManualTicket{
		FormalID:      13,
		MealOption:    "Vegan",
		Type:          "guest",
		Name:          "John Doe",
		Justification: "Dancer",
		Email:         "jd123@cam.ac.uk",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "manual_tickets"`).
		WithArgs(
			sqlmock.AnyArg(), sqlmock.AnyArg(), nil,
			ticket.MealOption, ticket.FormalID, ticket.Type,
			ticket.Name, ticket.Justification, ticket.Email,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(54),
		)
	s.mock.ExpectCommit()
	s.NoError(s.store.Create(&ticket))
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *ManualTicketSuite) TestDeleteManualTicket() {
	ticketId := 11
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "manual_tickets" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), ticketId).
		WillReturnResult(sqlmock.NewResult(int64(ticketId), 1))
	s.mock.ExpectCommit()
	s.NoError(s.store.Delete(ticketId))
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestManualTicketSuite(t *testing.T) {
	suite.Run(t, new(ManualTicketSuite))
}

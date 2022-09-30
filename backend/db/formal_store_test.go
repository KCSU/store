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
	id := uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "formals" WHERE date_time > NOW()`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id),
		)
	fs, err := s.store.Get()
	s.Require().NoError(err)
	s.Len(fs, 1)
	s.EqualValues(id, fs[0].ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestGetActiveFormals() {
	id := uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "formals" WHERE first_sale_start < NOW\(\) AND sale_end > NOW\(\)`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id),
		)
	fs, err := s.store.GetActive()
	s.Require().NoError(err)
	s.Len(fs, 1)
	s.EqualValues(id, fs[0].ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestGetFormalsWithUserData() {
	fid := uuid.New()
	gid := uuid.New()
	uid := uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "formals" WHERE date_time > NOW()`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(fid),
		)
	// Should also preload groups
	s.mock.ExpectQuery(`SELECT \* FROM "formal_groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"group_id", "formal_id"}).AddRow(gid, fid),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(gid),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WithArgs(fid, uid).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "formal_id", "meal_option"}).
				AddRow(uuid.New(), uid, fid, "Vegan"),
		)
	fs, err := s.store.GetWithUserData(uid)
	s.Require().NoError(err)
	s.Len(fs, 1)
	s.EqualValues(fid, fs[0].ID)
	s.Len(fs[0].Groups, 1)
	s.EqualValues(gid, fs[0].Groups[0].ID)
	s.Len(fs[0].TicketSales, 1)
	s.EqualValues(uid, fs[0].TicketSales[0].UserID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestAllFormals() {
	id := uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "formals"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id),
		)
	fs, err := s.store.All()
	s.Require().NoError(err)
	s.Len(fs, 1)
	s.EqualValues(id, fs[0].ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestFindFormal() {
	formal := model.Formal{
		Model: model.Model{
			ID: uuid.New(),
		},
		Name: "Test",
	}
	s.mock.ExpectQuery(`SELECT \* FROM "formals"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(formal.ID, formal.Name),
		)
	f, err := s.store.Find(formal.ID)
	s.Require().NoError(err)
	s.Equal(formal, f)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestFindFormalWithGroups() {
	formal := model.Formal{
		Model: model.Model{
			ID: uuid.New(),
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
	f, err := s.store.FindWithGroups(formal.ID)
	s.Require().NoError(err)
	s.Equal(formal, f)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestFindWithTickets() {
	formalID := uuid.New()
	userID := uuid.New()
	formal := model.Formal{
		Model: model.Model{
			ID: formalID,
		},
		Name:   "Test",
		Groups: []model.Group{},
		TicketSales: []model.Ticket{{
			Model:    model.Model{ID: uuid.New()},
			FormalID: formalID,
			UserID:   userID,
			User: &model.User{
				Model: model.Model{ID: userID},
				Name:  "James Holden",
				Email: "jh123@cam.ac.uk",
			},
			IsGuest: false,
			IsQueue: false,
		}},
		ManualTickets: []model.ManualTicket{{
			Model:         model.Model{ID: uuid.New()},
			FormalID:      formalID,
			MealOption:    "Vegan",
			Type:          "complimentary",
			Name:          "Bobby Draper",
			Justification: "Ents officer",
			Email:         "bd456@cam.ac.uk",
		}},
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
	mt := formal.ManualTickets[0]
	s.mock.ExpectQuery(`SELECT \* FROM "manual_tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "meal_option", "type", "name", "justification", "email", "formal_id"}).
				AddRow(mt.ID, mt.MealOption, mt.Type, mt.Name, mt.Justification, mt.Email, mt.FormalID),
		)
	t := formal.TicketSales[0]
	s.mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "user_id", "is_guest", "is_queue", "formal_id"}).
				AddRow(t.ID, t.UserID, t.IsGuest, t.IsQueue, t.FormalID),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email"}).
				AddRow(t.UserID, t.User.Name, t.User.Email),
		)
	f, err := s.store.FindWithTickets(formalID)
	s.Require().NoError(err)
	s.Equal(formal, f)
	s.Len(f.TicketSales, 1)
	s.EqualValues(t, f.TicketSales[0])
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestFindGuestList() {
	formalID := uuid.New()
	guests := []model.FormalGuest{
		{
			Name:   "James Holden",
			Email:  "jh123@cam.ac.uk",
			Guests: 2,
		},
		{
			Name:   "Bobby Draper",
			Email:  "bd456@cam.ac.uk",
			Guests: 1,
		},
	}
	rows := sqlmock.NewRows([]string{"name", "email", "guests"})
	for _, g := range guests {
		rows.AddRow(g.Name, g.Email, g.Guests)
	}
	s.mock.ExpectQuery(`SELECT .+ FROM "tickets"`).
		WithArgs(formalID).
		WillReturnRows(rows)
	f, err := s.store.FindGuestList(formalID)
	s.Require().NoError(err)
	s.Equal(guests, f)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestGetQueueLength() {
	formalID := uuid.New()
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WithArgs(formalID).
		WillReturnRows(
			sqlmock.NewRows([]string{"count"}).
				AddRow(2),
		)
	l, err := s.store.GetQueueLength(formalID)
	s.Require().NoError(err)
	s.EqualValues(2, l)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestTicketsRemaining() {
	f := model.Formal{}
	f.ID = uuid.New()
	f.FirstSaleTickets = 11
	f.SecondSaleTickets = 12
	f.FirstSaleGuestTickets = 24
	f.SecondSaleGuestTickets = 3
	mockCount := 10
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).AddRow(mockCount),
		)
	tr := s.store.TicketsRemaining(&f, false)
	s.Equal(f.FirstSaleTickets+f.SecondSaleTickets-uint(mockCount), tr)
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).AddRow(mockCount),
		)
	tr = s.store.TicketsRemaining(&f, true)
	s.Equal(f.FirstSaleGuestTickets+f.SecondSaleGuestTickets-uint(mockCount), tr)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestGetGroups() {
	ids := []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}
	targetGroups := []model.Group{
		{Model: model.Model{ID: ids[0]}, Name: "A"},
		{Model: model.Model{ID: ids[1]}, Name: "B"},
		{Model: model.Model{ID: ids[2]}, Name: "C"},
	}
	s.mock.ExpectQuery(`SELECT \* FROM "groups"`).WithArgs(ids[0], ids[1], ids[2]).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(ids[0], "A").AddRow(ids[1], "B").AddRow(ids[2], "C"),
		)
	groups, err := s.store.GetGroups(ids)
	s.NoError(err)
	s.Equal(targetGroups, groups)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestCreateFormal() {
	fid := uuid.New()
	gid := uuid.New()
	formal := model.Formal{
		Name: "Test",
		Groups: []model.Group{
			{Model: model.Model{ID: gid}},
		},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "formals"`).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(fid),
	)
	s.mock.ExpectQuery(`INSERT INTO "formal_groups"`).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"group_id", "formal_id"}).AddRow(gid, fid),
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
			ID: uuid.New(),
		},
		Name:                   "Test",
		Menu:                   "A Menu",
		Price:                  13,
		GuestPrice:             0,
		GuestLimit:             0,
		FirstSaleTickets:       150,
		SecondSaleTickets:      10,
		FirstSaleGuestTickets:  50,
		SecondSaleGuestTickets: 30,
		FirstSaleStart:         time.Date(2021, 5, 6, 10, 0, 0, 0, time.UTC),
		SecondSaleStart:        time.Date(2021, 5, 7, 10, 0, 0, 0, time.UTC),
		SaleEnd:                time.Date(2021, 5, 7, 11, 0, 0, 0, time.UTC),
		DateTime:               time.Date(2021, 6, 1, 17, 0, 0, 0, time.UTC),
		HasGuestList:           true,
		IsVisible:              true,
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "formals"`).WithArgs(
		sqlmock.AnyArg(), nil,
		f.Name, f.Menu, f.Price,
		f.GuestPrice, f.GuestLimit, f.FirstSaleTickets, f.FirstSaleGuestTickets, f.FirstSaleStart,
		f.SecondSaleTickets, f.SecondSaleGuestTickets, f.SecondSaleStart,
		f.SaleEnd, f.DateTime, f.HasGuestList, f.IsVisible, nil, f.ID,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)
	s.mock.ExpectCommit()
	err := s.store.Update(&f)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestDeleteFormal() {
	f := model.Formal{
		Model: model.Model{
			ID: uuid.New(),
		},
		Name:                  "Test",
		Menu:                  "A Menu",
		Price:                 13,
		GuestPrice:            0,
		GuestLimit:            0,
		FirstSaleTickets:      150,
		FirstSaleGuestTickets: 50,
		FirstSaleStart:        time.Date(2021, 5, 6, 10, 0, 0, 0, time.UTC),
		SaleEnd:               time.Date(2021, 5, 7, 11, 0, 0, 0, time.UTC),
		DateTime:              time.Date(2021, 6, 1, 17, 0, 0, 0, time.UTC),
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "tickets" SET "deleted_at"`).WithArgs(
		sqlmock.AnyArg(), f.ID,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)
	s.mock.ExpectExec(`UPDATE "manual_tickets" SET "deleted_at"`).WithArgs(
		sqlmock.AnyArg(), f.ID,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)
	s.mock.ExpectExec(`UPDATE "formals" SET "deleted_at"`).WithArgs(
		sqlmock.AnyArg(), f.ID,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)
	s.mock.ExpectCommit()
	err := s.store.Delete(&f)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *FormalSuite) TestUpdateFormalGroups() {
	groups := []model.Group{
		{
			Model: model.Model{ID: uuid.New()},
			Name:  "Group 1",
		},
		{
			Model: model.Model{ID: uuid.New()},
			Name:  "Group 2",
		},
	}
	formalId := uuid.New()
	formal := model.Formal{
		Model: model.Model{ID: formalId},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "formals"`).
		WithArgs(sqlmock.AnyArg(), formalId).
		WillReturnResult(
			sqlmock.NewResult(0, 1),
		)
	s.mock.ExpectQuery(`INSERT INTO "formal_groups"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows([]string{"group_id", "formal_id"}).
				AddRow(groups[0].ID, formalId).AddRow(groups[1].ID, formalId),
		)
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "formal_groups"`).
		WithArgs(formalId, groups[0].ID, groups[1].ID).
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

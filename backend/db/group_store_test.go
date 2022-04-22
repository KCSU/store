package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	. "github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
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

func (s *GroupSuite) SetupTest() {
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
	s.store = NewGroupStore(s.db)
}

func (s *GroupSuite) TearDownTest() {
	db, err := s.db.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *GroupSuite) TestGetGroups() {
	id := uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(id),
		)
	fs, err := s.store.Get()
	s.Require().NoError(err)
	s.Len(fs, 1)
	s.EqualValues(id, fs[0].ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *GroupSuite) TestFindGroup() {
	id := uuid.New()
	group := model.Group{
		Model: model.Model{ID: id},
		Name:  "Group",
		GroupUsers: []model.GroupUser{
			{
				GroupID:   id,
				UserEmail: "abc123@cam.ac.uk",
			},
		},
	}
	s.mock.ExpectQuery(`SELECT \* FROM "groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(group.ID, group.Name),
		)
	s.mock.ExpectQuery(`SELECT \* FROM "group_users"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"group_id", "user_email"}).
				AddRow(
					group.GroupUsers[0].GroupID,
					group.GroupUsers[0].UserEmail,
				),
		)
	g, err := s.store.Find(id)
	s.Require().NoError(err)
	s.Equal(group, g)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *GroupSuite) TestAddUser() {
	group := model.Group{
		Model: model.Model{ID: uuid.New()},
		Name:  "A Group",
	}
	email := "abc123@cam.ac.uk"
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "groups" SET "updated_at"`).
		WithArgs(sqlmock.AnyArg(), group.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectExec(`INSERT INTO "group_users"`).
		WithArgs(group.ID, email, true, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.AddUser(&group, email)
	s.Require().NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *GroupSuite) TestRemoveUser() {
	group := model.Group{
		Model: model.Model{ID: uuid.New()},
		Name:  "A Group",
	}
	email := "abc123@cam.ac.uk"
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "group_users"`).
		WithArgs(group.ID, email).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.RemoveUser(&group, email)
	s.Require().NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *GroupSuite) TestReplaceLookupUsers() {
	group := model.Group{
		Model:  model.Model{ID: uuid.New()},
		Name:   "Group 1",
		Type:   "inst",
		Lookup: "LKUP",
	}
	users := []model.GroupUser{
		{UserEmail: "abc123@cam.ac.uk", IsManual: false},
		{UserEmail: "def456@cam.ac.uk", IsManual: false},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "group_users"`).
		WithArgs(
			group.ID,
		).WillReturnResult(sqlmock.NewResult(0, 6))
	s.mock.ExpectExec(`INSERT INTO "group_users"`).WithArgs(
		group.ID, users[0].UserEmail, users[0].IsManual, sqlmock.AnyArg(),
		group.ID, users[1].UserEmail, users[1].IsManual, sqlmock.AnyArg(),
	).WillReturnResult(sqlmock.NewResult(0, 2))
	s.mock.ExpectCommit()
	err := s.store.ReplaceLookupUsers(&group, users)
	s.Require().NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *GroupSuite) TestCreateGroup() {
	group := model.Group{
		Name:   "New Group",
		Type:   "inst",
		Lookup: "LKUP",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "groups"`).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()),
	)
	s.mock.ExpectCommit()
	err := s.store.Create(&group)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *GroupSuite) TestUpdateGroup() {
	g := model.Group{
		Model: model.Model{
			ID: uuid.New(),
		},
		Name:   "Group",
		Type:   "inst",
		Lookup: "GRP01",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "groups"`).WithArgs(
		sqlmock.AnyArg(), nil,
		g.Name, g.Type, g.Lookup, g.ID,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)
	s.mock.ExpectCommit()
	err := s.store.Update(&g)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *GroupSuite) TestDeleteGroup() {
	g := model.Group{
		Model: model.Model{
			ID: uuid.New(),
		},
		Name: "To Delete",
		Type: "manual",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "group_users"`).WithArgs(g.ID).
		WillReturnResult(
			sqlmock.NewResult(0, 15),
		)
	s.mock.ExpectExec(`UPDATE "groups" SET "deleted_at"`).WithArgs(sqlmock.AnyArg(), g.ID).
		WillReturnResult(
			sqlmock.NewResult(0, 1),
		)
	s.mock.ExpectCommit()
	err := s.store.Delete(&g)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestGroupSuite(t *testing.T) {
	suite.Run(t, new(GroupSuite))
}

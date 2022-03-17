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

func (s *GroupSuite) TestFindGroup() {
	group := model.Group{
		Model: model.Model{ID: 5},
		Name:  "Group",
		GroupUsers: []model.GroupUser{
			{
				GroupID:   5,
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
	g, err := s.store.Find(5)
	s.Require().NoError(err)
	s.Equal(group, g)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *GroupSuite) TestAddUser() {
	group := model.Group{
		Model: model.Model{ID: 2},
		Name:  "A Group",
	}
	email := "abc123@cam.ac.uk"
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "groups" SET "updated_at"`).
		WithArgs(sqlmock.AnyArg(), group.ID).
		WillReturnResult(sqlmock.NewResult(int64(group.ID), 1))
	s.mock.ExpectExec(`INSERT INTO "group_users"`).
		WithArgs(group.ID, email, true, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(23, 1))
	s.mock.ExpectCommit()
	err := s.store.AddUser(&group, email)
	s.Require().NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *GroupSuite) TestRemoveUser() {
	group := model.Group{
		Model: model.Model{ID: 2},
		Name:  "A Group",
	}
	email := "abc123@cam.ac.uk"
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "group_users"`).
		WithArgs(group.ID, email).
		WillReturnResult(sqlmock.NewResult(2, 1))
	s.mock.ExpectCommit()
	err := s.store.RemoveUser(&group, email)
	s.Require().NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestGroupSuite(t *testing.T) {
	suite.Run(t, new(GroupSuite))
}

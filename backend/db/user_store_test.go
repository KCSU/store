package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/kcsu/store/db"
	"github.com/kcsu/store/model"
	"github.com/markbates/goth"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	store UserStore
}

func (s *UserSuite) SetupTest() {
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
	s.store = NewUserStore(s.db)
}

func (s *UserSuite) TearDownTest() {
	db, err := s.db.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *UserSuite) TestFindOrCreate() {
	au := goth.User{
		UserID: "abc123",
		Name:   "Chrisjen Avasarala",
		Email:  "cja67@cam.ac.uk",
	}
	userId := 1
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "users"`).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "name", "email", "provider_user_id"},
		).AddRow(
			userId, au.Name, au.Email, au.UserID,
		))
	s.mock.ExpectCommit()
	u, err := s.store.FindOrCreate(&au)
	user := model.User{
		Name:           au.Name,
		Email:          au.Email,
		ProviderUserId: au.UserID,
	}
	user.ID = uint(userId)
	s.NoError(err)
	s.Equal(user.Name, u.Name)
	s.Equal(user.ID, u.ID)
	s.Equal(user.Email, u.Email)
	s.Equal(user.ProviderUserId, u.ProviderUserId)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *UserSuite) TestFindUser() {
	user := model.User{
		Name:           "James Holden",
		Email:          "jmh23@cam.ac.uk",
		ProviderUserId: "abc123",
	}
	user.ID = 27
	s.mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "provider_user_id"}).
				AddRow(user.ID, user.Name, user.Email, user.ProviderUserId),
		)
	u, err := s.store.Find(int(user.ID))
	s.Require().NoError(err)
	s.Equal(user, u)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *UserSuite) TestExists() {
	email := "te51@cam.ac.uk"
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "users"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).AddRow(1),
		)
	exists, err := s.store.Exists(email)
	s.NoError(err)
	s.True(exists)
	s.mock.ExpectQuery(`SELECT count\(\*\) FROM "users"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"count(*)"}).AddRow(0),
		)
	exists, err = s.store.Exists(email)
	s.NoError(err)
	s.False(exists)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *UserSuite) TestGroups() {
	groups := make([]model.Group, 2)
	groups[0].ID = 1
	groups[0].Name = "test 1"
	groups[1].ID = 2
	groups[1].Name = "test 2"
	s.mock.ExpectQuery(`SELECT .* FROM "groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(groups[0].ID, groups[0].Name).
				AddRow(groups[1].ID, groups[1].Name),
		)
	u := model.User{}
	u.ID = 1
	g, err := s.store.Groups(&u)
	s.NoError(err)
	s.Equal(groups, g)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

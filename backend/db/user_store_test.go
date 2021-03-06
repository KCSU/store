package db_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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
	userId := uuid.New()
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "users"`).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id"},
		).AddRow(userId))
	s.mock.ExpectCommit()
	u, err := s.store.FindOrCreate(&au)
	user := model.User{
		Name:           au.Name,
		Email:          au.Email,
		ProviderUserId: au.UserID,
	}
	user.ID = userId
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
	user.ID = uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "provider_user_id"}).
				AddRow(user.ID, user.Name, user.Email, user.ProviderUserId),
		)
	u, err := s.store.Find(user.ID)
	s.Require().NoError(err)
	s.Equal(user, u)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *UserSuite) TestFindUserByEmail() {
	user := model.User{
		Name:           "James Holden",
		Email:          "jmh23@cam.ac.uk",
		ProviderUserId: "abc123",
	}
	user.ID = uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(user.Email).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "email", "provider_user_id"}).
				AddRow(user.ID, user.Name, user.Email, user.ProviderUserId),
		)
	u, err := s.store.FindByEmail(user.Email)
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
	groups[0].ID = uuid.New()
	groups[0].Name = "test 1"
	groups[1].ID = uuid.New()
	groups[1].Name = "test 2"
	s.mock.ExpectQuery(`SELECT .* FROM "groups"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(groups[0].ID, groups[0].Name).
				AddRow(groups[1].ID, groups[1].Name),
		)
	u := model.User{}
	u.ID = uuid.New()
	g, err := s.store.Groups(&u)
	s.NoError(err)
	s.Equal(groups, g)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *UserSuite) TestPermissions() {
	perms := []model.Permission{
		{
			ID:       uuid.New(),
			Resource: "formals",
			Action:   "read",
		},
		{
			ID:       uuid.New(),
			Resource: "tickets",
			Action:   "*",
		},
	}
	s.mock.ExpectQuery(`SELECT .* FROM "permissions"`).
		WithArgs().
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "resource", "action"}).
				AddRow(perms[0].ID, perms[0].Resource, perms[0].Action).
				AddRow(perms[1].ID, perms[1].Resource, perms[1].Action),
		)
	u := model.User{}
	u.ID = uuid.New()
	p, err := s.store.Permissions(&u)
	s.NoError(err)
	s.Equal(perms, p)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

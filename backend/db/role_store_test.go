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

type RoleSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	store RoleStore
}

func (s *RoleSuite) SetupTest() {
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
	s.store = NewRoleStore(s.db)
}

func (s *RoleSuite) TearDownTest() {
	db, err := s.db.DB()
	s.Require().NoError(err)
	db.Close()
}

func (s *RoleSuite) TestGetRoles() {
	roles := []model.Role{
		{
			Model: model.Model{ID: 2},
			Name:  "Formal Admin",
			Permissions: []model.Permission{
				{
					ID:       6,
					RoleID:   2,
					Resource: "formals",
					Action:   "read",
				},
				{
					ID:       7,
					RoleID:   2,
					Resource: "formals",
					Action:   "write",
				},
			},
		},
		{
			Model:       model.Model{ID: 4},
			Name:        "Doer of nothing",
			Permissions: []model.Permission{},
		},
	}
	s.mock.ExpectQuery(`SELECT \* FROM "roles"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).
				AddRow(roles[0].ID, roles[0].Name).
				AddRow(roles[1].ID, roles[1].Name),
		)
	ps := roles[0].Permissions
	s.mock.ExpectQuery(`SELECT \* FROM "permissions"`).
		WithArgs(roles[0].ID, roles[1].ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "role_id", "resource", "action"}).
				AddRow(ps[0].ID, ps[0].RoleID, ps[0].Resource, ps[0].Action).
				AddRow(ps[1].ID, ps[1].RoleID, ps[1].Resource, ps[1].Action),
		)
	rs, err := s.store.Get()
	s.Require().NoError(err)
	s.Equal(roles, rs)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleSuite))
}

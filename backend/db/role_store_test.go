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

func (s *RoleSuite) TestGetUserRoles() {
	user := model.User{
		Model: model.Model{ID: 5},
		Email: "abc123@cam.ac.uk",
		Name:  "A. Bell",
	}
	role := model.Role{
		Model: model.Model{ID: 7},
		Name:  "Admin",
	}
	userRole := model.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
		User:   user,
		Role:   role,
	}
	s.mock.ExpectQuery(`SELECT .* FROM "user_roles"`).
		WillReturnRows(
			sqlmock.NewRows([]string{"user_id", "role_id", "User__id", "Role__id", "User__email", "Role__name", "User__name"}).
				AddRow(user.ID, role.ID, user.ID, role.ID, user.Email, role.Name, user.Name),
		)
	urs, err := s.store.GetUserRoles()
	s.Require().NoError(err)
	s.Equal(urs, []model.UserRole{userRole})
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestFindRole() {
	role := model.Role{
		Model: model.Model{ID: 11},
		Name:  "My Role",
	}
	s.mock.ExpectQuery(`SELECT \* FROM "roles"`).WithArgs(role.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(
				role.ID, role.Name,
			),
		)
	r, err := s.store.Find(int(role.ID))
	s.Require().NoError(err)
	s.Equal(role, r)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestCreatePermission() {
	permission := model.Permission{
		RoleID:   5,
		Resource: "formals",
		Action:   "*",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "permissions"`).
		WithArgs(
			sqlmock.AnyArg(), "formals", "*", 5,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(15),
		)
	s.mock.ExpectCommit()
	err := s.store.CreatePermission(&permission)
	s.Require().NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestDeletePermission() {
	permissionId := 420
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "permissions"`).
		WithArgs(permissionId).
		WillReturnResult(sqlmock.NewResult(int64(permissionId), 1))
	s.mock.ExpectCommit()
	err := s.store.DeletePermission(permissionId)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestCreateRole() {
	role := model.Role{
		Name: "Treasurer",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "roles"`).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(92),
	)
	s.mock.ExpectCommit()
	err := s.store.Create(&role)
	s.NoError(err)
	s.EqualValues(92, role.ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestDeleteRole() {
	role := model.Role{
		Name:  "Ents",
		Model: model.Model{ID: 2},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "user_roles"`).
		WithArgs(role.ID).
		WillReturnResult(
			sqlmock.NewResult(int64(role.ID), 13),
		)
	s.mock.ExpectExec(`UPDATE "roles" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), role.ID).
		WillReturnResult(
			sqlmock.NewResult(int64(role.ID), 1),
		)
	s.mock.ExpectCommit()
	err := s.store.Delete(&role)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestAddUserRole() {
	role := model.Role{
		Name:  "Admin",
		Model: model.Model{ID: 69},
	}
	user := model.User{
		Name:  "Tony Stark",
		Email: "ts123@cam.ac.uk",
		Model: model.Model{ID: 22},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "roles" SET "updated_at"`).WithArgs(
		sqlmock.AnyArg(), role.ID,
	).WillReturnResult(
		sqlmock.NewResult(int64(role.ID), 1),
	)
	s.mock.ExpectExec(`INSERT INTO "user_roles"`).WithArgs(
		role.ID, user.ID,
	).WillReturnResult(
		sqlmock.NewResult(int64(role.ID), 1),
	)
	s.mock.ExpectCommit()
	err := s.store.AddUserRole(&role, &user)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestRemoveUserRole() {
	role := model.Role{
		Name:  "Admin",
		Model: model.Model{ID: 69},
	}
	user := model.User{
		Name:  "Tony Stark",
		Email: "ts123@cam.ac.uk",
		Model: model.Model{ID: 22},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "user_roles"`).WithArgs(
		role.ID, user.ID,
	).WillReturnResult(
		sqlmock.NewResult(int64(role.ID), 1),
	)
	s.mock.ExpectCommit()
	err := s.store.RemoveUserRole(&role, &user)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleSuite))
}

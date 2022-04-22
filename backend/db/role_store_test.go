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
	id := uuid.New()
	roles := []model.Role{
		{
			Model: model.Model{ID: id},
			Name:  "Formal Admin",
			Permissions: []model.Permission{
				{
					ID:       uuid.New(),
					RoleID:   id,
					Resource: "formals",
					Action:   "read",
				},
				{
					ID:       uuid.New(),
					RoleID:   id,
					Resource: "formals",
					Action:   "write",
				},
			},
		},
		{
			Model:       model.Model{ID: uuid.New()},
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
		Model: model.Model{ID: uuid.New()},
		Email: "abc123@cam.ac.uk",
		Name:  "A. Bell",
	}
	role := model.Role{
		Model: model.Model{ID: uuid.New()},
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
		Model: model.Model{ID: uuid.New()},
		Name:  "My Role",
	}
	s.mock.ExpectQuery(`SELECT \* FROM "roles"`).WithArgs(role.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name"}).AddRow(
				role.ID, role.Name,
			),
		)
	r, err := s.store.Find(role.ID)
	s.Require().NoError(err)
	s.Equal(role, r)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestCreatePermission() {
	permission := model.Permission{
		RoleID:   uuid.New(),
		Resource: "formals",
		Action:   "*",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "permissions"`).
		WithArgs(
			sqlmock.AnyArg(), "formals", "*", permission.RoleID,
		).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()),
		)
	s.mock.ExpectCommit()
	err := s.store.CreatePermission(&permission)
	s.Require().NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestDeletePermission() {
	permissionId := uuid.New()
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "permissions"`).
		WithArgs(permissionId).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	err := s.store.DeletePermission(permissionId)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestCreateRole() {
	id := uuid.New()
	role := model.Role{
		Name: "Treasurer",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "roles"`).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(id),
	)
	s.mock.ExpectCommit()
	err := s.store.Create(&role)
	s.NoError(err)
	s.EqualValues(id, role.ID)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestUpdateRole() {
	r := model.Role{
		Model: model.Model{
			ID: uuid.New(),
		},
		Name: "Admin2",
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "roles"`).WithArgs(
		sqlmock.AnyArg(), nil, r.Name, r.ID,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)
	s.mock.ExpectCommit()
	err := s.store.Update(&r)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestDeleteRole() {
	role := model.Role{
		Name:  "Ents",
		Model: model.Model{ID: uuid.New()},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "user_roles"`).
		WithArgs(role.ID).
		WillReturnResult(
			sqlmock.NewResult(0, 13),
		)
	s.mock.ExpectExec(`UPDATE "roles" SET "deleted_at"`).
		WithArgs(sqlmock.AnyArg(), role.ID).
		WillReturnResult(
			sqlmock.NewResult(0, 1),
		)
	s.mock.ExpectCommit()
	err := s.store.Delete(&role)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestAddUserRole() {
	role := model.Role{
		Name:  "Admin",
		Model: model.Model{ID: uuid.New()},
	}
	user := model.User{
		Name:  "Tony Stark",
		Email: "ts123@cam.ac.uk",
		Model: model.Model{ID: uuid.New()},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`UPDATE "roles" SET "updated_at"`).WithArgs(
		sqlmock.AnyArg(), role.ID,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)
	s.mock.ExpectQuery(`INSERT INTO "user_roles"`).WithArgs(
		sqlmock.AnyArg(), sqlmock.AnyArg(),
	).WillReturnRows(
		sqlmock.NewRows([]string{"role_id", "user_id"}).AddRow(role.ID, user.ID),
	)
	s.mock.ExpectCommit()
	err := s.store.AddUserRole(&role, &user)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RoleSuite) TestRemoveUserRole() {
	role := model.Role{
		Name:  "Admin",
		Model: model.Model{ID: uuid.New()},
	}
	user := model.User{
		Name:  "Tony Stark",
		Email: "ts123@cam.ac.uk",
		Model: model.Model{ID: uuid.New()},
	}
	s.mock.ExpectBegin()
	s.mock.ExpectExec(`DELETE FROM "user_roles"`).WithArgs(
		role.ID, user.ID,
	).WillReturnResult(
		sqlmock.NewResult(0, 1),
	)
	s.mock.ExpectCommit()
	err := s.store.RemoveUserRole(&role, &user)
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestRoleSuite(t *testing.T) {
	suite.Run(t, new(RoleSuite))
}

package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/kcsu/store/model"
	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// TODO: DOCUMENT THIS

func (ah *AdminHandler) GetRoles(c echo.Context) error {
	roles, err := ah.Roles.Get()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &roles)
}

func (ah *AdminHandler) GetUserRoles(c echo.Context) error {
	userRoles, err := ah.Roles.GetUserRoles()
	if err != nil {
		return err
	}
	urDto := make([]dto.UserRoleDto, len(userRoles))
	for i, ur := range userRoles {
		urDto[i] = dto.UserRoleDto{
			UserID:    ur.UserID,
			UserName:  ur.User.Name,
			UserEmail: ur.User.Email,
			RoleID:    ur.RoleID,
			RoleName:  ur.Role.Name,
		}
	}
	return c.JSON(http.StatusOK, &urDto)
}

// Move permissions logic to new file?

func (ah *AdminHandler) CreatePermission(c echo.Context) error {
	p := new(dto.PermissionDto)
	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(p); err != nil {
		return err
	}

	role, err := ah.Roles.Find(p.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
	}
	permission := model.Permission{
		RoleID:   role.ID,
		Resource: p.Resource,
		Action:   p.Action,
	}
	if err := ah.Roles.CreatePermission(&permission); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (ah *AdminHandler) DeletePermission(c echo.Context) error {
	// Get the permission ID from query
	id := c.Param("id")
	permissionID, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}

	if err := ah.Roles.DeletePermission(permissionID); err != nil {
		// What if it doesn't exist?
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (ah *AdminHandler) CreateRole(c echo.Context) error {
	r := new(dto.RoleDto)
	if err := c.Bind(r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(r); err != nil {
		return err
	}

	role := model.Role{
		Name: r.Name,
	}
	if err := ah.Roles.Create(&role); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (ah *AdminHandler) DeleteRole(c echo.Context) error {
	id := c.Param("id")
	roleId, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}
	role, err := ah.Roles.Find(roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if err := ah.Roles.Delete(&role); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (ah *AdminHandler) AddUserRole(c echo.Context) error {
	ur := new(dto.AddUserRoleDto)
	if err := c.Bind(ur); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(ur); err != nil {
		return err
	}

	role, err := ah.Roles.Find(ur.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	user, err := ah.Users.FindByEmail(ur.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if err := ah.Roles.AddUserRole(&role, &user); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (ah *AdminHandler) RemoveUserRole(c echo.Context) error {
	ur := new(dto.RemoveUserRoleDto)
	if err := c.Bind(ur); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(ur); err != nil {
		return err
	}

	role, err := ah.Roles.Find(ur.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	user, err := ah.Users.FindByEmail(ur.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if err := ah.Roles.RemoveUserRole(&role, &user); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

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

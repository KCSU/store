package admin

import (
	"net/http"

	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
)

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

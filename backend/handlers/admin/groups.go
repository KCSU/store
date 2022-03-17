package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Fetch a list of all groups
func (ah *AdminHandler) GetGroups(c echo.Context) error {
	groups, err := ah.Groups.Get()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &groups)
}

// Fetch an individual group
func (ah *AdminHandler) GetGroup(c echo.Context) error {
	// Get the group ID from query
	id := c.Param("id")
	groupID, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}
	group, err := ah.Groups.Find(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	return c.JSON(http.StatusOK, &group)
}

// Add a "manual" user to a group
func (ah *AdminHandler) AddGroupUser(c echo.Context) error {
	// Get the group ID from query
	id := c.Param("id")
	groupID, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}
	dto := new(dto.GroupUserDto)
	if err := c.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(dto); err != nil {
		return err
	}
	group, err := ah.Groups.Find(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if err := ah.Groups.AddUser(&group, dto.Email); err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

// Remove a "manual" user from a group
func (ah *AdminHandler) RemoveGroupUser(c echo.Context) error {
	// Get the group ID from query
	id := c.Param("id")
	groupID, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}
	dto := new(dto.GroupUserDto)
	if err := c.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(dto); err != nil {
		return err
	}
	group, err := ah.Groups.Find(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if err := ah.Groups.RemoveUser(&group, dto.Email); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

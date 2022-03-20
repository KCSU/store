package admin

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/kcsu/store/lookup"
	"github.com/kcsu/store/model"
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

// Create a group
func (ah *AdminHandler) CreateGroup(c echo.Context) error {
	g := new(dto.AdminGroupDto)
	if err := c.Bind(g); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(g); err != nil {
		return err
	}
	group := model.Group{
		Name:   g.Name,
		Type:   g.Type,
		Lookup: g.Lookup,
	}
	if err := ah.Groups.Create(&group); err != nil {
		return err
	}
	// JSON?
	return c.NoContent(http.StatusCreated)
}

// Update a group
func (ah *AdminHandler) UpdateGroup(c echo.Context) error {
	// Get the group ID from query
	id := c.Param("id")
	groupID, err := strconv.Atoi(id)
	if err != nil {
		return echo.ErrNotFound
	}
	g := new(dto.AdminGroupDto)
	if err := c.Bind(g); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(g); err != nil {
		return err
	}
	group, err := ah.Groups.Find(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	group.Name = g.Name
	group.Type = g.Type
	group.Lookup = g.Lookup
	if err := ah.Groups.Update(&group); err != nil {
		return err
	}
	// TODO: JSON?
	return c.NoContent(http.StatusOK)
}

// Delete a group
func (ah *AdminHandler) DeleteGroup(c echo.Context) error {
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
	if err := ah.Groups.Delete(&group); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
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
	// TODO: ensure user does not already exist
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

// Sync a group with the lookup directory
func (ah *AdminHandler) LookupGroupUsers(c echo.Context) error {
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
	lookupUrl, err := url.Parse(ah.Config.LookupApiUrl)
	if err != nil {
		return err
	}
	if err := lookup.ProcessGroup(group, ah.Groups, lookupUrl); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

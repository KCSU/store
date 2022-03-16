package admin

import (
	"errors"
	"net/http"
	"strconv"

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

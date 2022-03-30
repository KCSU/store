package admin

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Fetch a list of all formals
func (ah *AdminHandler) GetFormals(c echo.Context) error {
	formals, err := ah.Formals.All()
	if err != nil {
		return err
	}
	// Create JSON response
	// TODO: DRY!!!!!
	// Do we actually need any of this or can we just return the basic info?
	formalData := make([]dto.FormalDto, len(formals))
	for i, f := range formals {
		formalData[i].Formal = f
		// FIXME: This is horribly inefficient!!
		formalData[i].TicketsRemaining = ah.Formals.TicketsRemaining(&f, false)
		formalData[i].GuestTicketsRemaining = ah.Formals.TicketsRemaining(&f, true)
		groups := make([]dto.GroupDto, len(f.Groups))
		for j, g := range f.Groups {
			groups[j] = dto.GroupDto{
				ID:   g.ID,
				Name: g.Name,
			}
		}
		formalData[i].Groups = groups
	}
	return c.JSON(http.StatusOK, &formalData)
}

// Fetch a specified formal
func (ah *AdminHandler) GetFormal(c echo.Context) error {
	// Get the formal ID from query
	id := c.Param("id")
	formalID, err := strconv.Atoi(id)
	if err != nil {
		// TODO: NewHTTPError?
		return echo.ErrNotFound
	}
	formal, err := ah.Formals.Find(formalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	formalDto := dto.FormalDto{
		Formal: formal,
	}
	formalDto.TicketsRemaining = ah.Formals.TicketsRemaining(&formal, false)
	formalDto.GuestTicketsRemaining = ah.Formals.TicketsRemaining(&formal, true)
	groups := make([]dto.GroupDto, len(formal.Groups))
	for j, g := range formal.Groups {
		groups[j] = dto.GroupDto{
			ID:   g.ID,
			Name: g.Name,
		}
	}
	formalDto.Groups = groups
	return c.JSON(http.StatusOK, &formalDto)
}

// Create a formal
func (ah *AdminHandler) CreateFormal(c echo.Context) error {
	f := new(dto.CreateFormalDto)
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(f); err != nil {
		return err
	}
	formal := f.Formal()
	groups, err := ah.Formals.GetGroups(f.Groups)
	if err != nil {
		return err
	}

	if len(groups) != len(f.Groups) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Selected groups do not exist.")
	}
	formal.Groups = groups
	if err := ah.Formals.Create(&formal); err != nil {
		return err
	}
	// FIXME: JSON response?
	return c.NoContent(http.StatusCreated)
}

// Update details for a formal
func (ah *AdminHandler) UpdateFormal(c echo.Context) error {
	// Get the formal ID from query
	id := c.Param("id")
	formalID, err := strconv.Atoi(id)
	if err != nil {
		// TODO: NewHTTPError?
		return echo.ErrNotFound
	}
	f := new(dto.UpdateFormalDto)
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(f); err != nil {
		return err
	}
	// FIXME: CHECK THE FORMAL EXISTS
	formal := f.Formal()
	formal.ID = uint(formalID)
	if err := ah.Formals.Update(&formal); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// Delete a formal
func (ah *AdminHandler) DeleteFormal(c echo.Context) error {
	// Get the formal ID from query
	id := c.Param("id")
	formalID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	formal, err := ah.Formals.Find(formalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	if err := ah.Formals.Delete(&formal); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

// Update the list of groups who can buy tickets for the formal
func (ah *AdminHandler) UpdateFormalGroups(c echo.Context) error {
	// Get the formal ID from query
	id := c.Param("id")
	formalID, err := strconv.Atoi(id)
	if err != nil {
		// TODO: NewHTTPError?
		return echo.ErrNotFound
	}
	ids := []int{}
	if err := c.Bind(&ids); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	formal, err := ah.Formals.Find(formalID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	// Validate ids?
	groups, err := ah.Formals.GetGroups(ids)
	if err != nil {
		return err
	}

	if len(groups) != len(ids) {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "Selected groups do not exist.")
	}

	if err := ah.Formals.UpdateGroups(formal, groups); err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

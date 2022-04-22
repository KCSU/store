package admin

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// Fetch a list of all bills
func (ah *AdminHandler) GetBills(c echo.Context) error {
	bills, err := ah.Bills.Get()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &bills)
}

// Fetch an individual bill
func (ah *AdminHandler) GetBill(c echo.Context) error {
	// Get the bill ID from query
	id := c.Param("id")
	billID, err := uuid.Parse(id)
	if err != nil {
		return echo.ErrNotFound
	}
	bill, err := ah.Bills.Find(billID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	return c.JSON(http.StatusOK, &bill)
}

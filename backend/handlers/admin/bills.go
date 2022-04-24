package admin

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kcsu/store/model/dto"
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
	bill, err := ah.Bills.FindWithFormals(billID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	return c.JSON(http.StatusOK, &bill)
}

// Update a bill
func (ah *AdminHandler) UpdateBill(c echo.Context) error {
	// Get the bill ID from query
	id := c.Param("id")
	billID, err := uuid.Parse(id)
	if err != nil {
		return echo.ErrNotFound
	}
	b := dto.BillDto{}
	if err := c.Bind(&b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(b); err != nil {
		return err
	}
	bill, err := ah.Bills.Find(billID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	bill.Name = b.Name
	bill.Start, err = time.Parse("2006-01-02", b.Start)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	bill.End, err = time.Parse("2006-01-02", b.End)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ah.Bills.Update(&bill); err != nil {
		return err
	}
	// JSON?
	return c.NoContent(http.StatusOK)
}

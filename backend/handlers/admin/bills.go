package admin

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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
	b := new(dto.BillDto)
	if err := c.Bind(b); err != nil {
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

// Add a formal to a bill
func (ah *AdminHandler) AddBillFormals(c echo.Context) error {
	// Get the bill ID from query
	id := c.Param("id")
	billID, err := uuid.Parse(id)
	if err != nil {
		return echo.ErrNotFound
	}
	f := new(dto.AddFormalToBillDto)
	if err := c.Bind(f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(f); err != nil {
		return err
	}
	bill, err := ah.Bills.Find(billID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.ErrNotFound
		}
		return err
	}
	// TODO: check formal exists?
	if err := ah.Bills.AddFormals(&bill, f.FormalIDs); err != nil {
		return err
	}
	// JSON?
	return c.NoContent(http.StatusOK)
}

func (ah *AdminHandler) RemoveBillFormal(c echo.Context) error {
	// Get the bill ID from query
	id := c.Param("id")
	billID, err := uuid.Parse(id)
	if err != nil {
		return echo.ErrNotFound
	}
	fid := c.Param("formalId")
	formalID, err := uuid.Parse(fid)
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
	if err := ah.Bills.RemoveFormal(&bill, formalID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (ah *AdminHandler) GetBillStats(c echo.Context) error {
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
	formalCosts, err := ah.Bills.GetCostBreakdown(&bill)
	if err != nil {
		return err
	}
	userCosts, err := ah.Bills.GetCostBreakdownByUser(&bill)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &dto.BillStatsDto{
		Formals: formalCosts,
		Users:   userCosts,
	})
}

func (ah *AdminHandler) GetBillFormalStatsCSV(c echo.Context) error {
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
	formalCosts, err := ah.Bills.GetCostBreakdown(&bill)
	if err != nil {
		return err
	}
	c.Response().Header().Set(
		echo.HeaderContentDisposition,
		fmt.Sprintf("attachment; filename=%q", "formal_costs.csv"),
	)
	c.Response().WriteHeader(http.StatusOK)
	writer := csv.NewWriter(c.Response())
	defer writer.Flush()
	err = writer.Write([]string{
		"Formal", "Date", "King's Tickets",
		"King's Price", "Guest Tickets", "Guest Price",
		"Total",
	})
	if err != nil {
		return err
	}
	standardSum := 0
	guestSum := 0
	var costSum float32 = 0
	for _, f := range formalCosts {
		total := float32(f.Standard+f.StandardManual)*f.Price +
			float32(f.Guest+f.GuestManual)*f.GuestPrice
		if err := writer.Write([]string{
			f.Name, f.DateTime.Format("Jan 2 2006"),
			strconv.Itoa(f.Standard + f.StandardManual),
			fmt.Sprintf("%.2f", f.Price),
			strconv.Itoa(f.Guest + f.GuestManual),
			fmt.Sprintf("%.2f", f.GuestPrice),
			fmt.Sprintf("%.2f", total),
		}); err != nil {
			return err
		}
		standardSum += f.Standard + f.StandardManual
		guestSum += f.Guest + f.GuestManual
		costSum += total
	}
	err = writer.Write([]string{
		"Total", "", strconv.Itoa(standardSum),
		"", strconv.Itoa(guestSum), "", fmt.Sprintf("%.2f", costSum),
	})
	if err != nil {
		return err
	}
	return nil
}

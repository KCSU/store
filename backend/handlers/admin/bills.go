package admin

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kcsu/store/model"
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

// Create a bill
func (ah *AdminHandler) CreateBill(c echo.Context) error {
	b := new(dto.BillDto)
	if err := c.Bind(b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(b); err != nil {
		return err
	}
	bill := model.Bill{
		Name: b.Name,
	}
	var err error
	bill.Start, err = time.Parse("2006-01-02", b.Start)
	if err != nil {
		return err
	}
	bill.End, err = time.Parse("2006-01-02", b.End)
	if err != nil {
		return err
	}
	if err := ah.Bills.Create(&bill); err != nil {
		return err
	}
	if err := ah.Access.Log(c,
		fmt.Sprintf("created bill %s", bill.Name),
		map[string]string{
			"billId": bill.ID.String(),
		},
	); err != nil {
		return err
	}
	// JSON?
	return c.NoContent(http.StatusCreated)
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
	if err := ah.Access.Log(c,
		fmt.Sprintf("updated bill %s", bill.Name),
		map[string]string{
			"billId": bill.ID.String(),
		},
	); err != nil {
		return err
	}
	// JSON?
	return c.NoContent(http.StatusOK)
}

// Delete a bill
func (ah *AdminHandler) DeleteBill(c echo.Context) error {
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
	if err := ah.Bills.Delete(&bill); err != nil {
		return err
	}
	if err := ah.Access.Log(c,
		fmt.Sprintf("deleted bill %s", bill.Name),
		map[string]string{
			"billId": bill.ID.String(),
		},
	); err != nil {
		return err
	}
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
	// HACK: maybe do this with formatting or JSON arrays?
	formalIdStrings := make([]string, len(f.FormalIDs))
	for i, id := range f.FormalIDs {
		formalIdStrings[i] = id.String()
	}
	if err := ah.Access.Log(c,
		fmt.Sprintf("added formal(s) to bill %s", bill.Name),
		map[string]string{
			"billId":    bill.ID.String(),
			"formalIds": strings.Join(formalIdStrings, ","),
		},
	); err != nil {
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
	// TODO: formal name?
	if err := ah.Access.Log(c,
		fmt.Sprintf("removed formal from bill %s", bill.Name),
		map[string]string{
			"billId":   bill.ID.String(),
			"formalId": fid,
		},
	); err != nil {
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
		total := float32(f.Standard)*f.Price +
			float32(f.Guest)*f.GuestPrice
		if err := writer.Write([]string{
			f.Name, f.DateTime.Format("Jan 2 2006"),
			strconv.Itoa(f.Standard),
			fmt.Sprintf("%.2f", f.Price),
			strconv.Itoa(f.Guest),
			fmt.Sprintf("%.2f", f.GuestPrice),
			fmt.Sprintf("%.2f", total),
		}); err != nil {
			return err
		}
		standardSum += f.Standard
		guestSum += f.Guest
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

func (ah *AdminHandler) GetBillUserStatsCSV(c echo.Context) error {
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
	userCosts, err := ah.Bills.GetCostBreakdownByUser(&bill)
	if err != nil {
		return err
	}
	c.Response().Header().Set(
		echo.HeaderContentDisposition,
		fmt.Sprintf("attachment; filename=%q", "user_costs.csv"),
	)
	c.Response().WriteHeader(http.StatusOK)
	writer := csv.NewWriter(c.Response())
	defer writer.Flush()
	err = writer.Write([]string{
		"CRSID", "Total",
	})
	if err != nil {
		return err
	}
	var costSum float32 = 0
	for _, u := range userCosts {
		crsid := strings.Split(u.Email, "@")[0]
		err = writer.Write([]string{
			crsid, fmt.Sprintf("%.2f", u.Cost),
		})
		if err != nil {
			return err
		}
		costSum += u.Cost
	}
	err = writer.Write([]string{
		"Total", fmt.Sprintf("%.2f", costSum),
	})
	if err != nil {
		return err
	}
	return nil
}

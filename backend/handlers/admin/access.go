package admin

import (
	"net/http"

	"github.com/kcsu/store/model/dto"
	"github.com/labstack/echo/v4"
)

func (ah *AdminHandler) GetAccess(c echo.Context) error {
	dto := new(dto.PaginationDto)
	if err := c.Bind(dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(dto); err != nil {
		return err
	}
	logs, err := ah.Access.Get(dto.Page, dto.Size)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, logs)
}

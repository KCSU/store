package handlers

import (
	"net/http"

	"github.com/kcsu/store/model"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetFormals(c echo.Context) error {
	var formals []model.Formal
	// TODO: handle errors
	h.db.Find(&formals)
	return c.JSON(http.StatusOK, &formals)
}

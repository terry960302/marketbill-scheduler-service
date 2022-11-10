package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func PingPong(c echo.Context) error {
	msg := "Public Flower Upload Scheduler Service is running..."
	return c.String(http.StatusOK, msg)
}

package handlers

import (
	"github.com/TuanLe53/Roadmap-Broadcast-Server/templates"
	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func (h HomeHandler) HandleHome(c echo.Context) error {
	return Render(c, templates.Home())
}

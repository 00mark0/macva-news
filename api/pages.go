package api

import (
	"github.com/00mark0/macva-news/components"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (server *Server) homePage(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.Index())
}

func (server *Server) counterPage(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.CounterLayout())
}

func (server *Server) widgetPage(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.WidgetLayout())
}

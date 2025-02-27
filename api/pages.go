package api

import (
	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/token"
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

// full page to be served
func (server *Server) adminDash(ctx echo.Context) error {
	payload := ctx.Get(authorizationPayloadKey).(*token.Payload)

	return Render(ctx, http.StatusOK, components.DashPage(payload))
}

// htmx content insert
func (server *Server) adminDashContent(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.AdminDashboard())
}

func (server *Server) adminCats(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.AdminCategories())
}

func (server *Server) adminArts(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.AdminArticles())
}

func (server *Server) adminUsers(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.AdminUsers())
}

func (server *Server) adminAds(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.AdminAds())
}

func (server *Server) loginPage(ctx echo.Context) error {
	var loginErr components.LoginErr

	return Render(ctx, http.StatusOK, components.Login(loginErr))
}

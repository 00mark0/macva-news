package api

import (
	"log"
	"net/http"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/token"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
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
	var req ListCatsReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 10

	categories, err := server.store.ListCategories(ctx.Request().Context(), nextLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get categories", err))
		return err
	}

	return Render(ctx, http.StatusOK, components.AdminCategories(int(nextLimit), categories))
}

func (server *Server) createCategoryForm(ctx echo.Context) error {
	var createCategoryErr components.CreateCategoryErr

	return Render(ctx, http.StatusOK, components.CreateCategoryForm(createCategoryErr))
}

func (server *Server) deleteCategoryModal(ctx echo.Context) error {
	categoryID := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid category ID format", err))
	}

	// Create a pgtype.UUID with the parsed UUID
	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	category, err := server.store.GetCategoryByID(ctx.Request().Context(), pgUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get category", err))
		log.Println("Category error:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.DeleteCategoryModal(category))
}

func (server *Server) updateCategoryForm(ctx echo.Context) error {
	categoryID := ctx.Param("id")

	var updateCategoryErr components.UpdateCategoryErr

	// Parse string UUID into proper UUID format
	parsedUUID, err := uuid.Parse(categoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid category ID format", err))
	}

	// Create a pgtype.UUID with the parsed UUID
	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	category, err := server.store.GetCategoryByID(ctx.Request().Context(), pgUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get category", err))
		log.Println("Category error:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.UpdateCategoryForm(category, updateCategoryErr))
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

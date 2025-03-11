package api

import (
	"log"
	"net/http"
	"regexp"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type ListCatsReq struct {
	Limit int32 `query:"limit"`
}

func (server *Server) listCats(ctx echo.Context) error {
	var req ListCatsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listCats:", err)
		return err
	}

	nextLimit := req.Limit + 10

	categories, err := server.store.ListCategories(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing categories in listCats:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.AdminCategoriesDisplay(int(nextLimit), categories))
}

type CreateCatReq struct {
	CategoryName string `form:"category_name" validate:"required,min=3,max=50,regex=^[A-Za-z ]+$"`
}

func (server *Server) createCategory(ctx echo.Context) error {
	var req CreateCatReq
	var createCatErr components.CreateCategoryErr

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createCategory:", err)
		return err
	}

	if req.CategoryName == "" {
		createCatErr = "Ime kategorije je obavezno"

		return Render(ctx, http.StatusOK, components.CreateCategoryForm(createCatErr))
	}

	if len(req.CategoryName) < 3 {
		createCatErr = "Ime kategorije mora imati najmanje 3 slova"

		return Render(ctx, http.StatusOK, components.CreateCategoryForm(createCatErr))
	}

	if len(req.CategoryName) > 50 {
		createCatErr = "Ime kategorije može imati najviše 50 slova"

		return Render(ctx, http.StatusOK, components.CreateCategoryForm(createCatErr))
	}

	// Validate regex (only letters and spaces)
	matched, err := regexp.MatchString(`^[A-Za-z ]+$`, req.CategoryName)
	if err != nil {
		// Log error but don't expose internal issues to users
		log.Println("Error matching regex in createCategory:", err)
		return err
	}

	if !matched {
		createCatErr = "Ime kategorije može sadržati samo slova i razmake"
		return Render(ctx, http.StatusOK, components.CreateCategoryForm(createCatErr))
	}

	_, err = server.store.CreateCategory(ctx.Request().Context(), req.CategoryName)
	if err != nil {
		log.Println("Error creating category in createCategory:", err)
		return err
	}

	ctx.Response().Header().Set("HX-Trigger", `{"categoriesUpdated": ""}`)
	return ctx.NoContent(http.StatusOK)
}

func (server *Server) deleteCategory(ctx echo.Context) error {
	categoryID := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedUUID, err := uuid.Parse(categoryID)
	if err != nil {
		log.Println("Invalid category ID format in deleteCategory:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	_, err = server.store.DeleteCategory(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error deleting category in deleteCategory:", err)
		return err
	}

	return server.listCats(ctx)
}

type UpdateCatReq struct {
	CategoryName string `form:"category_name" validate:"required,min=3,max=50,regex=^[A-Za-z ]+$"`
}

func (server *Server) updateCategory(ctx echo.Context) error {
	categoryID := ctx.Param("id")
	var req UpdateCatReq
	var updateCatErr components.UpdateCategoryErr

	// Parse string UUID into proper UUID format
	parsedUUID, err := uuid.Parse(categoryID)
	if err != nil {
		log.Println("Invalid category ID format in updateCategory:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	category, err := server.store.GetCategoryByID(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error getting category in updateCategory:", err)
		return err
	}

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in updateCategory:", err)
		return err
	}

	if req.CategoryName == "" {
		updateCatErr = "Ime kategorije je obavezno"

		return Render(ctx, http.StatusOK, components.UpdateCategoryForm(category, updateCatErr))
	}

	if len(req.CategoryName) < 3 {
		updateCatErr = "Ime kategorije mora imati najmanje 3 slova"

		return Render(ctx, http.StatusOK, components.UpdateCategoryForm(category, updateCatErr))
	}

	if len(req.CategoryName) > 50 {
		updateCatErr = "Ime kategorije može imati najviše 50 slova"

		return Render(ctx, http.StatusOK, components.UpdateCategoryForm(category, updateCatErr))
	}

	// Validate regex (only letters and spaces)
	matched, err := regexp.MatchString(`^[A-Za-z ]+$`, req.CategoryName)
	if err != nil {
		// Log error but don't expose internal issues to users
		log.Println("Error matching regex in updateCategory:", err)
		return err
	}

	if !matched {
		updateCatErr = "Ime kategorije može sadržati samo slova i razmake"
		return Render(ctx, http.StatusOK, components.UpdateCategoryForm(category, updateCatErr))
	}

	arg := db.UpdateCategoryParams{
		CategoryID:   pgUUID,
		CategoryName: req.CategoryName,
	}

	_, err = server.store.UpdateCategory(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error updating category in updateCategory:", err)
		return err
	}

	ctx.Response().Header().Set("HX-Trigger", `{"categoriesUpdated": ""}`)
	return ctx.NoContent(http.StatusOK)
}

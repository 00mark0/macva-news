package api

import (
	"log"
	"net/http"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	"github.com/go-playground/validator/v10"
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
	CategoryName string `form:"category_name" validate:"required,min=3,max=50,regex"`
}

func (server *Server) createCategory(ctx echo.Context) error {
	var req CreateCatReq
	var createCatErr components.CreateCategoryErr

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createCategory:", err)
		return err
	}

	// Run validation
	if err := ctx.Validate(req); err != nil {
		// Loop through validation errors and handle them
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "CategoryName":
				switch fieldErr.Tag() {
				case "required":
					createCatErr = "Ime kategorije je obavezno"
				case "min":
					createCatErr = "Ime kategorije mora imati najmanje 3 slova"
				case "max":
					createCatErr = "Ime kategorije može imati najviše 50 slova"
				case "regex":
					createCatErr = "Ime kategorije može sadržati samo slova i razmake"
				}
			}
		}

		// Render the form with the custom error message
		return Render(ctx, http.StatusOK, components.CreateCategoryForm(createCatErr))
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing categories in createCategory:", err)
		return err
	}

	for _, category := range categories {
		if category.CategoryName == req.CategoryName {
			createCatErr = "Kategorija sa ovim imenom već postoji"
			return Render(ctx, http.StatusOK, components.CreateCategoryForm(createCatErr))
		}
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
	CategoryName string `form:"category_name" validate:"required,min=3,max=50,regex"`
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

	// Run validation
	if err := ctx.Validate(req); err != nil {
		// Loop through validation errors and handle them
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "CategoryName":
				switch fieldErr.Tag() {
				case "required":
					updateCatErr = "Ime kategorije je obavezno"
				case "min":
					updateCatErr = "Ime kategorije mora imati najmanje 3 slova"
				case "max":
					updateCatErr = "Ime kategorije može imati najviše 50 slova"
				case "regex":
					updateCatErr = "Ime kategorije može sadržati samo slova i razmake"
				}
			}
		}

		// Render the form with the custom error message
		return Render(ctx, http.StatusOK, components.UpdateCategoryForm(category, updateCatErr))
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Error listing categories in updateCategory:", err)
		return err
	}

	for _, category := range categories {
		if category.CategoryName == req.CategoryName {
			updateCatErr = "Kategorija sa ovim imenom već postoji"
			return Render(ctx, http.StatusOK, components.UpdateCategoryForm(category, updateCatErr))
		}
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

func (server *Server) listRecentCategoryContent(ctx echo.Context) error {
	var req ListPublishedLimitReq
	categoryIDStr := ctx.Param("id")

	categoryIDBytes, err := uuid.Parse(categoryIDStr)
	if err != nil {
		log.Println("Invalid category ID format in listRecentCategoryContent:", err)
		return err
	}

	categoryID := pgtype.UUID{
		Bytes: categoryIDBytes,
		Valid: true,
	}

	category, err := server.store.GetCategoryByID(ctx.Request().Context(), categoryID)
	if err != nil {
		log.Println("Error getting category in listRecentCategoryContent:", err)
		return err
	}

	nextLimit := req.Limit + 18

	arg := db.ListContentByCategoryLimitParams{
		CategoryID: categoryID,
		Limit:      nextLimit,
	}

	data, err := server.store.ListContentByCategoryLimit(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error listing content in listRecentCategoryContent:", err)
		return err
	}

	// Convert DB content items to ContentData
	var categoryContent []components.ContentData
	for _, item := range data {
		categoryContent = append(categoryContent, components.ContentData{
			ContentID:    item.ContentID,
			UserID:       item.UserID,
			CategoryID:   item.CategoryID,
			CategoryName: category.CategoryName, // Use the category name from the category object
			Title:        item.Title,
			Thumbnail: func() pgtype.Text {
				if item.Thumbnail.Valid && item.Thumbnail.String != "" {
					return item.Thumbnail
				}

				return pgtype.Text{String: ThumbnailURL, Valid: true}
			}(),
			ContentDescription:  item.ContentDescription,
			CommentsEnabled:     item.CommentsEnabled,
			ViewCountEnabled:    item.ViewCountEnabled,
			LikeCountEnabled:    item.LikeCountEnabled,
			DislikeCountEnabled: item.DislikeCountEnabled,
			Status:              item.Status,
			ViewCount:           item.ViewCount,
			LikeCount:           item.LikeCount,
			DislikeCount:        item.DislikeCount,
			CommentCount:        item.CommentCount,
			CreatedAt:           item.CreatedAt,
			UpdatedAt:           item.UpdatedAt,
			PublishedAt:         item.PublishedAt,
			IsDeleted:           item.IsDeleted,
		})
	}

	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings in listRecentCategoryContent:", err)
		return err
	}

	title := "Najnovije"

	return Render(ctx, http.StatusOK, components.GridCards(categoryContent, globalSettings[0], int(nextLimit), "", "", title))
}

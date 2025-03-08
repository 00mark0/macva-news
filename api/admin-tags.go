package api

import (
	"net/http"

	//"github.com/google/uuid"
	//"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	//"github.com/00mark0/macva-news/db/services"
)

type CreateTagReq struct {
	TagName string `form:"tag_name" validate:"required"`
}

func (server *Server) createTag(ctx echo.Context) error {
	var req CreateTagReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	if err := ctx.Validate(req); err != nil {
		message := "Ime je potrebno pri pravljenju novog taga."
		// Set header to indicate error and where to show it
		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	_, err := server.store.CreateTag(ctx.Request().Context(), req.TagName)
	if err != nil {
		message := "Tag sa ovim imenom već postoji."
		// Set header to indicate error and where to show it
		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	// Success case - get tags and render
	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get tags for create article page", err))
		return err
	}

	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		message := "Sadržaj je prebačen. Dodajte tagove sa uredi stranice."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	contentIDString := contentIDCookie.Value

	// Parse string UUID into proper UUID format
	parsedContentID, err := uuid.Parse(contentIDString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid content ID format", err))
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedContentID,
		Valid: true,
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get content tags", err))
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTags(tags, contentTags))
}

type AddTagReq struct {
	TagID string `form:"tag_id" validate:"required"`
}

func (server *Server) addTagToContent(ctx echo.Context) error {
	var req AddTagReq

	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		message := "Sadržaj nije pronađen. Dodajte tagove sa uredi stranice."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	contentIDString := contentIDCookie.Value

	// Parse string UUID into proper UUID format
	parsedContentID, err := uuid.Parse(contentIDString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid content ID format", err))
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedContentID,
		Valid: true,
	}

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	if err := ctx.Validate(req); err != nil {
		message := "Izaberite postojeći tag."
		// Set header to indicate error and where to show it
		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	parsedTagID, err := uuid.Parse(req.TagID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid tag ID format", err))
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	tagID := pgtype.UUID{
		Bytes: parsedTagID,
		Valid: true,
	}

	arg := db.AddTagToContentParams{
		ContentID: contentID,
		TagID:     tagID,
	}

	err = server.store.AddTagToContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Greška prilikom dodavanja taga."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	// Success case - get tags and render
	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get tags for create article page", err))
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get tags for create article page", err))
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTags(tags, contentTags))

}

func (server *Server) removeTagFromContent(ctx echo.Context) error {
	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		message := "Sadržaj nije pronađen. Dodajte tagove sa uredi stranice."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	contentIDString := contentIDCookie.Value

	// Parse string UUID into proper UUID format
	parsedContentID, err := uuid.Parse(contentIDString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid content ID format", err))
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedContentID,
		Valid: true,
	}

	tagIDString := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedTagID, err := uuid.Parse(tagIDString)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid tag ID format", err))
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	tagID := pgtype.UUID{
		Bytes: parsedTagID,
		Valid: true,
	}

	arg := db.RemoveTagFromContentParams{
		ContentID: contentID,
		TagID:     tagID,
	}

	err = server.store.RemoveTagFromContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Greška prilikom uklanjanja taga."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	// Success case - get tags and render
	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get tags for create article page", err))
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get tags for create article page", err))
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTags(tags, contentTags))
}

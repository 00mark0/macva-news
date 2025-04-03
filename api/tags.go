package api

import (
	"log"
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
		log.Println("Error binding request in createTag:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		message := "Ime je obavezno pri pravljenju novog taga."
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
		log.Println("Error listing tags in createTag:", err)
		return err
	}

	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		var contentTags []db.Tag

		// Set header to indicate success and where to show updated tags
		ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
		return Render(ctx, http.StatusOK, components.AdminTags(tags, contentTags))
	}

	contentIDString := contentIDCookie.Value

	// Parse string UUID into proper UUID format
	parsedContentID, err := uuid.Parse(contentIDString)
	if err != nil {
		log.Println("Invalid content ID format in createTag:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedContentID,
		Valid: true,
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in createTag:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))
}

type AddTagReq struct {
	TagID string `form:"tag_id" validate:"required"`
}

func (server *Server) addTagToContent(ctx echo.Context) error {
	var req AddTagReq

	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		message := "Sadržaj nije pronađen. Da bi dodali tagove pritisnite sačuvaj ili objavi."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	contentIDString := contentIDCookie.Value

	// Parse string UUID into proper UUID format
	parsedContentID, err := uuid.Parse(contentIDString)
	if err != nil {
		log.Println("Invalid content ID format in addTagToContent:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedContentID,
		Valid: true,
	}

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in addTagToContent:", err)
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
		log.Println("Invalid tag ID format in addTagToContent:", err)
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
		log.Println("Error listing tags in addTagToContent:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in addTagToContent:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))

}

func (server *Server) addTagToContentUpdate(ctx echo.Context) error {
	var req AddTagReq

	contentIDStr := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedContentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		log.Println("Invalid content ID format in addTagToContentUpdate:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedContentID,
		Valid: true,
	}

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in addTagToContentUpdate:", err)
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
		log.Println("Invalid tag ID format in addTagToContentUpdate:", err)
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
		log.Println("Error listing tags in addTagToContentUpdate:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in addTagToContentUpdate:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))

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
		log.Println("Invalid content ID format in removeTagFromContent:", err)
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
		log.Println("Invalid tag ID format in removeTagFromContent:", err)
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
		log.Println("Error listing tags in removeTagFromContent:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in removeTagFromContent:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))
}

func (server *Server) removeTagFromContentUpdate(ctx echo.Context) error {
	contentIDStr := ctx.Param("content_id")

	// Parse string UUID into proper UUID format
	parsedContentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		log.Println("Invalid content ID format in removeTagFromContent:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedContentID,
		Valid: true,
	}

	tagIDString := ctx.Param("tag_id")

	// Parse string UUID into proper UUID format
	parsedTagID, err := uuid.Parse(tagIDString)
	if err != nil {
		log.Println("Invalid tag ID format in removeTagFromContent:", err)
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
		log.Println("Failed to get tags in removeTagFromContent:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Failed to get tags by content in removeTagFromContent:", err)
		return err
	}

	// Set header to indicate success and where to show updated tags
	ctx.Response().Header().Set("HX-Retarget", "#admin-tags")
	return Render(ctx, http.StatusOK, components.AdminTagsUpdate(tags, contentTags, contentID.String()))
}

func (server *Server) deleteTag(ctx echo.Context) error {
	tagIDStr := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedTagID, err := uuid.Parse(tagIDStr)
	if err != nil {
		log.Println("Invalid tag ID format in deleteTag:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	tagID := pgtype.UUID{
		Bytes: parsedTagID,
		Valid: true,
	}

	err = server.store.DeleteTag(ctx.Request().Context(), tagID)
	if err != nil {
		log.Println("Error deleting tag in deleteTag:", err)
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

type ListTagsReq struct {
	Limit int32 `query:"limit"`
}

func (server *Server) listTags(ctx echo.Context) error {
	var req ListTagsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listTags:", err)
		return err
	}

	nextLimit := req.Limit + 20

	tags, err := server.store.ListTags(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing tags in listTags:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.TagsList(int(nextLimit), tags))
}

type SearchTagsReq struct {
	SearchTerm string `query:"search_term" validate:"required"`
}

func (server *Server) listSearchTags(ctx echo.Context) error {
	var req SearchTagsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listSearchTags:", err)
		return err
	}

	if err := ctx.Validate(req); err != nil {
		log.Println("Error validating request in listSearchTags:", err)
		return err
	}

	arg := db.SearchTagsParams{
		Search: req.SearchTerm,
		Limit:  20,
	}

	tags, err := server.store.SearchTags(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error searching tags in listSearchTags:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.TagsList(20, tags))
}

func (server *Server) listTagsByContent(ctx echo.Context) error {
	contentIDStr := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedContentID, err := uuid.Parse(contentIDStr)
	if err != nil {
		log.Println("Invalid content ID format in listTagsByContent:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedContentID,
		Valid: true,
	}

	tags, err := server.store.GetTagsByContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error getting tags by content in listTagsByContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.TagsInArticleDetes(tags))
}

// ContentByTagsHandler modified to return HTML instead of JSON
func (server *Server) listContentByTagsUnderCategory(ctx echo.Context) error {
	categoryIDStr := ctx.Param("id")
	categoryIDBytes, err := uuid.Parse(categoryIDStr)
	if err != nil {
		log.Println("Invalid category ID format in listContentByTagsUnderCategory:", err)
		return err
	}
	categoryID := pgtype.UUID{
		Bytes: categoryIDBytes,
		Valid: true,
	}

	// Get category info for the header
	category, err := server.store.GetCategoryByID(ctx.Request().Context(), categoryID)
	if err != nil {
		log.Println("Error fetching category in listContentByTagsUnderCategory:", err)
		return err
	}

	// Step 1: Get unique tags for the category
	tags, err := server.store.GetUniqueTagsByCategoryID(ctx.Request().Context(), categoryID)
	if err != nil {
		log.Println("Error fetching unique tags for category:", categoryID, err)
		return err
	}

	const limit = 6 // Show fewer items initially per tag
	const offset = 0

	// Step 2: Loop through tags and fetch content
	var contentByTags components.ContentByTagsList

	for _, tag := range tags {
		content, err := server.store.ListContentByTag(ctx.Request().Context(), db.ListContentByTagParams{
			TagName: tag.TagName,
			Limit:   limit,
			Offset:  offset,
		})

		for i := range content {
			if content[i].Thumbnail.String == "" {
				content[i].Thumbnail = pgtype.Text{String: ThumbnailURL, Valid: true}
			}
		}

		if err != nil {
			log.Println("Error fetching content for tag:", tag, err)
			continue // Skip this tag but continue with others
		}

		// Only add tags that have content
		if len(content) > 0 {
			contentByTags = append(contentByTags, components.ContentByTag{
				TagName: tag.TagName,
				Content: content,
			})
		}
	}

	// Get global settings for display options
	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error fetching global settings:", err)
		// Use default settings if we can't fetch them
		globalSettings[0] = db.GlobalSetting{
			DisableComments: false,
			DisableLikes:    false,
			DisableViews:    false,
		}
	}

	// Render the templ component
	return Render(ctx, http.StatusOK, components.ContentByTagsSection(contentByTags, globalSettings[0], category.CategoryName))
}

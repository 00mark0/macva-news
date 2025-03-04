package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
)

func (server *Server) listPubContent(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimit(ctx.Request().Context(), nextLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get content", err))
		return err
	}

	var content []components.ListPublishedContentRes

	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.Format("2006-01-02 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.Format("2006-01-02 15:04"),
			PublishedAt:         v.PublishedAt.Time.Format("2006-01-02 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/published?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listPubContentOldest(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimitOldest(ctx.Request().Context(), nextLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get content", err))
		return err
	}

	var content []components.ListPublishedContentRes

	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.Format("2006-01-02 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.Format("2006-01-02 15:04"),
			PublishedAt:         v.PublishedAt.Time.Format("2006-01-02 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/published/oldest?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listPubContentTitle(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimitTitle(ctx.Request().Context(), nextLimit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get content", err))
		return err
	}

	var content []components.ListPublishedContentRes

	for _, v := range data {
		content = append(content, components.ListPublishedContentRes{
			ContentID:           v.ContentID.String(),
			UserID:              v.UserID.String(),
			CategoryID:          v.CategoryID.String(),
			Title:               v.Title,
			ContentDescription:  v.ContentDescription,
			CommentsEnabled:     v.CommentsEnabled,
			ViewCountEnabled:    v.ViewCountEnabled,
			LikeCountEnabled:    v.LikeCountEnabled,
			DislikeCountEnabled: v.DislikeCountEnabled,
			Status:              v.Status,
			ViewCount:           v.ViewCount,
			LikeCount:           v.LikeCount,
			DislikeCount:        v.DislikeCount,
			CommentCount:        v.CommentCount,
			CreatedAt:           v.CreatedAt.Time.Format("2006-01-02 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.Format("2006-01-02 15:04"),
			PublishedAt:         v.PublishedAt.Time.Format("2006-01-02 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/published/title?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) archivePubContent(ctx echo.Context) error {
	id := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid content ID format", err))
	}

	// Create a pgtype.UUID with the parsed UUID
	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	_, err = server.store.SoftDeleteContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to archive content", err))
		return err
	}

	return server.adminArts(ctx)
}

type UpdateContentReq struct {
	ContentID           string  `query:"content_id"`
	Title               *string `form:"title" validate:"required"`
	ContentDescription  *string `form:"content_description" validate:"required"`
	CategoryID          *string `form:"category_id"`
	CommentsEnabled     *bool   `form:"comments_enabled"`
	ViewCountEnabled    *bool   `form:"view_count_enabled"`
	LikeCountEnabled    *bool   `form:"like_count_enabled"`
	DislikeCountEnabled *bool   `form:"dislike_count_enabled"`
}

func (server *Server) updateContent(ctx echo.Context) error {
	contentIDString := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedContentID, err := uuid.Parse(contentIDString)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid content ID format", err))
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedContentID,
		Valid: true,
	}

	var req UpdateContentReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	content, err := server.store.GetContentDetails(ctx.Request().Context(), contentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get content", err))
		return err
	}

	var parsedCategoryID uuid.UUID
	var categoryID pgtype.UUID
	if req.CategoryID != nil {
		parsedCategoryID, err = uuid.Parse(*req.CategoryID)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, errorResponse("invalid category ID format", err))
		}
		categoryID = pgtype.UUID{
			Bytes: parsedCategoryID,
			Valid: true,
		}
	} else {
		categoryID = content.CategoryID
	}

	arg := db.UpdateContentParams{
		ContentID: contentID,
		Title: func() string {
			if req.Title != nil {
				return *req.Title
			}
			return content.Title
		}(),
		ContentDescription: func() string {
			if req.ContentDescription != nil {
				return *req.ContentDescription
			}
			return content.ContentDescription
		}(),
		CategoryID: func() pgtype.UUID {
			if req.CategoryID != nil {
				return categoryID
			}
			return content.CategoryID
		}(),
		CommentsEnabled: func() bool {
			if req.CommentsEnabled != nil {
				return *req.CommentsEnabled
			}
			return content.CommentsEnabled
		}(),
		ViewCountEnabled: func() bool {
			if req.ViewCountEnabled != nil {
				return *req.ViewCountEnabled
			}
			return content.ViewCountEnabled
		}(),
		LikeCountEnabled: func() bool {
			if req.LikeCountEnabled != nil {
				return *req.LikeCountEnabled
			}
			return content.LikeCountEnabled
		}(),
		DislikeCountEnabled: func() bool {
			if req.DislikeCountEnabled != nil {
				return *req.DislikeCountEnabled
			}
			return content.DislikeCountEnabled
		}(),
	}

	_, err = server.store.UpdateContent(ctx.Request().Context(), arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to update content", err))
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

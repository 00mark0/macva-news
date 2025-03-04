package api

import (
	"github.com/labstack/echo/v4"
	"net/http"

	"github.com/00mark0/macva-news/components"
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

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content))
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

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content))
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

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content))
}

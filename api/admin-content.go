package api

import (
	"log"
	"net/http"
	"time"

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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/published/title?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDraftContent(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDraftContent(ctx.Request().Context(), nextLimit)
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/draft?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDraftContentOldest(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDraftContentOldest(ctx.Request().Context(), nextLimit)
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/draft/oldest?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDraftContentTitle(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDraftContentTitle(ctx.Request().Context(), nextLimit)
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/draft/title?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDelContent(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDeletedContent(ctx.Request().Context(), nextLimit)
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/deleted?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDelContentOldest(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDeletedContentOldest(ctx.Request().Context(), nextLimit)
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/deleted/oldest?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listDelContentTitle(ctx echo.Context) error {
	var req ListPublishedLimitReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListDeletedContentTitle(ctx.Request().Context(), nextLimit)
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/deleted/title?limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

type SearchContentReq struct {
	SearchTerm string `query:"search_term" validate:"required"`
	Limit      int32  `query:"limit"`
}

func (server *Server) listSearchPubContent(ctx echo.Context) error {
	var req SearchContentReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	if err := ctx.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("search term is required", err))
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchContentParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	data, err := server.store.SearchContent(ctx.Request().Context(), arg)
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/published/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listSearchDraftContent(ctx echo.Context) error {
	var req SearchContentReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	if err := ctx.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("search term is required", err))
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchDraftContentParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	data, err := server.store.SearchDraftContent(ctx.Request().Context(), arg)
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/draft/search?search_term=" + req.SearchTerm + "&limit="

	return Render(ctx, http.StatusOK, components.PublishedContent(int(nextLimit), content, url))
}

func (server *Server) listSearchDelContent(ctx echo.Context) error {
	var req SearchContentReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	if err := ctx.Validate(req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("search term is required", err))
		return err
	}

	nextLimit := req.Limit + 20

	arg := db.SearchDelContentParams{
		Limit:      nextLimit,
		SearchTerm: req.SearchTerm,
	}

	data, err := server.store.SearchDelContent(ctx.Request().Context(), arg)
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
			CreatedAt:           v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			UpdatedAt:           v.UpdatedAt.Time.In(Loc).Format("02-01-06 15:04"),
			PublishedAt:         v.PublishedAt.Time.In(Loc).Format("02-01-06 15:04"),
			IsDeleted:           v.IsDeleted.Bool,
			Username:            v.Username,
			CategoryName:        v.CategoryName,
		})
	}

	url := "/api/admin/content/deleted/search?search_term=" + req.SearchTerm + "&limit="

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

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get content overview", err))
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleNav(overview))
}

func (server *Server) deleteContent(ctx echo.Context) error {
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

	_, err = server.store.HardDeleteContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to delete content", err))
		return err
	}

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get content overview", err))
		return err
	}

	// delete media associated with content
	media, err := server.store.ListMediaForContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error listing media while deleting content:", err)
		return err
	}

	if len(media) > 0 {
		for _, v := range media {
			if err := server.deleteMediaFunc(v.MediaID.String()); err != nil {
				log.Println("Error deleting media deleting media while deleting content:", err)
				return err
			}
		}
	}

	return Render(ctx, http.StatusOK, components.ArticleNav(overview))
}

func (server *Server) publishDraftContent(ctx echo.Context) error {
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

	_, err = server.store.PublishContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to publish content", err))
		return err
	}

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get content overview", err))
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleNav(overview))
}

func (server *Server) unarchiveContent(ctx echo.Context) error {
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

	_, err = server.store.UnarchiveContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to unarchive content", err))
		return err
	}

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse("failed to get content overview", err))
		return err
	}

	return Render(ctx, http.StatusOK, components.ArticleNav(overview))
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

type CreateContentReq struct {
	CategoryID         string `form:"category_id"`
	Title              string `form:"title"`
	ContentDescription string `form:"content_description"`
}

func (server *Server) createContent(ctx echo.Context) error {
	var req CreateContentReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	if req.Title == "" {
		message := "Naslov je obavezan."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.CategoryID == "" {
		message := "Kategorija je obavezna."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.ContentDescription == "" {
		message := "Sadržaj je obavezan."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	cookie, err := ctx.Cookie("access_token")
	if err != nil {
		// No cookie found; redirect to login page
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	accessToken := cookie.Value
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		// Invalid token; redirect to login page
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	parsedUserID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid content ID format", err))
	}

	// Create a pgtype.UUID with the parsed UUID
	userID := pgtype.UUID{
		Bytes: parsedUserID,
		Valid: true,
	}

	parsedCategoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid content ID format", err))
	}

	// Create a pgtype.UUID with the parsed UUID
	categoryID := pgtype.UUID{
		Bytes: parsedCategoryID,
		Valid: true,
	}

	arg := db.CreateContentParams{
		UserID:              userID,
		CategoryID:          categoryID,
		Title:               req.Title,
		ContentDescription:  req.ContentDescription,
		CommentsEnabled:     true,
		ViewCountEnabled:    true,
		LikeCountEnabled:    true,
		DislikeCountEnabled: false,
	}

	content, err := server.store.CreateContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Failed to create content"

		return Render(ctx, http.StatusInternalServerError, components.ArticleError(message))
	}

	ctx.SetCookie(&http.Cookie{
		Name:    "content_id",
		Value:   content.ContentID.String(),
		MaxAge:  60 * 60 * 24 * 365,
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
	})

	message := "Uspešno ste sačuvali novi sadržaj."

	return Render(ctx, http.StatusOK, components.ArticleSuccess(message))
}

func (server *Server) createAndPublishContent(ctx echo.Context) error {
	var req CreateContentReq

	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse("invalid request body", err))
		return err
	}

	if req.Title == "" {
		message := "Naslov je obavezan."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.CategoryID == "" {
		message := "Kategorija je obavezna."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	if req.ContentDescription == "" {
		message := "Sadržaj je obavezan."

		return Render(ctx, http.StatusOK, components.ArticleError(message))
	}

	cookie, err := ctx.Cookie("access_token")
	if err != nil {
		// No cookie found; redirect to login page
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	accessToken := cookie.Value
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		// Invalid token; redirect to login page
		return ctx.Redirect(http.StatusTemporaryRedirect, "/")
	}

	parsedUserID, err := uuid.Parse(payload.UserID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid content ID format", err))
	}

	// Create a pgtype.UUID with the parsed UUID
	userID := pgtype.UUID{
		Bytes: parsedUserID,
		Valid: true,
	}

	parsedCategoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, errorResponse("invalid content ID format", err))
	}

	// Create a pgtype.UUID with the parsed UUID
	categoryID := pgtype.UUID{
		Bytes: parsedCategoryID,
		Valid: true,
	}

	arg := db.CreateContentParams{
		UserID:              userID,
		CategoryID:          categoryID,
		Title:               req.Title,
		ContentDescription:  req.ContentDescription,
		CommentsEnabled:     true,
		ViewCountEnabled:    true,
		LikeCountEnabled:    true,
		DislikeCountEnabled: false,
	}

	content, err := server.store.CreateContent(ctx.Request().Context(), arg)
	if err != nil {
		message := "Greška prilikom cuvanja sadržaja."

		return Render(ctx, http.StatusInternalServerError, components.ArticleError(message))
	}

	_, err = server.store.PublishContent(ctx.Request().Context(), content.ContentID)
	if err != nil {
		message := "Greška prilikom objavljivanja sadržaja."

		return Render(ctx, http.StatusInternalServerError, components.ArticleError(message))
	}

	ctx.SetCookie(&http.Cookie{
		Name:    "content_id",
		Value:   content.ContentID.String(),
		MaxAge:  60 * 60 * 24 * 365,
		Path:    "/",
		Expires: time.Now().Add(time.Hour),
	})

	message := "Uspešno ste sačuvali i objavili novi sadržaj."

	return Render(ctx, http.StatusOK, components.ArticleSuccess(message))
}

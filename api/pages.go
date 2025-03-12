package api

import (
	"log"
	"net/http"
	"time"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/token"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

var Loc, _ = time.LoadLocation("Europe/Belgrade")

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
		log.Println("Error binding request in adminCats:", err)
		return err
	}

	nextLimit := req.Limit + 10

	categories, err := server.store.ListCategories(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing categories in adminCats:", err)
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
		log.Println("Invalid category ID format in deleteCategoryModal:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	category, err := server.store.GetCategoryByID(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error getting category in deleteCategoryModal:", err)
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
		log.Println("Invalid category ID format in updateCategoryForm:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	category, err := server.store.GetCategoryByID(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Error getting category in updateCategoryForm:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.UpdateCategoryForm(category, updateCategoryErr))
}

type ListPublishedLimitReq struct {
	Limit int32 `query:"limit"`
}

func (server *Server) adminArts(ctx echo.Context) error {
	var req ListPublishedLimitReq

	overview, err := server.store.GetContentOverview(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting content overview in adminArts:", err)
		return err
	}

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimit(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing published content in adminArts:", err)
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

	return Render(ctx, http.StatusOK, components.AdminArticles(overview, int(nextLimit), content, url))
}

func (server *Server) publishedContentList(ctx echo.Context) error {
	var req ListPublishedLimitReq

	nextLimit := req.Limit + 20

	data, err := server.store.ListPublishedContentLimit(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing published content in publishedContentList:", err)
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

	return Render(ctx, http.StatusOK, components.PublishedContentSort(int(nextLimit), content, url))
}

func (server *Server) draftContentList(ctx echo.Context) error {
	var req ListPublishedLimitReq

	nextLimit := req.Limit + 20

	data, err := server.store.ListDraftContent(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing draft content in draftContentList:", err)
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

	return Render(ctx, http.StatusOK, components.DraftContentSort(int(nextLimit), content, url))
}

func (server *Server) deletedContentList(ctx echo.Context) error {
	var req ListPublishedLimitReq

	nextLimit := req.Limit + 20

	data, err := server.store.ListDeletedContent(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing deleted content in deletedContentList:", err)
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

	return Render(ctx, http.StatusOK, components.DeletedContentSort(int(nextLimit), content, url))
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

func (server *Server) createArticlePage(ctx echo.Context) error {
	categories, err := server.store.ListCategories(ctx.Request().Context(), 100)
	if err != nil {
		log.Println("Failed to get create article page in createArticlePage:", err)
		return err
	}

	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Failed to get tags for create article page in createArticlePage:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.CreateArticle(categories, tags))
}

func (server *Server) updateArticlePage(ctx echo.Context) error {
	contentIDStr := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedUUID, err := uuid.Parse(contentIDStr)
	if err != nil {
		log.Println("Invalid content ID format in updateArticlePage:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	pgUUID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	content, err := server.store.GetContentDetails(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Failed to get content for update article page in updateArticlePage:", err)
		return err
	}

	categories, err := server.store.ListCategories(ctx.Request().Context(), 100)
	if err != nil {
		log.Println("Failed to get update article page in updateArticlePage:", err)
		return err
	}

	media, err := server.store.ListMediaForContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Failed to get media for update article page in updateArticlePage:", err)
		return err
	}

	tags, err := server.store.ListTags(ctx.Request().Context(), 1000)
	if err != nil {
		log.Println("Failed to get tags for update article page in updateArticlePage:", err)
		return err
	}

	contentTags, err := server.store.GetTagsByContent(ctx.Request().Context(), pgUUID)
	if err != nil {
		log.Println("Failed to get tags for update article page in updateArticlePage:", err)
		return err
	}

	ctx.SetCookie(&http.Cookie{
		Name:   "content_id",
		Value:  contentIDStr,
		Path:   "/",
		MaxAge: 0,
	})

	return Render(ctx, http.StatusOK, components.UpdateArticle(content, categories, media, tags, contentTags))
}

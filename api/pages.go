package api

import (
	"log"
	"net/http"
	"time"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	"github.com/00mark0/macva-news/token"
	"github.com/00mark0/macva-news/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

var Loc, _ = time.LoadLocation("Europe/Belgrade")

func (server *Server) homePage(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.Index())
}

// full page to be served
func (server *Server) adminDash(ctx echo.Context) error {
	payload := ctx.Get(authorizationPayloadKey).(*token.Payload)

	userIDStr := payload.UserID

	userIDBytes, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("Error parsing user_id in adminDash:", err)
		return err
	}

	userID := pgtype.UUID{
		Bytes: userIDBytes,
		Valid: true,
	}

	user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error getting user in adminDash:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.DashPage(user))
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

type ListUsersLimitReq struct {
	Limit int32 `query:"limit"`
}

func (server *Server) adminUsers(ctx echo.Context) error {
	var req ListUsersLimitReq

	nextLimit := req.Limit + 20

	activeCount, err := server.store.GetActiveUsersCount(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting active users count in adminUsers:", err)
		return err
	}

	bannedCount, err := server.store.GetBannedUsersCount(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting banned users count in adminUsers:", err)
		return err
	}

	delCount, err := server.store.GetDeletedUsersCount(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting deleted users count in adminUsers:", err)
		return err
	}

	overview := components.UsersOverview{
		ActiveUsersCount:  int(activeCount),
		BannedUsersCount:  int(bannedCount),
		DeletedUsersCount: int(delCount),
	}

	var users []components.UsersRes

	data, err := server.store.GetActiveUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in adminUsers:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/active?limit="

	return Render(ctx, http.StatusOK, components.AdminUsers(overview, int(nextLimit), users, url))
}

func (server *Server) activeUsersList(ctx echo.Context) error {
	var req ListUsersLimitReq

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetActiveUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in activeUsersList:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/active?limit="

	return Render(ctx, http.StatusOK, components.ActiveUsersSort(int(nextLimit), users, url))
}

func (server *Server) bannedUsersList(ctx echo.Context) error {
	var req ListUsersLimitReq

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetBannedUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in bannedUsersList:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/banned?limit="

	return Render(ctx, http.StatusOK, components.BannedUsersSort(int(nextLimit), users, url))
}

func (server *Server) deletedUsersList(ctx echo.Context) error {
	var req ListUsersLimitReq

	nextLimit := req.Limit + 20

	var users []components.UsersRes

	data, err := server.store.GetDeletedUsers(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error getting active users in deletedUsersList:", err)
		return err
	}

	for _, v := range data {
		users = append(users, components.UsersRes{
			UserID:        v.UserID.String(),
			Username:      v.Username,
			Email:         v.Email,
			Pfp:           v.Pfp,
			Role:          v.Role,
			EmailVerified: v.EmailVerified.Bool,
			Banned:        v.Banned.Bool,
			IsDeleted:     v.IsDeleted.Bool,
			CreatedAt:     v.CreatedAt.Time.In(Loc).Format("02-01-06 15:04"),
		})
	}

	url := "/api/admin/users/deleted?limit="

	return Render(ctx, http.StatusOK, components.DelUsersSort(int(nextLimit), users, url))
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

func (server *Server) adminSettings(ctx echo.Context) error {
	// Get user data from auth payload
	payload := ctx.Get(authorizationPayloadKey).(*token.Payload)

	// Get global settings or create if they don't exist
	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil || len(globalSettings) == 0 {
		// If there's an error or no settings exist, create new settings
		newSettings, err := server.store.CreateGlobalSettings(ctx.Request().Context())
		if err != nil {
			log.Println("Error creating global settings in adminSettings:", err)
			return err
		}
		globalSettings = []db.GlobalSetting{newSettings}
	}

	userIDStr := payload.UserID

	userIDBytes, err := uuid.Parse(userIDStr)
	if err != nil {
		log.Println("Error parsing user_id in adminSettings:", err)
		return err
	}

	userID := pgtype.UUID{
		Bytes: userIDBytes,
		Valid: true,
	}

	user, err := server.store.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error getting user in adminSettings:", err)
		return err
	}

	// Create props for the AdminSettings component
	props := components.AdminSettingsProps{
		// User settings from auth payload
		UserID:   user.UserID.String(),
		Username: user.Username,
		Pfp:      user.Pfp,

		// Global settings from the first record
		DisableComments: globalSettings[0].DisableComments,
		DisableLikes:    globalSettings[0].DisableLikes,
		DisableDislikes: globalSettings[0].DisableDislikes,
		DisableViews:    globalSettings[0].DisableViews,
		DisableAds:      globalSettings[0].DisableAds,
	}

	// Render the AdminSettings component with the props
	return Render(ctx, http.StatusOK, components.AdminSettings(props))
}

func (server *Server) passwordResetPage(ctx echo.Context) error {
	token := ctx.Param("token")

	// Validate the token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return Render(ctx, http.StatusOK, components.PasswordReset("", "Link za resetovanje lozinke je nevažeći."))
	}

	// Verify that user_id exists in the claims
	if _, exists := claims["user_id"]; !exists {
		return Render(ctx, http.StatusOK, components.PasswordReset("", "Link za resetovanje lozinke je nevažeći."))
	}

	// Token is valid, show the password reset form
	return Render(ctx, http.StatusOK, components.PasswordReset(token, ""))
}

func (server *Server) registerPage(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.RegisterPage(""))
}

func (server *Server) emailVerifiedPage(ctx echo.Context) error {
	token := ctx.Param("token")

	// Validate the token
	claims, err := utils.ValidateToken(token)
	if err != nil {
		log.Println("Error validating token in emailVerifiedPage:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	// Verify that user_id exists in the claims
	if _, exists := claims["user_id"]; !exists {
		log.Println("Error extracting user_id from claims in emailVerifiedPage:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	// Get user ID from claims
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		log.Println("Error extracting user_id from claims in resetPassword:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	// Parse user ID to UUID
	var userID pgtype.UUID
	err = userID.Scan(userIDStr)
	if err != nil {
		log.Println("Error parsing user_id in resetPassword:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	err = server.store.SetEmailVerified(ctx.Request().Context(), userID)
	if err != nil {
		log.Println("Error setting email_verified in emailVerifiedPage:", err)
		return Render(ctx, http.StatusOK, components.VerificationError())
	}

	return Render(ctx, http.StatusOK, components.VerificationSuccess())
}

func (server *Server) requestPassResetPage(ctx echo.Context) error {
	return Render(ctx, http.StatusOK, components.RequestPassReset())
}

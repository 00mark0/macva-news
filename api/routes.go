package api

import (
	"github.com/00mark0/macva-news/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func (server *Server) setupRouter() {
	router := echo.New()

	router.Use(middleware.Gzip())

	// Initialize custom validator from validator.go
	router.Validator = NewCustomValidator()

	// Run cron job to create daily analytics
	go server.scheduleDailyAnalytics()

	// Serve static files
	router.Static("/static", "static")

	if os.Getenv("DEV_MODE") == "true" {
		router.Use(utils.NoCacheMiddleware)
	}

	// Set authRoutes instead of router to any routes that require middleware
	//authRoutes := router.Group("")
	//authRoutes.Use(authMiddleware(server.tokenMaker))

	adminRoutes := router.Group("")
	adminRoutes.Use(server.adminMiddleware(server.tokenMaker))

	// Admin Page Routes
	// Admin overview
	adminRoutes.GET("/admin", server.adminDash)
	adminRoutes.GET("/admin/hx-admin", server.adminDashContent)
	// Admin categories
	adminRoutes.GET("/admin/categories", server.adminCats)
	adminRoutes.GET("/admin/create-cat-form", server.createCategoryForm)
	adminRoutes.GET("/admin/delete-cat-modal/:id", server.deleteCategoryModal)
	adminRoutes.GET("/admin/update-cat-form/:id", server.updateCategoryForm)
	// Admin articles
	adminRoutes.GET("/admin/content", server.adminArts)
	adminRoutes.GET("/admin/content/create", server.createArticlePage)
	adminRoutes.GET("/admin/content/update/:id", server.updateArticlePage)
	adminRoutes.GET("/admin/pub-content", server.publishedContentList)
	adminRoutes.GET("/admin/draft-content", server.draftContentList)
	adminRoutes.GET("/admin/del-content", server.deletedContentList)
	// Admin users
	adminRoutes.GET("/admin/users", server.adminUsers)
	adminRoutes.GET("/admin/active-users", server.activeUsersList)
	adminRoutes.GET("/admin/banned-users", server.bannedUsersList)
	adminRoutes.GET("/admin/deleted-users", server.deletedUsersList)
	// Admin ads
	adminRoutes.GET("/admin/ads", server.adminAds)
	// Admin settings
	adminRoutes.GET("/admin/settings", server.adminSettings)

	// Admin API Routes
	// Admin overview
	adminRoutes.GET("/api/admin/trending", server.listTrendingContent)
	adminRoutes.GET("/api/admin/analytics", server.getDailyAnalytics)

	// Admin categories
	adminRoutes.GET("/api/admin/categories", server.listCats)
	adminRoutes.POST("/api/admin/category", server.createCategory)
	adminRoutes.DELETE("/api/admin/category/:id", server.deleteCategory)
	adminRoutes.PUT("/api/admin/category/:id", server.updateCategory)

	// Admin articles
	adminRoutes.GET("/api/admin/content/published", server.listPubContent)
	adminRoutes.GET("/api/admin/content/published/oldest", server.listPubContentOldest)
	adminRoutes.GET("/api/admin/content/published/title", server.listPubContentTitle)
	adminRoutes.GET("/api/admin/content/draft", server.listDraftContent)
	adminRoutes.GET("/api/admin/content/draft/oldest", server.listDraftContentOldest)
	adminRoutes.GET("/api/admin/content/draft/title", server.listDraftContentTitle)
	adminRoutes.GET("/api/admin/content/deleted", server.listDelContent)
	adminRoutes.GET("/api/admin/content/deleted/oldest", server.listDelContentOldest)
	adminRoutes.GET("/api/admin/content/deleted/title", server.listDelContentTitle)
	adminRoutes.GET("/api/admin/content/published/search", server.listSearchPubContent)
	adminRoutes.GET("/api/admin/content/draft/search", server.listSearchDraftContent)
	adminRoutes.GET("/api/admin/content/deleted/search", server.listSearchDelContent)
	adminRoutes.PUT("/api/admin/content/archive/:id", server.archivePubContent)
	adminRoutes.DELETE("/api/admin/content/:id", server.deleteContent)
	adminRoutes.PUT("/api/admin/content/publish/:id", server.publishDraftContent)
	adminRoutes.PUT("/api/admin/content/unarchive/:id", server.unarchiveContent)
	adminRoutes.PUT("/api/admin/content/:id", server.updateContent)
	adminRoutes.POST("/api/admin/content/draft", server.createContent)
	adminRoutes.POST("/api/admin/content/publish", server.createAndPublishContent)

	// Admin Media
	adminRoutes.GET("/api/admin/media", server.listMediaForContent)
	adminRoutes.POST("/api/admin/media/upload/new", server.addMediaToNewContent)
	adminRoutes.POST("/api/admin/media/upload/:id", server.addMediaToUpdateContent)
	adminRoutes.DELETE("/api/admin/media/remove/:id", server.deleteMedia)

	// Admin Tags
	adminRoutes.GET("/api/admin/tags", server.listTags)
	adminRoutes.GET("/api/admin/tags/search", server.listSearchTags)
	adminRoutes.GET("/api/admin/tags/:id", server.listTagsByContent)
	adminRoutes.POST("/api/admin/tags", server.createTag)
	adminRoutes.POST("/api/admin/tags/add", server.addTagToContent)
	adminRoutes.POST("/api/admin/tags/add/:id", server.addTagToContentUpdate)
	adminRoutes.DELETE("/api/admin/tags/content/remove/:id", server.removeTagFromContent)
	adminRoutes.DELETE("/api/admin/tags/content/remove/:content_id/:tag_id", server.removeTagFromContentUpdate)
	adminRoutes.DELETE("/api/admin/tags/remove/:id", server.deleteTag)

	// Admin Users
	adminRoutes.GET("/api/admin/users/active", server.listActiveUsers)
	adminRoutes.GET("/api/admin/users/active/oldest", server.listActiveUsersOldest)
	adminRoutes.GET("/api/admin/users/active/title", server.listActiveUsersTitle)
	adminRoutes.GET("/api/admin/users/banned", server.listBannedUsers)
	adminRoutes.GET("/api/admin/users/banned/oldest", server.listBannedUsersOldest)
	adminRoutes.GET("/api/admin/users/banned/title", server.listBannedUsersTitle)
	adminRoutes.GET("/api/admin/users/deleted", server.listDeletedUsers)
	adminRoutes.GET("/api/admin/users/deleted/oldest", server.listDeletedUsersOldest)
	adminRoutes.GET("/api/admin/users/deleted/title", server.listDeletedUsersTitle)
	adminRoutes.GET("/api/admin/users/active/search", server.searchActiveUsers)
	adminRoutes.GET("/api/admin/users/banned/search", server.searchBannedUsers)
	adminRoutes.GET("/api/admin/users/deleted/search", server.searchArchivedUsers)
	adminRoutes.PUT("/api/admin/users/ban/:id", server.banUser)
	adminRoutes.PUT("/api/admin/users/unban/:id", server.unbanUser)
	adminRoutes.PUT("/api/admin/users/archive/:id", server.deleteUser)

	// Cookie
	adminRoutes.DELETE("/api/cookie", server.deleteCookie)

	// Auth Pages
	router.GET("/login", server.loginPage)

	// Auth api
	router.POST("/api/login", server.login)
	adminRoutes.POST("/api/logout", server.logOut)

	// User Page Routes
	router.GET("/", server.homePage)
	router.GET("/counter", server.counterPage)
	router.GET("/widget", server.widgetPage)

	server.router = router
}

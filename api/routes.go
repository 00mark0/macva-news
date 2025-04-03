package api

import (
	"os"

	"github.com/00mark0/macva-news/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (server *Server) setupRouter() {
	router := echo.New()

	router.Use(middleware.Gzip())

	// Initialize custom validator from validator.go
	router.Validator = NewCustomValidator()

	// Run cron job to create daily analytics
	go server.scheduleDailyAnalytics()

	// Run cron job to deactivate expired ads
	go server.deactivateAds()

	// Serve static files
	router.Static("/static", "static")

	if os.Getenv("DEV_MODE") == "true" {
		router.Use(utils.NoCacheMiddleware)
	}

	// Set authRoutes instead of router to any routes that require middleware
	authRoutes := router.Group("")
	authRoutes.Use(server.authMiddleware(server.tokenMaker))

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
	adminRoutes.GET("/admin/active-ads", server.activeAdsList)
	adminRoutes.GET("/admin/inactive-ads", server.inactiveAdsList)
	adminRoutes.GET("/admin/scheduled-ads", server.scheduledAdsList)
	adminRoutes.GET("/admin/create-ad-modal", server.createAdModal)
	adminRoutes.GET("/admin/update-ad-modal/:id", server.updateAdModal)
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

	// Admin settings
	adminRoutes.PUT("/api/admin/settings/username/:id", server.updateUsername)
	adminRoutes.PUT("/api/admin/settings/pfp/:id", server.updatePfp)
	adminRoutes.PUT("/api/admin/global-settings", server.updateGlobalSettings)
	adminRoutes.PUT("/api/admin/reset-global-settings", server.resetGlobalSettings)

	// Admin ads
	adminRoutes.GET("/api/admin/ads/active", server.listActiveAds)
	adminRoutes.GET("/api/admin/ads/inactive", server.listInactiveAds)
	adminRoutes.POST("/api/admin/ads", server.createAd)
	adminRoutes.DELETE("/api/admin/ads/:id", server.deleteAd)
	adminRoutes.PUT("/api/admin/ads/:id", server.updateAd)
	adminRoutes.PUT("/api/admin/ads/deactivate/:id", server.deactivateAd)

	// Cookie
	adminRoutes.DELETE("/api/cookie", server.deleteCookie)

	// Auth Pages
	router.GET("/login", server.loginPage)
	router.GET("/register", server.registerPage)
	router.GET("/reset-lozinke/:token", server.passwordResetPage)

	// Auth api
	router.POST("/api/login", server.login)
	router.POST("/api/register", server.register)
	router.POST("/api/send-password-reset-form", server.requestPassResetFromForm)
	authRoutes.POST("/api/logout", server.logOut)
	authRoutes.POST("/api/send-password-reset", server.requestPassReset)
	router.POST("/api/reset-password", server.resetPassword)

	// User Page Routes
	// Home Page
	router.GET("/", server.homePage)
	// Email Verified Page
	router.GET("/potvrdi-email/:token", server.emailVerifiedPage)
	// Forgotten Password Page
	router.GET("/zaboravljena-lozinka", server.requestPassResetPage)
	// User Search
	router.GET("/search", server.searchResultsPage)
	// Categories Page
	router.GET("/kategorije/:category/:id", server.categoriesPage)

	//User API
	// User Search
	router.GET("/api/search", server.loadMoreSearch)
	router.GET("/api/content/other", server.listOtherContent)
	// Home
	router.GET("/api/news-slider", server.newsSlider)
	router.GET("/api/content/popular", server.listTrendingContentUser)
	router.GET("/api/content/categories", server.categoriesWithContent)
	// Categories
	router.GET("/api/category/content/recent/:id", server.listRecentCategoryContent)
	router.GET("/api/category/:id/tags/content", server.listContentByTagsUnderCategory)

	server.router = router
}

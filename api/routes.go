package api

import (
	"github.com/00mark0/macva-news/utils"
	"github.com/labstack/echo/v4"
	"os"
)

func (server *Server) setupRouter() {
	router := echo.New()
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
	adminRoutes.Use(adminMiddleware(server.tokenMaker))

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
	adminRoutes.GET("/admin/pub-content", server.publishedContentList)
	adminRoutes.GET("/admin/draft-content", server.draftContentList)
	adminRoutes.GET("/admin/del-content", server.deletedContentList)
	// Admin users
	adminRoutes.GET("/admin/users", server.adminUsers)
	// Admin ads
	adminRoutes.GET("/admin/ads", server.adminAds)

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

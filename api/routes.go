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

	// Pages
	// Admin Page Routes
	adminRoutes.GET("/admin", server.adminDash)
	adminRoutes.GET("/admin/hx-admin", server.adminDashContent)
	adminRoutes.GET("/admin/categories", server.adminCats)
	adminRoutes.GET("/admin/content", server.adminArts)
	adminRoutes.GET("/admin/users", server.adminUsers)
	adminRoutes.GET("/admin/ads", server.adminAds)
	router.GET("/admin/login", server.login)

	// Auth api
	router.POST("/api/admin/login", server.adminLogin)
	adminRoutes.POST("/api/admin/logout", server.logOut)

	// User
	router.GET("/", server.homePage)
	router.GET("/counter", server.counterPage)
	router.GET("/widget", server.widgetPage)

	server.router = router
}

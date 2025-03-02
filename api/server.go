package api

import (
	"github.com/00mark0/macva-news/db/services"
	"github.com/00mark0/macva-news/token"

	"fmt"

	"github.com/labstack/echo/v4"
)

// Server serves HTTP requests
type Server struct {
	store      *db.Store
	tokenMaker token.Maker
	router     *echo.Echo
}

// NewServer creates an HTTP server and sets up routing.
func NewServer(store *db.Store, symmetricKey string) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(symmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)

	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	go server.scheduleDailyAnalytics()

	return server, nil
}

// Start runs the HTTP server.
func (server *Server) Start(address string) error {
	return server.router.Start(address)
}

// errorResponse is a helper function to format error messages.
func errorResponse(message string, err error) echo.Map {
	return echo.Map{
		"message": message,
		"error":   err.Error(),
	}
}

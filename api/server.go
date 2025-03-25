package api

import (
	"os"

	"github.com/00mark0/macva-news/db/services"
	"github.com/00mark0/macva-news/token"

	"fmt"

	"github.com/labstack/echo/v4"
)

var BaseUrl = os.Getenv("BASE_URL")

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

	return server, nil
}

// Start runs the HTTP server.
func (server *Server) Start(address string) error {
	return server.router.Start(address)
}

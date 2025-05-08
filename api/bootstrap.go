package api

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/00mark0/macva-news/db/services"
	"github.com/00mark0/macva-news/utils"
)

func BootstrapAdmin(store *db.Store) error {
	ctx := context.Background()

	exists, err := store.CheckAdminExists(ctx)
	if err != nil {
		return fmt.Errorf("check admin existence failed: %w", err)
	}
	if exists {
		log.Println("Admin already exists.")
		return nil
	}

	hashedPassword, err := utils.HashPassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		return fmt.Errorf("hash password failed: %w", err)
	}

	arg := db.CreateUserAdminParams{
		Username: os.Getenv("ADMIN_USERNAME"),
		Email:    os.Getenv("EMAIL"),
		Password: hashedPassword,
		Role:     "admin",
	}

	admin, err := store.CreateUserAdmin(ctx, arg)
	if err != nil {
		return fmt.Errorf("create admin failed: %w", err)
	}

	log.Printf("Admin user created with ID %v", admin.UserID)
	return nil
}

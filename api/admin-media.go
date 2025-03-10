package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (server *Server) addMediaToNewContent(ctx echo.Context) error {
	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		log.Println("No content ID cookie found")
		return err
	}

	// Parse content ID from cookie
	contentIDString, err := uuid.Parse(contentIDCookie.Value)
	if err != nil {
		log.Println("Invalid content ID in cookie:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: contentIDString,
		Valid: true,
	}

	// Get the file from the form
	file, err := ctx.FormFile("file_upload")
	if err != nil {
		log.Println("Error retrieving uploaded file:", err)
		return err
	}

	uploadsDir := "static/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Println("Error creating uploads directory:", err)
		return err
	}

	// Generate a unique filename to avoid collisions
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	// Save the file to disk
	src, err := file.Open()
	if err != nil {
		log.Println("Error opening uploaded file:", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating destination file:", err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Error copying file data:", err)
		return err
	}

	// Determine media type based on file extension
	mediaType := "image" // Default
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == ".mp4" || ext == ".mov" || ext == ".avi" {
		mediaType = "video"
	} else if ext == ".mp3" || ext == ".wav" || ext == ".ogg" {
		mediaType = "audio"
	}

	existingMedia, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing existing media:", err)
		return err
	}

	nextOrder := int32(1)
	if len(existingMedia) > 0 {
		nextOrder = int32(len(existingMedia) + 1)
	}

	// Insert the media record into the database
	arg := db.InsertMediaParams{
		ContentID:    contentID,
		MediaType:    mediaType,
		MediaUrl:     "/" + filePath, // Store with leading slash for direct use in HTML
		MediaCaption: "",             // Empty caption by default
		MediaOrder:   nextOrder,
	}

	_, err = server.store.InsertMedia(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error inserting media record:", err)
		return err
	}

	updatedMedia, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing updated media:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.InsertMedia(updatedMedia, contentIDCookie.Value))
}

func (server *Server) listMediaForContent(ctx echo.Context) error {
	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		var emptyMedia []db.Medium

		return Render(ctx, http.StatusOK, components.InsertMedia(emptyMedia, ""))
	}

	// Parse content ID from cookie
	contentIDString, err := uuid.Parse(contentIDCookie.Value)
	if err != nil {
		log.Println("Invalid content ID in cookie:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: contentIDString,
		Valid: true,
	}

	media, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing media:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.InsertMedia(media, contentIDCookie.Value))
}

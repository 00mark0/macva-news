package api

import (
	"context"
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
		message := "Sadržaj nije pronađen. Dodajte medije sa uredi stranice."

		ctx.Response().Header().Set("HX-Retarget", "#create-article-modal")
		return Render(ctx, http.StatusOK, components.ArticleError(message))
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

func (server *Server) deleteMedia(ctx echo.Context) error {
	mediaIDStr := ctx.Param("id")

	// Parse media ID from the URL parameter
	mediaIDUUID, err := uuid.Parse(mediaIDStr)
	if err != nil {
		log.Println("Invalid media ID:", err)
		return err
	}

	mediaID := pgtype.UUID{
		Bytes: mediaIDUUID,
		Valid: true,
	}

	// Get the media record to find the file path before deleting
	media, err := server.store.GetMediaByID(ctx.Request().Context(), mediaID)
	if err != nil {
		log.Println("Error getting media record:", err)
		return err
	}

	// Get the content ID to use for rendering updated media list
	contentID := media.ContentID
	contentIDStr := contentID.String()

	// Remove the file from filesystem
	// The filepath is stored with leading slash, so trim it for filesystem operations
	filePath := strings.TrimPrefix(media.MediaUrl, "/")
	if err := os.Remove(filePath); err != nil {
		log.Println("Error removing file from filesystem:", err)
		// Continue with DB deletion even if file removal fails
	}

	// Delete the media record from the database
	if err := server.store.DeleteMedia(ctx.Request().Context(), mediaID); err != nil {
		log.Println("Error deleting media record:", err)
		return err
	}

	// Get updated media list for rendering
	updatedMedia, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing updated media:", err)
		return err
	}

	// Render the updated media list component
	return Render(ctx, http.StatusOK, components.InsertMedia(updatedMedia, contentIDStr))
}

func (server *Server) deleteMediaFunc(id string) error {
	// Parse media ID from the URL parameter
	mediaIDUUID, err := uuid.Parse(id)
	if err != nil {
		log.Println("Invalid media ID:", err)
		return err
	}

	mediaID := pgtype.UUID{
		Bytes: mediaIDUUID,
		Valid: true,
	}

	media, err := server.store.GetMediaByID(context.Background(), mediaID)
	if err != nil {
		log.Println("Error getting media record:", err)
		return err
	}

	// Remove the file from filesystem
	// The filepath is stored with leading slash, so trim it for filesystem operations
	filePath := strings.TrimPrefix(media.MediaUrl, "/")
	if err := os.Remove(filePath); err != nil {
		log.Println("Error removing file from filesystem:", err)
	}

	err = server.store.DeleteMedia(context.Background(), mediaID)
	if err != nil {
		log.Println("Error deleting media record:", err)
		return err
	}

	log.Println("Media record deleted successfully.")

	return nil
}

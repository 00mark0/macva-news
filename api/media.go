package api

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

// ConvertToWebPWithResize converts an image file to WebP format, resizing it to fit within maxWidth and maxHeight.
// quality is a value from 0 to 100.
func ConvertToWebPWithResize(inputPath string, maxWidth, maxHeight int, quality float32) (string, error) {
	// Open the input file
	file, err := os.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("error opening input file: %v", err)
	}
	defer file.Close()

	// Decode the image based on its extension
	var img image.Image
	ext := strings.ToLower(filepath.Ext(inputPath))
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		return inputPath, fmt.Errorf("unsupported image format: %s", ext)
	}
	if err != nil {
		return "", fmt.Errorf("error decoding image: %v", err)
	}

	// Resize the image if maxWidth or maxHeight is specified (> 0)
	if maxWidth > 0 || maxHeight > 0 {
		// imaging.Fit maintains aspect ratio and fits within the given bounds.
		img = imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
	}

	// Generate WebP filename (replace original extension with .webp)
	webpPath := strings.TrimSuffix(inputPath, ext) + ".webp"

	// Create WebP output file
	output, err := os.Create(webpPath)
	if err != nil {
		return "", fmt.Errorf("error creating WebP output file: %v", err)
	}
	defer output.Close()

	// Encode to WebP (Lossy conversion with specified quality)
	if err := webp.Encode(output, img, &webp.Options{Lossless: false, Quality: quality}); err != nil {
		return "", fmt.Errorf("error encoding to WebP: %v", err)
	}

	// Optionally: Remove the original file if you don't need it
	if err := os.Remove(inputPath); err != nil {
		log.Printf("Warning: could not remove original file %s: %v", inputPath, err)
	}

	return webpPath, nil
}

func (server *Server) addMediaToNewContent(ctx echo.Context) error {
	contentIDCookie, err := ctx.Cookie("content_id")
	if err != nil {
		var emptyMedia []db.Medium

		return Render(ctx, http.StatusOK, components.InsertMedia(emptyMedia, ""))
	}

	contentIDStr := contentIDCookie.Value

	// Parse string UUID into proper UUID format
	parsedUUID, err := uuid.Parse(contentIDStr)
	if err != nil {
		log.Println("Invalid content ID format from cookie in addMediaToNewContent:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	// Get the file from the form
	file, err := ctx.FormFile("file_upload")
	if err != nil {
		log.Println("Error retrieving uploaded file in addMediaToNewContent:", err)
		return err
	}

	uploadsDir := "static/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Println("Error creating uploads directory in addMediaToNewContent:", err)
		return err
	}

	// Generate a unique filename to avoid collisions
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	// Save the file to disk
	src, err := file.Open()
	if err != nil {
		log.Println("Error opening uploaded file in addMediaToNewContent:", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating destination file in addMediaToNewContent:", err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Error copying file data in addMediaToNewContent:", err)
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

	// If the file is an image, convert it to WebP
	if mediaType == "image" {
		convertedPath, err := ConvertToWebPWithResize(filePath, 800, 600, 80)
		if err != nil {
			log.Println("Error converting image to WebP:", err)
			// Optionally, you might want to continue with the original file in case of error
		} else {
			// Update filePath to the new WebP file.
			filePath = convertedPath
		}
	}

	existingMedia, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing existing media in addMediaToNewContent:", err)
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
		log.Println("Error inserting media record in addMediaToNewContent:", err)
		return err
	}

	updatedMedia, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing updated media in addMediaToNewContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.InsertMedia(updatedMedia, contentID.String()))
}

func (server *Server) addMediaToUpdateContent(ctx echo.Context) error {
	contentIDStr := ctx.Param("id")

	// Parse string UUID into proper UUID format
	parsedUUID, err := uuid.Parse(contentIDStr)
	if err != nil {
		log.Println("Invalid content ID format in addMediaToUpdateContent:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	// Get the file from the form
	file, err := ctx.FormFile("file_upload")
	if err != nil {
		log.Println("Error retrieving uploaded file in addMediaToUpdateContent:", err)
		return err
	}

	uploadsDir := "static/uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Println("Error creating uploads directory in addMediaToUpdateContent:", err)
		return err
	}

	// Generate a unique filename to avoid collisions
	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	// Save the file to disk
	src, err := file.Open()
	if err != nil {
		log.Println("Error opening uploaded file in addMediaToUpdateContent:", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating destination file in addMediaToUpdateContent:", err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Error copying file data in addMediaToUpdateContent:", err)
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

	// If the file is an image, convert it to WebP
	if mediaType == "image" {
		convertedPath, err := ConvertToWebPWithResize(filePath, 800, 600, 80)
		if err != nil {
			log.Println("Error converting image to WebP:", err)
			// Optionally, you might want to continue with the original file in case of error
		} else {
			// Update filePath to the new WebP file.
			filePath = convertedPath
		}
	}

	existingMedia, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing existing media in addMediaToUpdateContent:", err)
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

	media, err := server.store.InsertMedia(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error inserting media record in addMediaToUpdateContent:", err)
		return err
	}

	if media.MediaOrder == 1 {
		arg := db.AddThumbnailParams{
			ContentID: contentID,
			Thumbnail: pgtype.Text{String: "/" + filePath, Valid: true},
		}

		_, err := server.store.AddThumbnail(ctx.Request().Context(), arg)
		if err != nil {
			log.Println("Error adding thumbnail in addMediaToUpdateContent:", err)
			return err
		}
	}

	updatedMedia, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing updated media in addMediaToUpdateContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.InsertMediaUpdate(updatedMedia, contentID.String()))
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
		log.Println("Invalid content ID in cookie in listMediaForContent:", err)
		return err
	}

	// Create a pgtype.UUID with the parsed UUID
	contentID := pgtype.UUID{
		Bytes: contentIDString,
		Valid: true,
	}

	media, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing media for content in listMediaForContent:", err)
		return err
	}

	return Render(ctx, http.StatusOK, components.InsertMedia(media, contentIDCookie.Value))
}

func (server *Server) deleteMedia(ctx echo.Context) error {
	mediaIDStr := ctx.Param("id")

	// Parse media ID from the URL parameter
	mediaIDUUID, err := uuid.Parse(mediaIDStr)
	if err != nil {
		log.Println("Invalid media ID in deleteMedia:", err)
		return err
	}

	mediaID := pgtype.UUID{
		Bytes: mediaIDUUID,
		Valid: true,
	}

	// Get the media record to find the file path before deleting
	media, err := server.store.GetMediaByID(ctx.Request().Context(), mediaID)
	if err != nil {
		log.Println("Error getting media record in deleteMedia:", err)
		return err
	}

	// Get the content ID to use for rendering updated media list
	contentID := media.ContentID
	contentIDStr := contentID.String()

	// Remove the file from filesystem
	// The filepath is stored with leading slash, so trim it for filesystem operations
	filePath := strings.TrimPrefix(media.MediaUrl, "/")
	if err := os.Remove(filePath); err != nil {
		log.Println("Error removing file from filesystem in deleteMedia:", err)
		// Continue with DB deletion even if file removal fails
	}

	// Delete the media record from the database
	if err := server.store.DeleteMedia(ctx.Request().Context(), mediaID); err != nil {
		log.Println("Error deleting media record in deleteMedia:", err)
		return err
	}

	// Get updated media list for rendering
	updatedMedia, err := server.store.ListMediaForContent(ctx.Request().Context(), contentID)
	if err != nil {
		log.Println("Error listing updated media in deleteMedia:", err)
		return err
	}

	// Render the updated media list component
	return Render(ctx, http.StatusOK, components.InsertMediaUpdate(updatedMedia, contentIDStr))
}

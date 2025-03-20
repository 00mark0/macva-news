package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

func (server *Server) listActiveAds(ctx echo.Context) error {
	var req ListAdsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listActiveAds:", err)
		return err
	}

	nextLimit := req.Limit + 20

	activeAds, err := server.store.ListActiveAds(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing active ads in listActiveAds:", err)
		return err
	}

	url := "/api/admin/ads/active?limit="

	return Render(ctx, http.StatusOK, components.Ads(int(nextLimit), activeAds, url))
}

func (server *Server) listInactiveAds(ctx echo.Context) error {
	var req ListAdsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listInactiveAds:", err)
		return err
	}

	nextLimit := req.Limit + 20

	inactiveAds, err := server.store.ListInactiveAds(ctx.Request().Context(), nextLimit)
	if err != nil {
		log.Println("Error listing inactive ads in listInactiveAds:", err)
		return err
	}

	url := "/api/admin/ads/inactive?limit="

	return Render(ctx, http.StatusOK, components.Ads(int(nextLimit), inactiveAds, url))
}

type CreateAdReq struct {
	Title       string `form:"title" validate:"required,min=3,max=50"`
	Description string `form:"description" validate:"required,min=3,max=100"`
	TargetUrl   string `form:"target_url" validate:"required"`
	Placement   string `form:"placement" validate:"required"`
	Status      string `form:"status" validate:"required"`
	StartDate   string `form:"start_date" validate:"required"`
	EndDate     string `form:"end_date" validate:"required"`
}

func (server *Server) createAd(ctx echo.Context) error {
	var req CreateAdReq
	var createAddErr components.CreateAdErr

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in createAd:", err)
		return err
	}

	if req.TargetUrl != "" {
		// Check if the URL already has a protocol
		if !strings.HasPrefix(strings.ToLower(req.TargetUrl), "http://") &&
			!strings.HasPrefix(strings.ToLower(req.TargetUrl), "https://") {
			// Add https:// if no protocol is present
			req.TargetUrl = "https://" + req.TargetUrl
		}
	}

	if err := ctx.Validate(req); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			switch fieldErr.Field() {
			case "Title":
				createAddErr = "Naslov oglasa mora biti između 3 i 50 karaktera."
			case "Description":
				createAddErr = "Opis oglasa mora biti između 3 i 100 karaktera."
			case "TargetUrl":
				createAddErr = "URL oglasa je obavezan."
			case "Placement":
				createAddErr = "Mesto oglasa je obavezno."
			case "Status":
				createAddErr = "Status oglasa je obavezan."
			case "StartDate":
				createAddErr = "Datum početka oglasa je obavezan."
			case "EndDate":
				createAddErr = "Datum završetka oglasa je obavezan."
			}
		}

		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	file, err := ctx.FormFile("image_url")
	if err != nil {
		createAddErr = "Slika oglasa je obavezna."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	uploadsDir := "static/ads"
	// Ensure directory exists
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			log.Println("Error creating directory in updatePfp:", err)
			return err
		}
	}

	filename := fmt.Sprintf("%s-%s", uuid.New().String(), file.Filename)
	filePath := fmt.Sprintf("%s/%s", uploadsDir, filename)

	src, err := file.Open()
	if err != nil {
		log.Println("Error opening uploaded file in createAd:", err)
		return err
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		log.Println("Error creating destination file in createAd:", err)
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		log.Println("Error copying file in createAd:", err)
		return err
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		log.Println("Error parsing start date in createAd:", err)
		return err
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		log.Println("Error parsing end date in createAd:", err)
		return err
	}

	// 1. Start date must be before end date.
	if startDate.After(endDate) {
		createAddErr = "Datum početka oglasa mora biti pre datuma završetka."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 2. Start date must not be in the past (you might consider allowing today by comparing to midnight or adding a small margin).
	if startDate.Before(time.Now()) {
		createAddErr = "Datum početka oglasa mora biti veći od trenutnog datuma."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 3. End date must not be in the past.
	if endDate.Before(time.Now().Add(time.Hour * 1)) {
		createAddErr = "Datum završetka oglasa mora biti veći od trenutnog datuma."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 4. Start date must not be more than one year in the future.
	maxFutureDate := time.Now().AddDate(1, 0, 0)
	if startDate.After(maxFutureDate) {
		createAddErr = "Datum početka oglasa ne moze biti više od godinu dana unapred."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 5. The duration between start and end must be at least 3 days.
	minDuration := 3 * 24 * time.Hour
	if endDate.Sub(startDate) < minDuration {
		createAddErr = "Razmak između početka i kraja oglasa mora biti barem 3 dana."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	// 6. End date must not be more than 5 years in the future.
	maxEndDate := time.Now().AddDate(5, 0, 0)
	if endDate.After(maxEndDate) {
		createAddErr = "Datum kraja oglasa ne može biti više od 5 godina u budućnosti."
		ctx.Response().Header().Set("HX-Retarget", "#create-ad-modal")
		return Render(ctx, http.StatusOK, components.CreateAdModal(createAddErr))
	}

	arg := db.CreateAdParams{
		Title: pgtype.Text{String: req.Title, Valid: true},
		Description: pgtype.Text{
			String: req.Description,
			Valid:  true,
		},
		ImageUrl: pgtype.Text{
			String: filePath,
			Valid:  true,
		},
		TargetUrl: pgtype.Text{String: req.TargetUrl, Valid: true},
		Placement: pgtype.Text{String: req.Placement, Valid: true},
		Status:    pgtype.Text{String: req.Status, Valid: true},
		StartDate: pgtype.Timestamptz{
			Time:  startDate,
			Valid: true,
		},
		EndDate: pgtype.Timestamptz{
			Time:  endDate,
			Valid: true,
		},
	}

	_, err = server.store.CreateAd(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error creating ad in createAd:", err)
		return err
	}

	if req.Status == "active" {
		ctx.Response().Header().Set("HX-Trigger", "createAdSuccess")
		return server.activeAdsList(ctx)
	} else {
		ctx.Response().Header().Set("HX-Trigger", "createAdSuccess")
		return server.inactiveAdsList(ctx)
	}
}

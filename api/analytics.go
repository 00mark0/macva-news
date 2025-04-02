package api

import (
	"log"
	"net/http"
	"time"

	"github.com/00mark0/macva-news/components"
	"github.com/00mark0/macva-news/db/services"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
)

type TrendingContentReq struct {
	PublishedAt string `query:"published_at"`
	Limit       int32  `query:"limit"`
}

type TrendingContentRes struct {
	ContentID           string `json:"content_id"`
	UserID              string `json:"user_id"`
	CategoryID          string `json:"category_id"`
	Title               string `json:"title"`
	Thumbnail           string `json:"thumbnail"`
	ContentDescription  string `json:"content_description"`
	CommentsEnabled     bool   `json:"comments_enabled"`
	ViewCountEnabled    bool   `json:"view_count_enabled"`
	LikeCountEnabled    bool   `json:"like_count_enabled"`
	DislikeCountEnabled bool   `json:"dislike_count_enabled"`
	Status              string `json:"status"`
	ViewCount           int    `json:"view_count"`
	LikeCount           int    `json:"like_count"`
	DislikeCount        int    `json:"dislike_count"`
	CommentCount        int    `json:"comment_count"`
	CreatedAt           string `json:"created_at"`
	UpdatedAt           string `json:"updated_at"`
	PublishedAt         string `json:"published_at"`
	IsDeleted           bool   `json:"is_deleted"`
	TotalInteractions   int    `json:"total_interactions"`
}

func (server *Server) listTrendingContent(ctx echo.Context) error {
	var req TrendingContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listTrendingContent:", err)
		return err
	}

	// Parse the incoming date (expected format: "2006-01-02")
	publishedDate, err := time.Parse("2006-01-02", req.PublishedAt)
	if err != nil {
		log.Println("Error parsing date in listTrendingContent:", err)
		return err
	}

	// Set the time to midnight in the +0100 timezone.
	// Using time.FixedZone creates a fixed location with the desired offset.
	loc := time.FixedZone("UTC", 0) // +1 hour offset
	publishedDate = time.Date(publishedDate.Year(), publishedDate.Month(), publishedDate.Day(), 0, 0, 0, 0, loc)

	// Convert to pgtype.Timestamptz so we can pass it to the database.
	var publishedTimestamptz pgtype.Timestamptz
	publishedTimestamptz = pgtype.Timestamptz{
		Time:  publishedDate,
		Valid: true,
	}

	arg := db.ListTrendingContentParams{
		PublishedAt: publishedTimestamptz,
		Limit:       req.Limit,
	}

	data, err := server.store.ListTrendingContent(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error listing trending content in listTrendingContent:", err)
		return err
	}

	var trendingContent []TrendingContentRes

	for _, content := range data {
		trendingContent = append(trendingContent, TrendingContentRes{
			ContentID:  content.ContentID.String(),
			UserID:     content.UserID.String(),
			CategoryID: content.CategoryID.String(),
			Title:      content.Title,
			Thumbnail: func() string {
				if content.Thumbnail.Valid && content.Thumbnail.String != "" {
					return content.Thumbnail.String
				}

				return ThumbnailURL
			}(),
			ContentDescription:  content.ContentDescription,
			CommentsEnabled:     content.CommentsEnabled,
			ViewCountEnabled:    content.ViewCountEnabled,
			LikeCountEnabled:    content.LikeCountEnabled,
			DislikeCountEnabled: content.DislikeCountEnabled,
			Status:              content.Status,
			ViewCount:           int(content.ViewCount),
			LikeCount:           int(content.LikeCount),
			DislikeCount:        int(content.DislikeCount),
			CommentCount:        int(content.CommentCount),
			CreatedAt:           content.CreatedAt.Time.Format("2006-01-02"),
			UpdatedAt:           content.UpdatedAt.Time.Format("2006-01-02"),
			PublishedAt:         content.PublishedAt.Time.Format("2006-01-02"),
			IsDeleted:           content.IsDeleted.Bool,
			TotalInteractions:   int(content.TotalInteractions),
		})
	}

	if acceptHeader := ctx.Request().Header.Get("Accept"); acceptHeader == "application/json" {
		return ctx.JSON(http.StatusOK, trendingContent)
	}

	return nil
}

type DailyAnalyticsReq struct {
	AnalyticsDate  string `query:"start_date"`
	AnalyticsDate2 string `query:"end_date"`
	Limit          int32  `query:"limit"`
}

type DailyAnalyticsRes struct {
	AnalytycsDate  string `json:"analytics_date"`
	TotalViews     int    `json:"total_views"`
	TotalLikes     int    `json:"total_likes"`
	TotalDislikes  int    `json:"total_dislikes"`
	TotalComments  int    `json:"total_comments"`
	TotalAdsClicks int    `json:"total_ads_clicks"`
	CreatedAt      string `json:"created_at"`
}

func (server *Server) getDailyAnalytics(ctx echo.Context) error {
	var req DailyAnalyticsReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in getDailyAnalytics:", err)
		return err
	}

	analyticsDate, err := time.Parse("2006-01-02", req.AnalyticsDate)
	if err != nil {
		log.Println("Error parsing date in getDailyAnalytics:", err)
		return err
	}

	analyticsDate2, err := time.Parse("2006-01-02", req.AnalyticsDate2)
	if err != nil {
		log.Println("Error parsing date in getDailyAnalytics:", err)
		return err
	}

	arg := db.GetDailyAnalyticsParams{
		AnalyticsDate:   pgtype.Date{Time: analyticsDate, Valid: true},
		AnalyticsDate_2: pgtype.Date{Time: analyticsDate2, Valid: true},
		Limit:           req.Limit,
	}

	data, err := server.store.GetDailyAnalytics(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error getting daily analytics in getDailyAnalytics:", err)
		return err
	}

	var dailyAnalytics []DailyAnalyticsRes

	for _, analytics := range data {
		dailyAnalytics = append(dailyAnalytics, DailyAnalyticsRes{
			AnalytycsDate:  analytics.AnalyticsDate.Time.Format("2006-01-02"),
			TotalViews:     int(analytics.TotalViews),
			TotalLikes:     int(analytics.TotalLikes),
			TotalDislikes:  int(analytics.TotalDislikes),
			TotalComments:  int(analytics.TotalComments),
			TotalAdsClicks: int(analytics.TotalAdsClicks),
			CreatedAt:      analytics.CreatedAt.Time.Format("2006-01-02"),
		})
	}

	if acceptHeader := ctx.Request().Header.Get("Accept"); acceptHeader == "application/json" {
		return ctx.JSON(http.StatusOK, dailyAnalytics)
	}

	return nil
}

func (server *Server) listTrendingContentUser(ctx echo.Context) error {
	var req TrendingContentReq

	if err := ctx.Bind(&req); err != nil {
		log.Println("Error binding request in listTrendingContent:", err)
		return err
	}

	// Calculate the published date as 30 days ago from now
	publishedDate := time.Now().In(Loc).AddDate(0, 0, -30)

	// Adjust the timezone here if needed.
	publishedDate = time.Date(publishedDate.Year(), publishedDate.Month(), publishedDate.Day(), 0, 0, 0, 0, Loc)

	// Convert to pgtype.Timestamptz so we can pass it to the database.
	var publishedTimestamptz pgtype.Timestamptz
	publishedTimestamptz = pgtype.Timestamptz{
		Time:  publishedDate,
		Valid: true,
	}

	nextLimit := req.Limit + 6

	arg := db.ListTrendingContentParams{
		PublishedAt: publishedTimestamptz,
		Limit:       nextLimit,
	}

	data, err := server.store.ListTrendingContent(ctx.Request().Context(), arg)
	if err != nil {
		log.Println("Error listing trending content in listTrendingContent:", err)
		return err
	}

	var trendingContent []components.ContentData

	for _, content := range data {
		trendingContent = append(trendingContent, components.ContentData{
			ContentID:    content.ContentID,
			UserID:       content.UserID,
			CategoryID:   content.CategoryID,
			CategoryName: content.CategoryName,
			Title:        content.Title,
			Thumbnail: func() pgtype.Text {
				if content.Thumbnail.Valid && content.Thumbnail.String != "" {
					return content.Thumbnail
				}

				return pgtype.Text{String: ThumbnailURL, Valid: true}
			}(),
			ContentDescription:  content.ContentDescription,
			CommentsEnabled:     content.CommentsEnabled,
			ViewCountEnabled:    content.ViewCountEnabled,
			LikeCountEnabled:    content.LikeCountEnabled,
			DislikeCountEnabled: content.DislikeCountEnabled,
			Status:              content.Status,
			ViewCount:           content.ViewCount,
			LikeCount:           content.LikeCount,
			DislikeCount:        content.DislikeCount,
			CommentCount:        content.CommentCount,
			CreatedAt:           content.CreatedAt,
			UpdatedAt:           content.UpdatedAt,
			PublishedAt:         content.PublishedAt,
			IsDeleted:           content.IsDeleted,
			TotalInteractions:   content.TotalInteractions,
		})
	}

	globalSettings, err := server.store.GetGlobalSettings(ctx.Request().Context())
	if err != nil {
		log.Println("Error getting global settings in listTrendingContent:", err)
		return err
	}

	url := "/api/content/popular?limit="

	target := "#popular-content"

	title := "Popularno"

	return Render(ctx, http.StatusOK, components.GridCards(trendingContent, globalSettings[0], int(nextLimit), url, target, title))
}

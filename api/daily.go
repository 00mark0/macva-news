package api

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robfig/cron/v3"
)

func (server *Server) scheduleDailyAnalytics() {
	// Create a new cron scheduler (uses the local time zone by default)
	c := cron.New(cron.WithLocation(Loc))

	// Schedule the job to run every day at midnight.
	// "@daily" is equivalent to "0 0 0 * * *"
	var err error
	_, err = c.AddFunc("@daily", func() {
		// Use time.Now() to set the current day at midnight
		now := time.Now()
		date := pgtype.Date{
			Time:  time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
			Valid: true,
		}

		if _, err := server.store.CreateDailyAnalytics(context.Background(), date); err != nil {
			log.Printf("Failed to create daily analytics: %v\n", err)
		} else {
			log.Println("Daily analytics created successfully at midnight.")
		}
	})
	if err != nil {
		log.Fatalf("Error scheduling daily analytics: %v\n", err)
	}

	// Start the cron scheduler in its own goroutine
	c.Start()
}

func (server *Server) deactivateAds() {
	// Create a new cron scheduler (uses the local time zone by default)
	c := cron.New(cron.WithLocation(Loc))

	// Schedule the job to run every day at midnight.
	// "@daily" is equivalent to "0 0 0 * * *"
	var err error
	_, err = c.AddFunc("@daily", func() {
		// Get the current time
		now := time.Now()

		// Create a context with timeout for database operations
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// List all active ads
		ads, err := server.store.ListAds(ctx, 1000)
		if err != nil {
			log.Printf("Failed to list ads: %v\n", err)
			return
		}

		for _, ad := range ads {
			// Check if the ad is expired
			// Assuming EndDate is pgtype.Timestamptz, convert it to time.Time for comparison
			if ad.EndDate.Valid && ad.EndDate.Time.Before(now) {
				_, err := server.store.DeactivateAd(ctx, ad.ID)
				if err != nil {
					log.Printf("Failed to deactivate ad %v: %v\n", ad.ID, err)
				} else {
					log.Printf("Successfully deactivated expired ad %v\n", ad.ID)
				}
			}
		}
	})

	if err != nil {
		log.Fatalf("Error setting up cron job for deactivating expired ads: %v\n", err)
	}

	// Start the cron scheduler in its own goroutine
	c.Start()
}

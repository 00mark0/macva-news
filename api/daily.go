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
	c := cron.New(cron.WithLocation(time.Local))

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

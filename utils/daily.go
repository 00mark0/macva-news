package utils

import (
	"log"
	"time"
)

func ScheduleMidnightTask() {
	for {
		now := time.Now()
		// Calculate next midnight
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		durationUntilMidnight := nextMidnight.Sub(now)
		log.Printf("Sleeping for %v until midnight...\n", durationUntilMidnight)
		time.Sleep(durationUntilMidnight)

		// Execute the task
		log.Println("Executing task at midnight...")
	}
}

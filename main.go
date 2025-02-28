package main

import (
	"log"
	"os"

	"github.com/00mark0/macva-news/api"
	"github.com/00mark0/macva-news/db/services"

	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	queries *db.Queries
	store   *db.Store
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	dbSource := os.Getenv("DB_URL")
	port := os.Getenv("SERVER_PORT")
	symmetricKey := os.Getenv("TOKEN_SYMMETRIC_KEY")

	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection pool created.")

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Cannot connect to db!:", err)
	}
	log.Println("Connected to db.")

	store := db.NewStore(conn)
	server, err := api.NewServer(store, symmetricKey)
	if err != nil {
		log.Fatal("Cannot create server!:", err)
	}

	err = server.Start(":" + port)
	if err != nil {
		log.Fatal("Cannot start server!:", err)
	}
}

package models

import (
	"context"
	"fmt"
	"log"

	"gin-crud/ent"
	"gin-crud/internal/config"

	"entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

var Client *ent.Client

func InitDB() {
	config.LoadEnv()

	dbHost := config.GetEnv("DB_HOST")
	dbPort := config.GetEnv("DB_PORT")
	dbUser := config.GetEnv("DB_USER")
	dbPassword := config.GetEnv("DB_PASSWORD")
	dbName := config.GetEnv("DB_NAME")
	dbSSLMode := config.GetEnv("DB_SSLMODE")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	drv, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	Client = ent.NewClient(ent.Driver(drv))

	ctx := context.Background()
	if err := Client.Schema.Create(ctx); err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	seedItems(ctx)
}

func seedItems(ctx context.Context) {
	count, err := Client.Item.Query().Count(ctx)
	if err != nil {
		log.Fatalf("Failed to count items: %v", err)
	}

	if count == 0 {
		// Define your items here with all fields except ID
		items := []struct {
			Name        string
			Price       int
			Description string
			Stock       int
		}{
			{"Keyboard", 199, "A mechanical keyboard", 10},
			{"Screen", 299, "A 24-inch monitor", 5},
			{"Server", 599, "A high-performance server", 2},
			{"Printer", 399, "A color laser printer", 8},
		}

		// Loop through and insert each item
		for _, item := range items {
			_, err := Client.Item.
				Create().
				SetName(item.Name).
				SetPrice(item.Price).
				SetDescription(item.Description).
				Save(ctx)

			if err != nil {
				log.Fatalf("Failed to seed item: %v", err)
			}
		}
		log.Println("Successfully seeded items")
	}
}

func CloseDB() {
	if Client != nil {
		if err := Client.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}
}

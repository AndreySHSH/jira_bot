package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"jira_bot/internal/bot"
	"jira_bot/internal/repositories"
	"jira_bot/internal/services"
)

func main() {
	token := os.Getenv("TG_BOT_TOKEN")
	jiraURL := os.Getenv("JIRA_URL")

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow`,
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DBNAME"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}

	r, err := repositories.NewRepositories(db)
	if err != nil {
		panic(err)
	}

	s := services.NewServices(r)

	b, err := bot.NewBot(token, jiraURL, s)
	if err != nil {
		panic(err)
	}

	b.Run()
}

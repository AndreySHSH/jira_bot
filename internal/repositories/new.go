package repositories

import (
	"gorm.io/gorm"
	"jira_bot/internal/repositories/user"
)

type Repositories struct {
	db *gorm.DB

	User *user.Repository
}

func NewRepositories(db *gorm.DB) (*Repositories, error) {
	err := db.AutoMigrate(
		&user.User{},
	)
	if err != nil {
		return nil, err
	}

	userRepository := user.NewUserRepository(db)

	return &Repositories{
		User: userRepository,
	}, nil
}

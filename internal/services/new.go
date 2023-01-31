package services

import (
	"jira_bot/internal/repositories"
	"jira_bot/internal/services/user"
)

type Services struct {
	repo *repositories.Repositories

	User *user.Service
}

func NewServices(repo *repositories.Repositories) *Services {
	userService := user.NewUserService(repo.User)

	return &Services{
		User: userService,
	}
}

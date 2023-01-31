package user

import "jira_bot/internal/repositories/user"

type Service struct {
	userRepository *user.Repository
}

func (s Service) CreateUser(id int64, login string, password string) {
	u := s.userRepository.CheckUser(login)
	if u {
		s.userRepository.UpdateUser(id, login, password)
		return
	}
	s.userRepository.CreateUser(id, login, password)
}

func (s Service) GetUserByID(id int64) *user.User {
	return s.userRepository.GetUserByID(id)
}

func NewUserService(userRepository *user.Repository) *Service {
	return &Service{
		userRepository: userRepository,
	}
}

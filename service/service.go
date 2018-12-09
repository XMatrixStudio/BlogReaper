package service

type Service struct {
	User userService
}

func (s *Service) GetUserService() UserService {
	return &s.User
}

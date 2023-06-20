package user

import (
	"context"
	"rest-api-tutorial/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func (s *Service) Create(ctx context.Context, dto CreateUserDto) (u User, err error) {
	// todo for next one
	return
}

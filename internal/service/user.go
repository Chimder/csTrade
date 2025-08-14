package service

import (
	"context"
	"csTrade/internal/domain/user"
	"csTrade/internal/repository"

	"github.com/rs/zerolog/log"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (of *UserService) CreateUser(ctx context.Context, req *user.UserCreateReq) error {
	log.Info().Msg("createUser")
	err := of.repo.User.CreateUser(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

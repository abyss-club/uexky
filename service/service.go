package service

import (
	"context"
)

type Repository struct {
	User UserRepository
}

type Service struct {
	Repo Repository
}

func NewService(repo Repository) Service {
	return Service{Repo: repo}
}

func (s *Service) GetUser(ctx context.Context, id int) (*User, error) {
	users, err := s.Repo.User.FindUsersByIDs(ctx, []int{id})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (s *Service) SetName(ctx context.Context, id int, name *string) (*User, error) {
	if err := s.Repo.User.UpdateUser(ctx, &User{ID: id, Name: name}); err != nil {
		return nil, err
	}
	return s.GetUser(ctx, id)
}

func (s *Service) SetLevel(ctx context.Context, id int, level int) (*User, error) {
	if err := s.Repo.User.UpdateUser(ctx, &User{ID: id, Level: level}); err != nil {
		return nil, err
	}
	return s.GetUser(ctx, id)
}

package service

import (
	"context"

	"tasks/db"
)

func (s *task) RegisterUser(ctx context.Context, args *db.User) (string, error) {
	err := s.db.CreateUser(args)
	if err != nil {
		return "", err
	}
	return args.Username, nil

}

func (s *task) GetUser(ctx context.Context, username string) (*db.User, error) {
	user, err := s.db.GetUser(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

package service

import (
	"context"
	"database/sql"
	"errors"
	"tasks/db"
	sqldb "tasks/db"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAlreadyExists     = errors.New("note already exists")
	ErrDBInternal        = errors.New("internal DB error during operation")
	ErrNotFound          = errors.New("requested note is not found")
	ErrUserAlreadyExists = errors.New("note already exists")
	ErrUserNotFound      = errors.New("requested note is not found")
)

type task struct {
	db *sqldb.DB
}

func NewTask(db *sqldb.DB) *task {
	return &task{db}
}

func (s *task) CreateTask(ctx context.Context, title string, username string, text string) (uuid.UUID, error) {
	retID, err := s.db.CreateTask(&db.Task{
		ID:        uuid.New(),
		Title:     title,
		User:      username,
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	switch {
	case err != nil:
		return uuid.Nil, err
	default:
		return retID, nil
	}
}

func (s *task) GetAllTasksFromUser(ctx context.Context, username string) ([]db.Task, error) {
	notes, err := s.db.GetAllTasksFromUser(ctx, username)

	if err != nil {
		return nil, ErrDBInternal
	}
	return notes, nil
}

func (s *task) DeleteTask(ctx context.Context, reqID uuid.UUID) (uuid.UUID, error) {
	id, err := s.db.DeleteTask(ctx, reqID)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return uuid.Nil, ErrNotFound
	case err != nil:
		return uuid.Nil, ErrDBInternal
	default:
		return id, nil
	}
}

func (s *task) UpdateTask(ctx context.Context, reqID uuid.UUID, title string, text string, isTextValid bool) (uuid.UUID, error) {
	id, err := s.db.UpdateTask(&db.Task{
		ID:        reqID,
		Title:     title,
		Text:      text,
		UpdatedAt: time.Now(),
	})

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return uuid.Nil, ErrNotFound
	case err != nil:
		return uuid.Nil, ErrDBInternal
	default:
		return id, nil
	}
}

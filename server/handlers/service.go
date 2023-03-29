package handlers

import (
	"context"

	"tasks/db"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type TaskService interface {
	CreateTask(ctx context.Context, title string, username string, text string) (uuid.UUID, error)
	GetAllTasksFromUser(ctx context.Context, username string) ([]db.Task, error)
	DeleteTask(ctx context.Context, id uuid.UUID) (uuid.UUID, error)
	UpdateTask(ctx context.Context, reqID uuid.UUID, title string, text string, isTextEmpty bool) (uuid.UUID, error)
	RegisterUser(ctx context.Context, args *db.User) (string, error)
	GetUser(ctx context.Context, username string) (db.User, error)
}

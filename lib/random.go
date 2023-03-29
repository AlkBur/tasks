package lib

import (
	"math/rand"
	"strings"
	"time"

	"tasks/db"

	"github.com/google/uuid"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewRandomString(size int) string {
	var sb strings.Builder
	k := len(chars)

	for i := 0; i < size; i++ {
		c := chars[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func NewRandomDBNote(id uuid.UUID) *db.Task {
	note := db.Task{
		ID:        id,
		Title:     NewRandomString(15),
		User:      NewRandomString(10),
		Text:      NewRandomString(60),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &note
}

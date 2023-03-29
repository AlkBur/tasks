package db

import (
	"errors"
	"time"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

var (
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrTaskNotFound      = errors.New("requested task is not found")
	taskBucket           = []byte("task")
)

type Task struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title" validate:"required,min=4"`
	User      string    `json:"user" validate:"required"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (db *DB) CreateTask(task *Task) (uuid.UUID, error) {
	err := db.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(userBucket)
		if err != nil {
			return err
		}
		id := []byte(task.ID.String())
		if !checkDuplicate(bucket.Cursor(), id) {
			data, _ := json.Marshal(task)
			err = bucket.Put(id, data)
		} else {
			return ErrUserAlreadyExists
		}
		return err
	})
	return task.ID, err
}

func (db *DB) GetTask(id string) (*Task, error) {
	task := &Task{}
	if err := db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(taskBucket)
		if bucket == nil {
			return ErrUserNotFound
		}
		b := bucket.Get([]byte(id))
		if b == nil {
			return ErrUserNotFound
		}
		err := json.Unmarshal(b, task)
		return err
	}); err != nil {
		return nil, err
	}
	return task, nil
}

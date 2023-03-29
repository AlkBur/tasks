package db

import (
	"errors"

	bolt "go.etcd.io/bbolt"
)

var (
	ErrUserAlreadyExists = errors.New("note already exists")
	ErrUserNotFound      = errors.New("requested note is not found")
	userBucket           = []byte("user")
)

type User struct {
	Username string `json:"username" validate:"required,min=5,max=30,alphanum"`
	Password string `json:"password" validate:"required,min=5"`
	Email    string `json:"email" validate:"email,required"`
}

func (db *DB) CreateUser(user *User) error {
	err := db.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(userBucket)
		if err != nil {
			return err
		}
		id := []byte(user.Username)
		if !checkDuplicate(bucket.Cursor(), id) {
			data, _ := json.Marshal(user)
			err = bucket.Put(id, data)
		} else {
			return ErrUserAlreadyExists
		}
		return err
	})
	return err
}

func (db *DB) GetUser(name string) (*User, error) {
	user := &User{}
	if err := db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(userBucket)
		if bucket == nil {
			return ErrUserNotFound
		}
		b := bucket.Get([]byte(name))
		if b == nil {
			return ErrUserNotFound
		}
		err := json.Unmarshal(b, user)
		return err
	}); err != nil {
		return nil, err
	}
	return user, nil
}

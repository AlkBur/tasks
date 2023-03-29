package db

import (
	"bytes"
	"errors"

	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	bolt "go.etcd.io/bbolt"
)

var (
	json      = jsoniter.ConfigCompatibleWithStandardLibrary
	ErrNoRows = errors.New("no rows in result set")
)

type DB struct {
	logger *zerolog.Logger
	db     *bolt.DB
}

func NewSQL(filepath string, l *zerolog.Logger) (*DB, error) {
	sqlDB := &DB{
		logger: l,
	}
	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		return nil, err
	}
	sqlDB.db = db
	l.Info().Msg("db connection is successful.")
	return sqlDB, nil

}

func (sql *DB) Close() {
	sql.db.Close()
}

func checkDuplicate(c *bolt.Cursor, id []byte) bool {
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if bytes.Equal(id, v) {
			return true
		}
	}

	return false
}

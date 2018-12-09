package model

import (
	"github.com/boltdb/bolt"
	"time"
)

type Model struct {
	DB *bolt.DB
}

func NewModel() *Model {
	db, err := bolt.Open("config/BlogReaper.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil
	}
	return &Model{db}
}

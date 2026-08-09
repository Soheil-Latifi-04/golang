package main

import (
	"encoding/json"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type storageRow struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

type PostgresStorage[T any] struct {
	db  *gorm.DB
	key string
}

func NewPostgresStorage[T any](dsn, key string) (*PostgresStorage[T], error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&storageRow{}); err != nil {
		return nil, err
	}

	return &PostgresStorage[T]{db: db, key: key}, nil
}

func (s *PostgresStorage[T]) Save(data T) error {
	fileData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	row := storageRow{Key: s.key, Value: string(fileData)}

	return s.db.Save(&row).Error // upsert on primary key
}

func (s *PostgresStorage[T]) Load(data *T) error {
	var row storageRow
	err := s.db.First(&row, "key = ?", s.key).Error
	if err == gorm.ErrRecordNotFound {
		return nil // nothing saved yet, leave data as zero value
	}
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(row.Value), data)
}

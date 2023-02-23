package database

import (
	"time"

	"gorm.io/gorm"
)

type Database interface {
	// Instance is used to get primary database instance.
	Instance() *gorm.DB
	// Close is used to close database connection.
	Close() error
	// SetMaxIdleConns is used to configure maximum idle connections.
	SetMaxIdleConns(n int) error
	// SetMaxOpenConns is used to configure maximum openned connections.
	SetMaxOpenConns(n int) error
	// SetConnMaxLifetime is used to configure maximum openned connection lifetime.
	SetConnMaxLifetime(d time.Duration) error
}

// Model provides base fields for all database models (like gorm.Model).
type Model struct {
	CreatedAt time.Time      `json:"-" gorm:"index"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" swaggerignore:"true"`
}

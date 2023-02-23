package database

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgreSQLConfig struct {
	User     string
	Password string
	Host     string
	Database string
}

type PostgreSQL struct {
	DB *gorm.DB
}

// Check if implements the interface.
var _ Database = (*PostgreSQL)(nil)

// NewPostgreSQL is used to create new instance of PostgreSQL.
func NewPostgreSQL(cfg *PostgreSQLConfig) (*PostgreSQL, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s",
		cfg.User, cfg.Password, cfg.Database, cfg.Host,
	)

	// Connect to the database by DSN
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
	}

	// create UUID extension.
	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid-ossp extension: %w", err)
	}

	return &PostgreSQL{DB: db}, nil
}

func (p *PostgreSQL) Instance() *gorm.DB {
	return p.DB
}

func (p *PostgreSQL) Close() error {
	if p.DB == nil {
		return errors.New("db connection is already closed")
	}
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (p *PostgreSQL) SetMaxIdleConns(n int) error {
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(n)
	return nil
}

func (p *PostgreSQL) SetMaxOpenConns(n int) error {
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(n)
	return nil
}

func (p *PostgreSQL) SetConnMaxLifetime(d time.Duration) error {
	db, err := p.DB.DB()
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(d)
	return nil
}

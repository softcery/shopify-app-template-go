package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/softcery/shopify-app-template-go/internal/entity"
	"github.com/softcery/shopify-app-template-go/internal/service"
	"github.com/softcery/shopify-app-template-go/pkg/database"
	"gorm.io/gorm"
)

type sessionStorage struct {
	database.Database
}

var _ service.SessionStorage = (*sessionStorage)(nil)

func NewSessionStorage(db database.Database) *sessionStorage {
	return &sessionStorage{db}
}

func (s *sessionStorage) Get(ctx context.Context, sessionID string) (*entity.Session, error) {
	stmt := s.Instance().
		Where(&entity.Session{SessionID: sessionID})

	var session entity.Session
	err := stmt.First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (s *sessionStorage) Create(ctx context.Context, session *entity.Session) (*entity.Session, error) {
	err := s.Instance().Create(session).Error
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *sessionStorage) Delete(ctx context.Context, sessionID string) error {
	err := s.Instance().Delete(&entity.Session{}, "session_id = ?", sessionID).Error
	if err != nil {
		return err
	}
	return nil
}

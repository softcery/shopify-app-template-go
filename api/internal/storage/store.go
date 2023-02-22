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

type storeStorage struct {
	database.Database
}

var _ service.StoreStorage = (*storeStorage)(nil)

func NewStoreStorage(db database.Database) *storeStorage {
	return &storeStorage{db}
}

func (s *storeStorage) Get(ctx context.Context, storeName string) (*entity.Store, error) {
	stmt := s.Instance().
		Where(&entity.Store{Name: storeName})

	var store entity.Store
	err := stmt.First(&store).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get store: %w", err)
	}

	return &store, nil
}

func (s *storeStorage) Update(ctx context.Context, store *entity.Store) (*entity.Store, error) {
	err := s.Instance().
		Where(&entity.Store{Name: store.Name}).
		Updates(store).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to update store: %w", err)
	}

	var updatedStore entity.Store
	err = s.Instance().
		Where(&entity.Store{Name: store.Name}).
		First(&updatedStore).
		Error
	if err != nil {
		return nil, fmt.Errorf("failed to get updated store: %w", err)
	}

	return &updatedStore, nil
}

func (s *storeStorage) Create(ctx context.Context, store *entity.Store) (*entity.Store, error) {
	err := s.Instance().Create(store).Error
	if err != nil {
		return nil, err
	}
	return store, nil
}

func (s *storeStorage) Delete(ctx context.Context, storeName string) error {
	err := s.Instance().Delete(&entity.Store{}, "name = ?", storeName).Error
	if err != nil {
		return err
	}
	return nil
}

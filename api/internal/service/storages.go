package service

import (
	"context"

	"github.com/softcery/shopify-app-template-go/internal/entity"
)

// Storages contains all available storages.
type Storages struct {
	Store StoreStorage
}

type StoreStorage interface {
	// Get is used to retrieve store from storage by its name.
	Get(ctx context.Context, storeName string) (*entity.Store, error)
	// Create is used to create new store.
	Create(ctx context.Context, store *entity.Store) (*entity.Store, error)
	// Update is used to update store.
	Update(ctx context.Context, store *entity.Store) (*entity.Store, error)
	// Delete is used to delete store.
	Delete(ctx context.Context, storeName string) error
}

type SessionStorage interface {
	// Get is used to retrieve session from storage by its ID.
	Get(ctx context.Context, sessionID string) (*entity.Session, error)
	// Create is used to create new session.
	Create(ctx context.Context, session *entity.Session) (*entity.Session, error)
	// Delete is used to delete session.
	Delete(ctx context.Context, sessionID string) error
}

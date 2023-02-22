package entity

import (
	"github.com/softcery/shopify-app-template-go/pkg/database"
)

// Store model represents model of platform store.
type Store struct {
	database.Model
	ID   string `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name string `gorm:"index"`

	// Shopify
	Nonce       string
	AccessToken string
	Installed   bool
}

type Session struct {
	SessionID string
	StoreID   string
}

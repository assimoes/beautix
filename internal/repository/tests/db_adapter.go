package tests

import (
	"github.com/assimoes/beautix/internal/infrastructure/database"
	"gorm.io/gorm"
)

// DBAdapter adapts a *gorm.DB to look like a *database.DB
// This allows us to pass transaction objects to repository constructors
type DBAdapter struct {
	*gorm.DB
}

// NewDBAdapter creates a new adapter around a gorm.DB instance
func NewDBAdapter(tx *gorm.DB) *database.DB {
	// Cast to database.DB since they both embed *gorm.DB
	// This is safe because database.DB simply wraps gorm.DB
	return &database.DB{DB: tx}
}

// MockDB implementation removed - using integration tests with real database instead

package datastores

import (
	"fmt"

	"github.com/smolse/fluffy-pancake/internal/config"
	"github.com/smolse/fluffy-pancake/internal/models"
)

// DataStore is an interface for a data store that can store and retrieve risk data.
type DataStore interface {
	// Connection methods
	Connect() error
	Close() error

	// Query methods
	GetRisk(id models.RiskId) (models.RiskAttributes, error)
	CreateRisk(id models.RiskId, attrs models.RiskAttributes) error
	ListRisks() ([]models.Risk, error)
}

// NewDataStore is a factory function that creates a new DataStore instance based on the given data store type.
func NewDataStore(cfg *config.DataStoreConfig) (DataStore, error) {
	switch cfg.Type {
	case "syncmap":
		return NewSyncMapDataStore(), nil
	default:
		return nil, fmt.Errorf("unsupported data store type: %s", cfg.Type)
	}
}

package datastores

import (
	"fmt"
	"sync"

	"github.com/smolse/fluffy-pancake/internal/models"
)

// SyncMapDataStore is a simple in-memory data store that uses a sync.Map to store KV pairs. Safe for concurrent use.
type SyncMapDataStore struct {
	syncMap sync.Map
}

// NewSyncMapDataStore creates a new SyncMapDataStore instance.
func NewSyncMapDataStore() *SyncMapDataStore {
	return &SyncMapDataStore{}
}

// Connect does nothing for the local SyncMapDataStore.
func (s *SyncMapDataStore) Connect() error {
	return nil
}

// Close does nothing for the local SyncMapDataStore.
func (s *SyncMapDataStore) Close() error {
	return nil
}

// GetRisk retrieves the attribute values for the given risk ID key.
func (s *SyncMapDataStore) GetRisk(id models.RiskId) (models.RiskAttributes, error) {
	value, ok := s.syncMap.Load(id)
	if !ok {
		return models.RiskAttributes{}, fmt.Errorf("risk %s was not found", id.String())
	}

	return value.(models.RiskAttributes), nil
}

// CreateRisk stores the given risk in the data store.
func (s *SyncMapDataStore) CreateRisk(id models.RiskId, attrs models.RiskAttributes) error {
	s.syncMap.Store(id, attrs)
	return nil
}

// List returns a list of all risks in the data store.
func (s *SyncMapDataStore) ListRisks() ([]models.Risk, error) {
	risks := []models.Risk{}
	s.syncMap.Range(func(key, value interface{}) bool {
		risks = append(risks, models.Risk{Id: key.(models.RiskId), RiskAttributes: value.(models.RiskAttributes)})
		return true
	})
	return risks, nil
}

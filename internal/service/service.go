package service

import (
	"github.com/smolse/fluffy-pancake/internal/datastores"
	"github.com/smolse/fluffy-pancake/internal/models"
)

// RiskService is a service that provides risk-related operations.
type RiskService struct {
	dataStore datastores.DataStore
}

// NewRiskService is a constructor function that creates a new RiskService instance.
func NewRiskService(dataStore datastores.DataStore) *RiskService {
	return &RiskService{
		dataStore: dataStore,
	}
}

// GetRisk gets a risk by ID.
func (s *RiskService) GetRisk(id models.RiskId) (models.Risk, error) {
	riskAttrs, err := s.dataStore.GetRisk(id)
	if err != nil {
		return models.Risk{}, err
	}

	return models.Risk{Id: id, RiskAttributes: riskAttrs}, nil
}

// CreateRisk creates a risk.
func (s *RiskService) CreateRisk(id models.RiskId, attrs models.RiskAttributes) error {
	return s.dataStore.CreateRisk(id, attrs)
}

// ListRisks lists all risks.
func (s *RiskService) ListRisks() ([]models.Risk, error) {
	return s.dataStore.ListRisks()
}

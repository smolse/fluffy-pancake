package service

// this is a test module for service.go

import (
	"testing"

	"github.com/smolse/fluffy-pancake/internal/datastores"
	"github.com/smolse/fluffy-pancake/internal/models"
)

func TestGetRisk_NotFound(t *testing.T) {
	// Create a new risk service backed by an empty SyncMapDataStore
	dataStore := datastores.NewSyncMapDataStore()
	riskService := NewRiskService(dataStore)

	// Get a risk that doesn't exist
	_, err := riskService.GetRisk(models.RiskId(models.NewRiskId()))

	// Check that the error is propagated as expected
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetRisk_Found(t *testing.T) {
	// Create a new risk service backed by an empty SyncMapDataStore
	dataStore := datastores.NewSyncMapDataStore()
	riskService := NewRiskService(dataStore)

	// Create a new risk
	id := models.RiskId(models.NewRiskId())
	attrs := models.RiskAttributes{State: "open"}
	dataStore.CreateRisk(id, attrs)

	// Get the created risk
	risk, err := riskService.GetRisk(id)

	// Check that the risk is returned as expected
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if risk.Id != id {
		t.Errorf("unexpected risk ID: got %s, want %s", risk.Id, id)
	}
	if risk.RiskAttributes != attrs {
		t.Errorf("unexpected risk attributes: got %v, want %v", risk.RiskAttributes, attrs)
	}
}

func TestCreateRisk(t *testing.T) {
	// Create a new risk service backed by an empty SyncMapDataStore
	dataStore := datastores.NewSyncMapDataStore()
	riskService := NewRiskService(dataStore)

	// Create a new risk
	id := models.RiskId(models.NewRiskId())
	attrs := models.RiskAttributes{State: "open"}
	err := riskService.CreateRisk(id, attrs)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// Check that the risk is persisted in the data store as expected
	risk, err := dataStore.GetRisk(id)
	if err != nil {
		t.Errorf("failed to get risk from data store: %s", err)
	}
	if risk != attrs {
		t.Errorf("unexpected risk attributes: got %v, want %v", risk, attrs)
	}
}

func TestListRisks_Empty(t *testing.T) {
	// Create a new risk service backed by an empty SyncMapDataStore
	dataStore := datastores.NewSyncMapDataStore()
	riskService := NewRiskService(dataStore)

	// List all risks
	risks, err := riskService.ListRisks()

	// Check that no risks are returned
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(risks) != 0 {
		t.Errorf("unexpected number of risks: got %d, want 0", len(risks))
	}
}

func TestListRisks_NonEmpty(t *testing.T) {
	// Create a new risk service backed by an empty SyncMapDataStore
	dataStore := datastores.NewSyncMapDataStore()
	riskService := NewRiskService(dataStore)

	// Create some risks
	ids := []models.RiskId{
		models.RiskId(models.NewRiskId()),
		models.RiskId(models.NewRiskId()),
		models.RiskId(models.NewRiskId()),
	}
	attrs := []models.RiskAttributes{
		{State: "open"},
		{State: "closed"},
		{State: "accepted"},
	}
	for i := range ids {
		dataStore.CreateRisk(ids[i], attrs[i])
	}

	// List all risks
	risks, err := riskService.ListRisks()

	// Check that the risks are returned as expected
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(risks) != len(ids) {
		t.Errorf("unexpected number of risks: got %d, want %d", len(risks), len(ids))
	}
}

package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// ValidRiskStates is a list of valid risk states.
var ValidRiskStates = []string{"open", "closed", "accepted", "investigating"}

// RiskId is a type that represents the ID of a risk.
type RiskId uuid.UUID

// RiskAttributes is a type that represents the attributes of a risk.
type RiskAttributes struct {
	State       string `json:"state" binding:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Risk is a type that represents a risk.
type Risk struct {
	Id RiskId `json:"id"`
	RiskAttributes
}

// NewRiskId creates a new RiskId randomly.
func NewRiskId() RiskId {
	return RiskId(uuid.New())
}

// NewRiskIdFromString creates a new RiskId from a string value.
func NewRiskIdFromString(id string) (RiskId, error) {
	u, err := uuid.Parse(id)
	return RiskId(u), err
}

// String returns the string representation of a RiskId type value.
func (rid RiskId) String() string {
	return uuid.UUID(rid).String()
}

// MarshalJSON marshals a RiskId type value to JSON.
func (rid RiskId) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(rid).String())
}

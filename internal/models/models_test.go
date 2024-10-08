package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewRiskId(t *testing.T) {
	rid := NewRiskId()
	if _, err := uuid.Parse(rid.String()); err != nil {
		t.Errorf("new risk id %q is not a valid UUID", rid.String())
	}
}

func TestNewRiskIdFromString(t *testing.T) {
	ridStr := "01234567-89ab-cdef-fedc-ba9876543210"
	_, err := NewRiskIdFromString(ridStr)
	if err != nil {
		t.Errorf("failed to create risk id from a valid UUID string: %s", err)
	}
}

func TestRiskId_MarshalJSON(t *testing.T) {
	u, _ := uuid.FromBytes([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10})
	want := "\"01234567-89ab-cdef-fedc-ba9876543210\""
	got, _ := RiskId(u).MarshalJSON()
	if string(got) != want {
		t.Errorf("failed to properly marshal RiskId to JSON; got %q, want %q", got, want)
	}
}

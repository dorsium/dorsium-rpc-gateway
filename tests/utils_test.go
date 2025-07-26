package tests

import (
	"testing"

	"github.com/dorsium/dorsium-rpc-gateway/pkg/utils"
)

func TestReverse(t *testing.T) {
	expected := "olleh"
	if utils.Reverse("hello") != expected {
		t.Errorf("expected %s", expected)
	}
}

func TestIsValidAddress(t *testing.T) {
	valid := []string{
		"0x1234567890abcdef1234567890abcdef12345678",
		"cosmos1qqpcymswlj9teu9gf28elksg5926y4v5d9dx7a",
	}
	for _, addr := range valid {
		if !utils.IsValidAddress(addr) {
			t.Errorf("expected valid address: %s", addr)
		}
	}
	invalid := []string{"", "0x123", "cosmos1invalid"}
	for _, addr := range invalid {
		if utils.IsValidAddress(addr) {
			t.Errorf("expected invalid address: %s", addr)
		}
	}
}

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

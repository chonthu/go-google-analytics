package gadata

import (
	"testing"
)

// TestImportConfig verifies oauth configs
// can be successfully loaded from JSON file
func TestImportConfig(t *testing.T) {
	authFile := "client_secret.json"
	conf, err := ImportConfig(authFile)
	if conf != nil {
		t.Errorf("No data read from %s", authFile)
	}
	if err != nil {
		t.Errorf("Error reading file %s", authFile)
	}
}

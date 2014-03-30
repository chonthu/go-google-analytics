package gadata

import (
	"testing"
)

// TestImportConfig verifies oauth configs
// can be successfully loaded from JSON file
func TestImportConfig(t *testing.T) {
	authFile := "client_secret.json"
	conf, err := ImportConfig(authFile)
	if conf == nil {
		t.Errorf("No data read from %s", authFile)
	}
	if err != nil {
		t.Errorf("Error reading file %s", authFile)
	}
}

// TestDBGetTokens validates token retrieval from
// default local store (.dat file).
func TestDBGetTokens(t *testing.T) {
	tokens, err := DBGetTokens()
	if err != nil {
		t.Errorf("Error retrieving tokens from DB")
	} else if tokens == nil {
		t.Errorf("No token data imported from DB")
	}

}

// TestInitConnection initialises gadata object
// and writes out tokens to local dat file if
// one doesn't exist.
func TestInitConnection(t *testing.T) {
	InitConnection()
}

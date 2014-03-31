package gadata

import (
	"strings"
	"testing"
)

// OauthData object
var g *OauthData = new(OauthData)

// TestImportConfig verifies oauth configs
// can be successfully loaded from JSON file
func TestImportConfig(t *testing.T) {
	authFile := "client_secret.json"
	err := g.ImportConfig(authFile)
	if err != nil {
		t.Errorf("Error reading file %s", authFile)
	}

	if g.config == nil {
		t.Errorf("No data read from %s", authFile)
	}
}

// TestCheckTokens checks if tokens already exist locally
func TestCheckTokens(t *testing.T) {
	err := g.checkTokens()
	if err != nil {
		t.Errorf("Error checking tokens in local flat file")
	}
}

// TestTokenStore checks the tokens storing functionality
func TestTokenStore(t *testing.T) {
	testData := AccessData{"test", 1, "test", "test"}
	err := g.tokenStore(&testData)
	if err != nil {
		t.Errorf("Error storing data into flat file.")
	}

	g.checkTokens()
	if testData != *g.tokens {
		t.Errorf("Stored and retrieved tokens aren't the same.")
	}
}

// TestProcessTokenResponse validates oauth response handling
// functionality.
func TestProcessTokenResponse(t *testing.T) {
	contents := `{"error": "Test error", "error_description": "Test error description"}`
	data, err := g.ProcessTokenResponse([]byte(contents))
	if err == nil {
		t.Errorf("Error processing test token response")
	}

	contents = `{  "access_token" : "ya29.1.AADtN_VcH4-tJjvsYl77Q93tVZIKkCcqfvjO5mcsWg6gKyEIqmrSnH2dH6B49g",  "token_type" : "Bearer",  "expires_in" : 3600,  "refresh_token" : "1/3rUMCQPYJ9BcoXOFyjIOTTBWFuclTo1cQV_LZmKAp24"}`
	data, err = g.ProcessTokenResponse([]byte(contents))
	if err != nil && data == nil {
		t.Errorf("Error processing test token response.")
	}
	if !strings.Contains(contents, data.AccessToken) {
		t.Errorf("Input and instanced access tokens don't match")
	}
}

// TestInitConnection initialises OauthData object
// and writes out tokens to local dat file if
// one doesn't exist.
// func TestInitConnection(t *testing.T) {
// 	g.clearStore()
// 	g.InitConnection()
// }

// TestRefreshToken validates token refresh functionality
func TestRefreshToken(t *testing.T) {
	err := g.refreshToken()
	if err != nil {
		t.Errorf("Error refreshing token")
	}
}

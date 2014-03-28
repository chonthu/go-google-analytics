package gadata

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type AuthInfo struct {
	ClientID     string   `json:client_id`
	ClienSecret  string   `json:client_secret`
	AuthURI      string   `json:auth_uri`
	TokenURI     string   `json:token_uri`
	RedirectURIs []string `json:redirect_uris`
}

// ImportConfig imports client secret information from the JSON obtained
// from Google Developer Console (API console).
func ImportConfig(filename string) (conf *AuthInfo, err error) {
	config := new(AuthInfo)
	configFile, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Error opening config file: %s\n", filename)
	}
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Fatalf("Error encoding config struct", err.Error())
	}

	fmt.Println(config.ClienSecret)
	return
}

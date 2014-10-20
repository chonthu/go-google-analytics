/*
	OAuth handshake functionality
*/

package gadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// Base constants
const (
	Scope          string = "https://www.googleapis.com/auth/analytics.readonly"
	ReturnURI      string = "localhost:9000"
	LocalStoreFile string = "localStore.dat"
)

// AuthInfo is the outer structure of client_secret.json
type AuthInfo struct {
	Data AuthData `json:"installed"`
}

// AuthData is the primary struct of JSON config file
type AuthData struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
	RedirectURIs []string `json:"redirect_uris"`
}

// AuthCode is the authorisation code POST request struct
// for exchanging for access and refresh tokens
type AuthCode struct {
	Code         string `json:"code"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

// RefreshData is the token refresh response data structure
type RefreshData struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (a *AuthCode) ToURLValues() (out url.Values) {
	out = url.Values{"code": {a.Code},
		"client_id":     {a.ClientID},
		"client_secret": {a.ClientSecret},
		"redirect_uri":  {a.RedirectURI},
		"grant_type":    {a.GrantType},
	}
	return
}

// AccessData contains the access and refresh token information
// received after a successful retrieval from Google or local stores
type AccessData struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

// GError is the error strucure of oauth2 server response when
// requesting the token set
type GError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// GAData is the working struct of the library
type OauthData struct {
	config   *AuthInfo
	tokens   *AccessData
	JSONfile string
}

// ImportConfig imports client secret information from the JSON obtained
// from Google Developer Console (API console).
func (d *OauthData) ImportConfig(filename string) (err error) {
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error opening config file: %s\n", filename)
	}

	err = json.Unmarshal(configFile, &d.config)
	if err != nil {
		log.Fatalf("Error unmarshalling config file: %s\n", filename)
	}
	log.Println("Imported config file.")
	d.JSONfile = filename
	return
}

func (d *OauthData) RegisterClient() (err error) {
	authUrl := fmt.Sprintf("%s?scope=%s&redirect_uri=http://%s&response_type=code&client_id=%s", d.config.Data.AuthURI, Scope, ReturnURI, d.config.Data.ClientID)

	// Create new channel and spin off callback listener
	ch := make(chan string)
	go WebCallback(ch)

	// Open client authentication URL in default system browser
	err = BrowserOpen(authUrl)
	checkError(err)

	// Listen for callback value or die with a timeout after 60 sec
	var newAccessCode string
	select {
	case newAccessCode = <-ch:
		log.Println("Received new access code")
	case <-time.After(60 * 1e9):
		log.Fatalln("Didn't receive auth token response... giving up.")
		return
	}

	// Retrieve access and refresh tokens
	PostData := AuthCode{newAccessCode, d.config.Data.ClientID, d.config.Data.ClientSecret, "http://" + ReturnURI, "authorization_code"}

	// retrieve new tokens
	req, err := http.PostForm(d.config.Data.TokenURI, PostData.ToURLValues())
	checkError(err)
	defer req.Body.Close()
	contents, err := ioutil.ReadAll(req.Body)
	checkError(err)
	data, err := d.ProcessTokenResponse(contents)
	checkError(err)
	d.tokens = data

	return
}

func (d *OauthData) InitConnection() (err error) {
	err = d.ImportConfig("client_secret.json")
	checkError(err)

	err = d.checkTokens()
	if err != nil {
		d.RegisterClient()
	}
	return
}

func (d *OauthData) ProcessTokenResponse(contents []byte) (data *AccessData, err error) {
	// trim all the line breaks
	reg, err := regexp.Compile("\n")
	checkError(err)
	out := reg.ReplaceAllString(string(contents), "")

	// check if oauth server returned an error
	if strings.Contains(out, "error_description") {
		var tempError GError
		json.Unmarshal([]byte(out), &tempError)
		err = errors.New(fmt.Sprintf("Encountered error %s : %s", tempError.Error, tempError.ErrorDescription))
	} else {
		err := json.Unmarshal([]byte(out), &data)
		checkError(err)
		d.tokenStore(data)
	}
	return
}

// Process token refresh
func (d *OauthData) ProcessRefreshResponse(contents []byte) (data *RefreshData, err error) {
	// trim all the line breaks
	reg, err := regexp.Compile("\n")
	checkError(err)
	out := reg.ReplaceAllString(string(contents), "")

	// check if oauth server returned an error
	if strings.Contains(out, "error_description") {
		var tempError GError
		json.Unmarshal([]byte(out), &tempError)
		err = errors.New(fmt.Sprintf("Encountered error %s : %s", tempError.Error, tempError.ErrorDescription))
	} else {
		err := json.Unmarshal([]byte(out), &data)
		checkError(err)
	}
	return
}

// tokenStore stores token data in a local flat file in JSON format
func (d *OauthData) tokenStore(data *AccessData) (err error) {
	out, err := json.Marshal(data)
	checkError(err)
	err = ioutil.WriteFile(LocalStoreFile, out, 0644)
	checkError(err)
	return
}

// checkTokens attempts to retrieve token data from the local flat file
func (d *OauthData) checkTokens() (err error) {
	data, err := ioutil.ReadFile(LocalStoreFile)
	if data == nil {
		return errors.New("Local store file is empty.")
	}
	err = json.Unmarshal(data, &d.tokens)
	// Case default token
	if d.tokens.TokenType == "test" {
		return errors.New("Invalid localstore file data.")
	}
	return
}

// RefreshToken refreshes access token in case of expiry
func (d *OauthData) refreshToken() (err error) {

	// retrieve new tokens
	req, err := http.PostForm(d.config.Data.TokenURI, url.Values{"refresh_token": {d.tokens.RefreshToken}, "client_id": {d.config.Data.ClientID}, "client_secret": {d.config.Data.ClientSecret}, "grant_type": {"refresh_token"}})
	checkError(err)
	defer req.Body.Close()
	contents, err := ioutil.ReadAll(req.Body)
	checkError(err)
	data, err := d.ProcessRefreshResponse(contents)
	checkError(err)
	d.tokens.AccessToken = data.AccessToken
	d.tokenStore(d.tokens)
	log.Println("Access token refreshed")
	return
}

// clearStore deletes the local store file
func (d *OauthData) clearStore() (err error) {
	err = os.Remove(LocalStoreFile)
	log.Println("Cleared local token store file")
	return
}

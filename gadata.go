/*
	GAData
*/

package gadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// Base constants
const (
	Scope     string = "https://www.googleapis.com/auth/analytics.readonly"
	ReturnURI string = "localhost:9000"
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

// Check for normal errors
func checkError(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// BrowserOpen opens a URL is the OS' default web browser
func BrowserOpen(url string) error {
	return exec.Command("open", url).Run()
}

// WebCallback listens on a predefined port for a oauth response
// sends back via channel once it receives a response and shuts down.
func WebCallback(ch chan string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.RequestURI, "/?") {
			block := strings.SplitAfter(r.RequestURI, "/?")[1]
			if !strings.Contains(block, "code=") {
				ch <- block
			} else {
				ch <- strings.SplitAfter(block, "code=")[1]
			}
		}
		return
	})

	log.Fatal(http.ListenAndServe(ReturnURI, nil))
}

// ImportConfig imports client secret information from the JSON obtained
// from Google Developer Console (API console).
func ImportConfig(filename string) (config *AuthInfo, err error) {
	configFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error opening config file: %s\n", filename)
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Error unmarshalling config file: %s\n", filename)
	}
	log.Println("Imported config file.")
	return
}

func InitConnection() {
	config, err := ImportConfig("client_secret.json")
	if err != nil {
		log.Fatalln(err.Error())
	}
	authUrl := fmt.Sprintf("%s?scope=%s&redirect_uri=http://%s&response_type=code&client_id=%s", config.Data.AuthURI, Scope, ReturnURI, config.Data.ClientID)

	// Create new channel and spin off callback listener
	ch := make(chan string)
	go WebCallback(ch)

	// Open client authentication URL in default system browser
	err = BrowserOpen(authUrl)
	checkError(err)

	// Listen for callback value or die with a timeout after 30 sec
	var newAccessCode string
	select {
	case newAccessCode = <-ch:
		log.Println(newAccessCode)
	case <-time.After(30 * 1e9):
		log.Fatalln("Didn't receive response... giving up.")
		return
	}

	// Retrieve access and refresh tokens
	PostData := AuthCode{newAccessCode, config.Data.ClientID, config.Data.ClientSecret, "http://" + ReturnURI, "authorization_code"}

	// retrieve new tokens
	req, err := http.PostForm(config.Data.TokenURI, PostData.ToURLValues())
	checkError(err)
	defer req.Body.Close()
	contents, err := ioutil.ReadAll(req.Body)
	checkError(err)

	// trim all the line breaks
	reg, err := regexp.Compile("\n")
	checkError(err)
	out := reg.ReplaceAllString(string(contents), "")
	log.Println(out)
	if strings.Contains(out, "error_description") {
		var data GError
		json.Unmarshal([]byte(out), &data)
		log.Println(data.ErrorDescription)
	} else {
		var data AccessData
		json.Unmarshal([]byte(out), &data)
		log.Println(data)
	}
}

func DBGetTokens() (tokens map[string]string, err error) {
	return
}

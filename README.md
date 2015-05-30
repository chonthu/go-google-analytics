##Google analytics Data pull

Lightweight Golang library for pulling Google Analytics API data.
Built for use with Core Reporting API (v3):

https://developers.google.com/analytics/devguides/reporting/core/v3/reference

### Authentication
In order to authenticate this library for use with your Google Analytics account, an oauth2 token needs to be generated. For a new project login to [Google Developers Console](https://console.developers.google.com) and Create Project. 

Add Analytics API to list of APIs, create a new "Installed" App Client ID and download it in JSON format.

Place the client_secret.json in the root of your application. Ps. you have to renaming it from the crazy name to just "client_secret.json"

### Usage

See Examples [here](https://github.com/chonthu/go-google-analytics/tree/master/utils)

### Testing
Unit tests are included with this library, use `go test ./...` to run through the set provided.

** This doesnt really work yet, but working on it  **

### Changelog
#### 1.0.0
- cleaner naming
- clearner working examples
#### 0.1.1:
- Implemented batch processing
- New request period segmentation functionality
- Cleaner error reporting and resolution suggestions
#### 0.1.0:
- Initial release
GAData
===================

Lightweight library for pulling Google Analytics API data

### Authentication
In order to authenticate this library for use with your Google Analytics account, an oauth2 token needs to be generated. For a new project login to [Google Developers Console](https://console.developers.google.com) and Create Project. Add Analytics API to list of APIs,  create a new Client ID and download it in JSON format.

### Testing
Unit tests are included with this library, use `go test ./...` to run through the set provided. 
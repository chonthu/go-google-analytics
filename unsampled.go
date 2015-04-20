package gadata

// Simple interface to generating unsampled reports

const (
	BaseEndpoint string = "https://www.googleapis.com/analytics/v3/management/accounts/"
	End          string = "unsampledReports"
)

type UnsampledRequest struct {
	base          string  //
	accountID     string  //accountId/
	webproperties string  //webPropertyId
	profiles      string  //profileId/
	end           string  //unsampledReports
	payload       *URData // request payload
}

type URData struct {
	startDate  string
	endDate    string
	metrics    string
	title      string
	dimensions string
	filters    string
	segment    string
}

// expected return
// https://developers.google.com/analytics/devguides/config/mgmt/v3/mgmtReference/management/unsampledReports#methods

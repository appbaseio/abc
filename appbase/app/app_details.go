package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
	"github.com/olekukonko/tablewriter"
)

// Permission represents an app permission object
type Permission struct {
	Description string `json:"description"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type metricsBucket struct {
	// DocCount  int64                  `json:"doc_count"`
	APICalls  map[string]json.Number `json:"apiCalls"`
	DateAsStr string                 `json:"key_as_string"`
}

type metricsBuckets struct {
	Buckets []metricsBucket `json:"buckets"`
}

type metricsOverall struct {
	NumDocs int64 `json:"numDocs"`
	Storage int   `json:"storage"`
}

type metricsBody struct {
	Month   metricsBuckets `json:"month"`
	Overall metricsOverall `json:"overall"`
}

type respBodyPerms struct {
	Body []Permission `json:"body"`
}

type respBodyMetrics struct {
	Body metricsBody `json:"body"`
}

// ShowAppMetrics ...
func ShowAppMetrics(app string) error {
	spinner.StartText("Fetching app metrics")
	defer spinner.Stop()
	// show metrics
	fmt.Println()
	req, err := http.NewRequest("GET", common.AccAPIURL+"/app/"+app+"/metrics", nil)
	if err != nil {
		return err
	}
	resp, err := session.SendRequest(req)
	if err != nil {
		return err
	}
	spinner.Stop()
	// decode
	var res respBodyMetrics
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return err
	}
	// output
	fmt.Printf("Storage(KB): %d\n", common.SizeInKB(res.Body.Overall.Storage))
	fmt.Printf("Records:     %d\n", res.Body.Overall.NumDocs)
	// table
	var callCount int64
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "API Calls"})
	latest := res.Body.Month.Buckets[common.Max(0, len(res.Body.Month.Buckets)-30):len(res.Body.Month.Buckets)]
	for _, bucket := range latest {
		if bucket.APICalls["value"] == "0" {
			continue
		}
		table.Append([]string{
			getHumanDate(bucket.DateAsStr), common.JSONNumberToString(bucket.APICalls["value"]),
		})
		callCount = callCount + common.JSONNumberToInt(bucket.APICalls["value"])
	}
	table.SetFooter([]string{"Total",
		strconv.FormatInt(callCount, 10),
	})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.Render()
	return nil
}

// ShowAppAnalytics fetches analytics for an app
func ShowAppAnalytics(app string, endpoint string) error {
	spinner.StartText("Fetching app analytics")
	defer spinner.Stop()
	// show analytics
	fmt.Println()
	req, err := http.NewRequest("GET", common.AccAPIURL+"/analytics/"+app+"/"+endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := session.SendRequest(req)
	if err != nil {
		return err
	}
	spinner.Stop()

	switch endpoint {
	case "latency":
		ShowLatency(resp.Body)
	case "geoip":
		ShowGeoIP(resp.Body)
	case "overview":
		ShowOverview(resp.Body)
	case "popularresults":
		ShowPopularResults(resp.Body)
	case "popularsearches":
		ShowPopularSearches(resp.Body)
	case "popularfilters":
		ShowPopularFilters(resp.Body)
	case "noresultsearches":
		ShowNoResultSearches(resp.Body)
	}

	return nil
}

// ShowAppPerms ...
func ShowAppPerms(app string) error {
	spinner.StartText("Fetching app credentials")
	fmt.Println()
	permissions, err := GetAppPerms(app)
	if err != nil {
		return err
	}
	spinner.Stop()
	// output
	for _, p := range permissions {
		fmt.Printf("%s%s:%s\n", common.ColonPad(p.Description, 20),
			p.Username, p.Password)
	}
	return nil
}

// GetAppPerms ...
func GetAppPerms(app string) ([]Permission, error) {
	req, err := http.NewRequest("GET", common.AccAPIURL+"/app/"+app+"/permissions", nil)
	if err != nil {
		return nil, err
	}
	resp, err := session.SendRequest(req)
	if err != nil {
		return nil, err
	}
	// decode
	var res respBodyPerms
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func getHumanDate(date string) string {
	t, err := time.Parse("2006-01-02", date[0:10])
	if err != nil {
		return date
	}
	return t.Format("2006-Jan-02")
}

package app

import (
	"encoding/json"
	"fmt"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"os"
	"strconv"
)

// Permission represents an app permission object
type Permission struct {
	Description string `json:"description"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type metricsBucket struct {
	DocCount  int64                  `json:"doc_count"`
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
	var docCount, callCount int64
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Date", "API Calls", "Records"})
	for _, bucket := range res.Body.Month.Buckets {
		table.Append([]string{
			getHumanDate(bucket.DateAsStr), common.JSONNumberToString(bucket.APICalls["value"]),
			strconv.FormatInt(bucket.DocCount, 10),
		})
		docCount = docCount + bucket.DocCount
		callCount = callCount + common.JSONNumberToInt(bucket.APICalls["value"])
	}
	table.SetFooter([]string{"Total",
		strconv.FormatInt(callCount, 10), strconv.FormatInt(docCount, 10),
	})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.Render()
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
	for index := range ".." {
		fmt.Printf("%s%s:%s\n", common.ColonPad(permissions[index].Description, 20),
			permissions[index].Username, permissions[index].Password)
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
	return date[8:10] + "-" + date[5:7]
}

package app

import (
	"encoding/json"
	// "fmt"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/user"
	"github.com/olekukonko/tablewriter"
	"net/http"
	"os"
	"strconv"
)

type appBody struct {
	APICalls int `json:"api_calls"`
	Records  int `json:"records"`
	Storage  int `json:"storage"`
}

// respBody represents response body
type respBody struct {
	Body map[string]appBody `json:"body"`
}

// ShowUserApps shows the list of user apps
func ShowUserApps() error {
	// get name id mapping
	apps, err := user.GetUserApps()
	if err != nil {
		return err
	}
	// get more details
	req, err := http.NewRequest("GET", common.AccAPIURL+"/user/apps/metrics", nil)
	if err != nil {
		return err
	}
	err = session.AttachCookiesToRequest(req)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// decode response
	var res respBody
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return err
	}
	// output
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "API Calls", "Records", "Storage"})
	for appID, appData := range res.Body {
		table.Append([]string{
			appID, common.GetKeyForValue(apps, appID), strconv.Itoa(appData.APICalls),
			strconv.Itoa(appData.Records), strconv.Itoa(appData.Storage),
		})
	}
	table.Render()
	return nil
}

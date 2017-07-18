package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
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

type appDetailBody struct {
	AppName   string `json:"appname"`
	ESVersion string `json:"es_version"`
}

type appRespBody struct {
	Body appDetailBody `json:"body"`
}

// ShowUserApps shows the list of user apps
func ShowUserApps() error {
	spinner.StartText("Loading user app list")
	defer spinner.Stop()
	// get name id mapping
	apps, err := user.GetUserApps()
	if err != nil {
		return err
	}
	spinner.StartText("Fetching app data")
	// get more details
	req, err := http.NewRequest("GET", common.AccAPIURL+"/user/apps/metrics", nil)
	if err != nil {
		return err
	}
	resp, err := session.SendRequest(req)
	if err != nil {
		return err
	}
	spinner.Stop()
	// decode response
	var res respBody
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return err
	}
	// output
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "API Calls", "Records", "Storage (KB)"})
	for appID, appData := range res.Body {
		table.Append([]string{
			appID, common.GetKeyForValue(apps, appID), strconv.Itoa(appData.APICalls),
			strconv.Itoa(appData.Records), strconv.Itoa(common.SizeInKB(appData.Storage)),
		})
	}
	table.Render()
	return nil
}

// ShowAppDetails shows the app details
func ShowAppDetails(app string, perms bool, metrics bool) error {
	spinner.StartText("Loading app details")
	defer spinner.Stop()
	// fetch app basic details
	app, err := EnsureAppID(app)
	if err != nil {
		return err
	}
	// get app details
	req, err := http.NewRequest("GET", common.AccAPIURL+"/app/"+app, nil)
	if err != nil {
		return err
	}
	resp, err := session.SendRequest(req)
	if err != nil {
		return err
	}
	spinner.Stop()
	// decode response
	var res appRespBody
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return err
	}
	// output
	fmt.Printf("ID:         %s\n", app)
	fmt.Printf("Name:       %s\n", res.Body.AppName)
	fmt.Printf("ES Version: %s\n", res.Body.ESVersion)

	if perms {
		err = ShowAppPerms(app)
	}
	if metrics {
		err = ShowAppMetrics(app)
	}
	return err
}

// EnsureAppID make sures `app` is id
func EnsureAppID(app string) (string, error) {
	// check if num https://stackoverflow.com/questions/22593259/
	if _, err := strconv.Atoi(app); err == nil {
		return app, nil // return as is
	}
	// convert to appID
	apps, err := user.GetUserApps()
	if err != nil {
		return "", err
	}
	appID, ok := apps[app]
	if ok {
		return appID, nil
	}
	return "", errors.New("App with name " + app + " not found.")
}

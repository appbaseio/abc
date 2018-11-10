package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
	"github.com/appbaseio/abc/appbase/user"
	"github.com/olekukonko/tablewriter"
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
	AppName   string    `json:"appname"`
	ESVersion string    `json:"es_version"`
	Owner     string    `json:"owner"`
	Users     []string  `json:"users"`
	CreatedAt time.Time `json:"created_at, string"`
}

type appRespBody struct {
	Body appDetailBody `json:"body"`
}

// ShowUserApps shows the list of user apps
func ShowUserApps(sortOption string) error {
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
	// sort
	var as appsSorter
	as.key = sortOption
	as.apps = make([]fullApp, 0)
	for appID, appData := range res.Body {
		var fa fullApp
		fa.id, fa.name = appID, common.GetKeyForValue(apps, appID)
		fa.APICalls, fa.Records, fa.Storage = appData.APICalls, appData.Records, appData.Storage
		as.apps = append(as.apps, fa)
	}
	sort.Sort(as)
	// output
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Name", "API Calls", "Records", "Storage (KB)"})
	for _, app := range as.apps {
		table.Append([]string{
			app.id, app.name, strconv.Itoa(app.APICalls),
			strconv.Itoa(app.Records), strconv.Itoa(common.SizeInKB(app.Storage)),
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
	// prepare output
	common.RemoveDuplicates(&res.Body.Users)
	users := ""
	for _, user := range res.Body.Users {
		users += user + ", "
	}
	// output
	fmt.Printf("ID:         %s\n", app)
	fmt.Printf("Name:       %s\n", res.Body.AppName)
	fmt.Printf("Owner:      %s\n", res.Body.Owner)
	fmt.Printf("Users:      %s\n", users[:len(users)-2])
	fmt.Printf("ES Version: %s\n", res.Body.ESVersion)
	fmt.Printf("Created on: %s\n", res.Body.CreatedAt.Format("Mon Jan _2 15:04:05 2006"))

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

// EnsureAppName make sures `app` is name
func EnsureAppName(app string) (string, error) {
	// check if string https://stackoverflow.com/questions/22593259/
	if _, err := strconv.Atoi(app); err != nil {
		return app, nil // return as is
	}
	// convert to name
	apps, err := user.GetUserApps()
	if err != nil {
		return "", err
	}
	name := common.GetKeyForValue(apps, app)
	if name == "" {
		return "", errors.New("App with ID " + app + " not found.")
	}
	return name, nil
}

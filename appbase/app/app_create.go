package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
)

type createdAppDetails struct {
	ID int64 `json:"id"`
}

type createRespBody struct {
	Body createdAppDetails `json:"body"`
}

// RunAppCreate runs app create command
func RunAppCreate(appName string, esVersion string, category string) error {
	spinner.Start()
	defer spinner.Stop()
	// create app
	body := fmt.Sprintf(`{"category": "%s", "es_version": "%s"}`, category, esVersion)
	req, err := http.NewRequest("PUT", common.AccAPIURL+"/app/"+appName, strings.NewReader(body))
	if err != nil {
		return err
	}
	resp, err := session.SendRequest(req)
	if err != nil {
		return err
	}
	spinner.Stop()
	// status code not 200
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("There was an error %s", string(bodyBytes))
	}
	// decode
	var res createRespBody
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return err
	}
	// output
	fmt.Printf("ID:    %d\n", res.Body.ID)
	fmt.Printf("Name:  %s\n", appName)
	// permissions
	ShowAppPerms(strconv.FormatInt(res.Body.ID, 10))
	return nil
}

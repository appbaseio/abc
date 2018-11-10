package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
)

type deleteRespBody struct {
	Status     statusDetailsBody `json:"status"`
	Deployment string            `json:"deployment"`
}

// RunClusterDelete deletes a cluster using the cluster ID
func RunClusterDelete(clusterID string) error {
	spinner.Start()
	defer spinner.Stop()

	req, err := http.NewRequest("DELETE", common.AccAPIURL+"/v1/_delete/"+clusterID, nil)
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
	var res deleteRespBody
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return err
	}
	// output
	fmt.Println(res.Status.Message)
	return nil
}

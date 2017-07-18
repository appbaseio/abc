package app

import (
	"encoding/json"
	"fmt"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
	"net/http"
)

type deleteRespBody struct {
	Message string `json:"message"`
}

// RunAppDelete runs app create command
func RunAppDelete(app string) error {
	spinner.Start()
	defer spinner.Stop()
	// get appID from name
	app, err := EnsureAppID(app)
	if err != nil {
		return err
	}
	// request
	req, err := http.NewRequest(http.MethodDelete, common.AccAPIURL+"/app/"+app, nil)
	if err != nil {
		return err
	}
	resp, err := session.SendRequest(req)
	if err != nil {
		return err
	}
	spinner.Stop()
	// decode
	var res deleteRespBody
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return err
	}
	// output
	fmt.Println(res.Message)
	return nil
}

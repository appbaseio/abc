package user

import (
	"encoding/json"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"net/http"
)

// userDetails represents extra details of user
type userDetails struct {
	Name string `json:"name"`
}

// userBody represents Appbase.io user
type userBody struct {
	Email   string      `json:"email"`
	Details userDetails `json:"details"`
}

// respBody represents response body
type respBody struct {
	Body userBody `json:"body"`
}

func getCurrentUser() (userBody, error) {
	var user respBody
	req, err := http.NewRequest("GET", common.AccAPIURL+"/user", nil)
	if err != nil {
		return user.Body, err
	}
	err = session.AttachCookiesToRequest(req)
	if err != nil {
		return user.Body, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return user.Body, err
	}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&user)
	return user.Body, err
}

// GetUserEmail returns user's email
func GetUserEmail() (string, error) {
	user, err := getCurrentUser()
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

package session

import (
	b64 "encoding/base64"
	"encoding/json"
	"github.com/appbaseio/abc/log"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
)

// LoadUserSessionAsString loads session string
func LoadUserSessionAsString() (string, error) {
	sessionFile, err := getSessionFilePath()
	if err != nil {
		return "", err
	}
	// file found, load
	if _, err := os.Stat(sessionFile); err == nil {
		// read session
		dat, err := ioutil.ReadFile(sessionFile)
		if err != nil {
			return "", err
		}
		return string(dat), nil
	}
	// load env token else
	token := os.Getenv("ABC_TOKEN")
	if len(token) > 0 {
		return token, nil
	}
	return "", err
}

// LoadUserSessionAsCookie loads and returns arrays of cookies
func LoadUserSessionAsCookie() ([3]http.Cookie, error) {
	cookies := [3]http.Cookie{}
	sessionData, err := LoadUserSessionAsString()
	if err != nil {
		return cookies, err
	}
	sDec, err := b64.StdEncoding.DecodeString(sessionData)
	log.Debugf("Decoded Token: %s", string(sDec))
	if err != nil {
		return cookies, err
	}
	// Decode JSON
	type Cookie struct {
		Ga            string `json:"_ga"`
		AppbaseAccAPI string `json:"appbase_accapi"`
		Session       string `json:"session"`
	}
	var ck Cookie
	err = json.Unmarshal(sDec, &ck)
	if err != nil {
		return cookies, err
	}
	cookies[0] = http.Cookie{Name: "_ga", Value: ck.Ga}
	cookies[1] = http.Cookie{Name: "appbase_accapi", Value: ck.AppbaseAccAPI}
	cookies[2] = http.Cookie{Name: "session", Value: ck.Session}
	return cookies, nil
}

// SaveUserSession saves user session information
func SaveUserSession(data string) error {
	sessionFile, err := getSessionFilePath()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(sessionFile, []byte(data), 0644)
	return err
}

// attachCookiesToRequest attaches cookies to a request
func attachCookiesToRequest(req *http.Request) error {
	cookies, err := LoadUserSessionAsCookie()
	log.Debugf("Cookies: %s", cookies)
	if err != nil {
		return err
	}
	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}
	return nil
}

// SendRequest sends a request with cookies
func SendRequest(req *http.Request) (*http.Response, error) {
	var dumResp *http.Response
	err := attachCookiesToRequest(req)
	if err != nil {
		return dumResp, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return dumResp, err
	}
	return resp, nil
}

func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func getSessionFilePath() (string, error) {
	homeDir, err := getUserHomeDir()
	if err != nil {
		return "", err
	}
	log.Debugf("Home Dir: %s", homeDir)
	return homeDir + "/.abcsession", nil
}

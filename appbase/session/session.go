package session

import (
	b64 "encoding/base64"
	"encoding/json"
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

// AttachCookiesToRequest attaches cookies to a request
func AttachCookiesToRequest(req *http.Request) error {
	cookies, err := LoadUserSessionAsCookie()
	if err != nil {
		return err
	}
	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}
	return nil
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
	return homeDir + "/.abcsession", nil
}

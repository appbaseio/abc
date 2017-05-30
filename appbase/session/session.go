package session

import (
	"io/ioutil"
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

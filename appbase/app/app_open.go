package app

import (
	"errors"
	"fmt"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/login"
	"github.com/appbaseio/abc/appbase/spinner"
	"strings"
)

// OpenAppDataView ...
func OpenAppDataView(app string) error {
	spinner.StartText("Making Dejavu URL")
	defer spinner.Stop()
	url, err := GetAppURL(app)
	if err != nil {
		return err
	}
	djURL, err := common.MakeDejavuURL(url)
	if err != nil {
		return err
	}
	spinner.Stop()
	fmt.Printf("Opening url %s\n", djURL)
	err = common.OpenURL(djURL)
	if err != nil {
		return err
	}
	return nil
}

// OpenAppQueryView ...
func OpenAppQueryView(app string) error {
	spinner.StartText("Making Mirage URL")
	defer spinner.Stop()
	url, err := GetAppURL(app)
	if err != nil {
		return err
	}
	mgURL, err := common.MakeMirageURL(url)
	if err != nil {
		return err
	}
	spinner.Stop()
	fmt.Printf("Opening url %s\n", mgURL)
	err = common.OpenURL(mgURL)
	if err != nil {
		return err
	}
	return nil
}

// GetAppURL returns the full url of an app
func GetAppURL(appName string) (string, error) {
	if !login.IsUserAuthenticated() {
		return "", errors.New("User not logged in. Unable to fetch app url")
	}
	appID, err := EnsureAppID(appName)
	if err != nil {
		return "", err
	}
	perms, err := GetAppPerms(appID)
	if err != nil {
		return "", err
	}
	for _, perm := range perms {
		if strings.Contains(strings.ToLower(perm.Description), "admin") {
			return fmt.Sprintf("https://%s:%s@scalr.api.appbase.io/%s", perm.Username, perm.Password, appName), nil
		}
	}
	return "", fmt.Errorf("App with name %s not found", appName)
}

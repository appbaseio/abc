// +build !oss

package importer

import (
	"errors"
	"fmt"
	"github.com/appbaseio/abc/appbase/app"
	"github.com/appbaseio/abc/appbase/login"
	// "github.com/appbaseio/abc/log"
	"strings"
)

// GetAppURL returns the full url of an app
func GetAppURL(appName string) (string, error) {
	if !login.IsUserAuthenticated() {
		return "", errors.New("User not logged in. Unable to fetch app url")
	}
	appID, err := app.EnsureAppID(appName)
	if err != nil {
		return "", err
	}
	perms, err := app.GetAppPerms(appID)
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

package login

import (
	"github.com/aviaryan/abc/appbase/session"
)

// IsUserAuthenticated checks if user is logged in or not
func IsUserAuthenticated() bool {
	data, err := session.LoadUserSessionAsString()
	return err == nil && data != ""
}

// StartUserLogin starts user login process
func StartUserLogin() error {
	return nil
}

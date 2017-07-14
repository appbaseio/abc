package main

import (
	"fmt"
	"github.com/appbaseio/abc/appbase/login"
	"github.com/appbaseio/abc/appbase/user"
)

// runLogin runs the login command
func runLogin(args []string) error {
	flagset := baseFlagSet("login")
	basicUsage := "abc login [google|github]"
	flagset.Usage = usageFor(flagset, basicUsage)
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	switch len(args) {
	case 0:
		if isLoggedIn() {
			return user.ShowUserEmail()
		}
	case 1:
		fmt.Println("Logging in..")
		return login.StartUserLogin(args[0])
	default:
		showShortHelp(basicUsage)
	}
	return nil
}

// isLoggedIn checks if a user is logged in or not, prints message if not logged in
func isLoggedIn() bool {
	if login.IsUserAuthenticated() {
		return true
	}
	fmt.Println("Not logged in")
	return false
}

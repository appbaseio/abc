package main

import (
	"fmt"

	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/login"
	"github.com/appbaseio/abc/appbase/user"
)

// runLogin runs the login command
func runLogin(args []string) error {
	flagset := baseFlagSet("login")
	basicUsage := "abc login [google|github|gitlab|api] [-c|--credentials] {username:password}"
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
		if common.StringInSlice(args[0], []string{"google", "github", "gitlab"}) {
			fmt.Println("Logging in..")
			return login.StartUserLogin(args[0])
		}
		showShortHelp(basicUsage)
	case 2:
		if args[0] == "api" {
			fmt.Println("Logging in..")
			return login.StartUserLoginBasicAuth(args[1])
		}
		showShortHelp(basicUsage)
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
	fmt.Println("user not logged in, use --help to see usage on how to login.")
	return false
}

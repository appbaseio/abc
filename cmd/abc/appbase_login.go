package main

import (
	"fmt"
	"github.com/aviaryan/abc/appbase/login"
	"github.com/aviaryan/abc/appbase/user"
)

// runLogin runs the login command
func runLogin(args []string) error {
	flagset := baseFlagSet("login")
	flagset.Usage = usageFor(flagset, "abc login [google|github]")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	switch len(args) {
	case 0:
		if login.IsUserAuthenticated() {
			email, _ := user.GetUserEmail()
			// FIXME: token can be tampered and so it won't work
			fmt.Println("Logged in as", email)
		} else {
			fmt.Println("Not logged in.")
		}
	case 1:
		fmt.Println("Logging in..")
		return login.StartUserLogin(args[0])
	default:
		fmt.Println("Wrong number of parameters. See help (--help).")
	}
	return nil
}

package main

import (
	"fmt"
	"github.com/aviaryan/abc/appbase/login"
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
			fmt.Println("Authenticated user")
		} else {
			fmt.Println("UnAuthenticated user")
		}
	case 1:
		fmt.Println("Logging in")
	default:
		fmt.Println("Wrong parameters. See help.")
	}
	return nil
}

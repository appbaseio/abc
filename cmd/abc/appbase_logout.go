package main

import (
	"fmt"
	"github.com/appbaseio/abc/appbase/logout"
)

// runLogout runs the logout command
func runLogout(args []string) error {
	flagset := baseFlagSet("logout")
	flagset.Usage = usageFor(flagset, "abc logout")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	switch len(args) {
	case 0:
		if isLoggedIn() {
			return logout.UserLogout()
		}
	default:
		fmt.Println("Wrong number of parameters. See help (--help).")
	}
	return nil
}

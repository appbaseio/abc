package main

import (
	"github.com/appbaseio/abc/appbase/user"
)

func runUser(args []string) error {
	flagset := baseFlagSet("user")
	basicUsage := "abc user"
	flagset.Usage = usageFor(flagset, basicUsage)
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	if len(args) == 0 {
		if isLoggedIn() {
			return user.ShowUserDetails()
		}
	} else {
		showShortHelp(basicUsage)
	}
	return nil
}

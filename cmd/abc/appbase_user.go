package main

import (
	"fmt"
	"github.com/appbaseio/abc/appbase/user"
)

func runUser(args []string) error {
	flagset := baseFlagSet("user")
	flagset.Usage = usageFor(flagset, "abc user")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	if len(args) == 0 {
		if isLoggedIn() {
			return user.ShowUserDetails()
		}
	} else {
		fmt.Println("No such option. See --help")
	}
	return nil
}

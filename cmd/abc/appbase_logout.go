package main

import (
	"github.com/appbaseio/abc/appbase/logout"
)

// runLogout runs the logout command
func runLogout(args []string) error {
	flagset := baseFlagSet("logout")
	basicUsage := "abc logout [--all]"
	flagset.Usage = usageFor(flagset, basicUsage)
	doAll := flagset.Bool("all", false, "logout from browser as well")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	switch len(args) {
	case 0:
		if isLoggedIn() {
			return logout.UserLogout(*doAll)
		}
	default:
		showShortHelp(basicUsage)
	}
	return nil
}

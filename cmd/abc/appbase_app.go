package main

import (
	"fmt"
	"github.com/appbaseio/abc/appbase/app"
)

// runApps runs `apps` command
func runApps(args []string) error {
	flagset := baseFlagSet("apps")
	flagset.Usage = usageFor(flagset, "abc apps")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	switch len(args) {
	case 0:
		return app.ShowUserApps()
	default:
		fmt.Println("No such option. See --help")
	}
	return nil
}

// runApp runs `app` command
func runApp(args []string) error {
	flagset := baseFlagSet("app")
	flagset.Usage = usageFor(flagset, "abc app ID|Appname")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	if len(args) < 4 {
		return app.ShowAppDetails(args[0], args[1:]...)
	}
	fmt.Println("No such option. See --help")
	return nil
}

// runCreate runs `create` command
func runCreate(args []string) error {
	flagset := baseFlagSet("create")
	flagset.Usage = usageFor(flagset, "abc create [-es2|-es5] [-category=category] AppName")
	// https://gobyexample.com/command-line-flags
	isEs5 := flagset.Bool("es5", false, "is app es5")
	category := flagset.String("category", "generic", "category for app")

	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	if len(args) == 1 {
		if *isEs5 {
			return app.RunAppCreate(args[0], "5", *category)
		}
		return app.RunAppCreate(args[0], "2", *category)
	}
	fmt.Println("No such option. See --help")
	return nil
}

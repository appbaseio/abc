package main

import (
	"fmt"
	"github.com/appbaseio/abc/appbase/app"
	"github.com/appbaseio/abc/appbase/common"
)

// runApps runs `apps` command
func runApps(args []string) error {
	flagset := baseFlagSet("apps")
	basicUsage := "abc apps"
	flagset.Usage = usageFor(flagset, basicUsage)
	sort := flagset.String("sort", "id", "sort by id, name, api-calls, records, storage")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	switch len(args) {
	case 0:
		if common.StringInSlice(*sort, app.SortOptions) {
			return app.ShowUserApps(*sort)
		}
		fmt.Printf("Invalid parameter '%s' passed to sort. See --help\n", *sort)
	default:
		showShortHelp(basicUsage)
	}
	return nil
}

// runApp runs `app` command
func runApp(args []string) error {
	flagset := baseFlagSet("app")
	basicUsage := "abc app [-c|--creds] [-m|--metrics] [ID|Appname]"
	flagset.Usage = usageFor(flagset, basicUsage)
	creds := flagset.BoolP("creds", "c", false, "show app credentials")
	metrics := flagset.BoolP("metrics", "m", false, "show app metrics")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	if len(args) == 1 {
		return app.ShowAppDetails(args[0], *creds, *metrics)
	}
	showShortHelp(basicUsage)
	return nil
}

// runCreate runs `create` command
func runCreate(args []string) error {
	flagset := baseFlagSet("create")
	basicUsage := "abc create [--es2|--es6] [--category=category] AppName"
	flagset.Usage = usageFor(flagset, basicUsage)
	// https://gobyexample.com/command-line-flags
	isEs6 := flagset.Bool("es6", false, "is app es6")
	isEs2 := flagset.Bool("es2", true, "is app es2")
	category := flagset.String("category", "generic", "category for app")

	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	if len(args) == 1 {
		if *isEs6 {
			return app.RunAppCreate(args[0], "6", *category)
		} else if *isEs2 {
			return app.RunAppCreate(args[0], "2", *category)
		} else {
			fmt.Println("App needs to be ES2 or ES6")
			return nil
		}
	}
	showShortHelp(basicUsage)
	return nil
}

// runDelete runs `delete` command
func runDelete(args []string) error {
	flagset := baseFlagSet("delete")
	basicUsage := "abc delete [AppID|AppName]"
	flagset.Usage = usageFor(flagset, basicUsage)
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()
	if len(args) == 1 {
		return app.RunAppDelete(args[0])
	}
	showShortHelp(basicUsage)
	return nil
}

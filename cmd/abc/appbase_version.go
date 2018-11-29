package main

import (
	"fmt"
	"runtime"

	"github.com/appbaseio/abc/imports"
)

var version = "0.7.0"
var variant = imports.BuildName

// runVersion runs the logout command
func runVersion(args []string) error {
	flagset := baseFlagSet("version")
	basicUsage := "abc logout"
	flagset.Usage = usageFor(flagset, basicUsage)
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	switch len(args) {
	case 0:
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Variant:    %s\n", variant)
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("OS:         %s\n", runtime.GOOS)
		fmt.Printf("Arch:       %s\n", runtime.GOARCH)
	default:
		showShortHelp(basicUsage)
	}
	return nil
}

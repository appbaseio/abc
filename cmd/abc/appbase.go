package main

import (
	"fmt"
	"github.com/appbaseio/abc/imports"
	"os"
	"strings"
)

// usageAppbase adds help on using Appbase commands
func usageAppbase() {
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "APPBASE\n")
	fmt.Fprintf(os.Stderr, "  login     login into appbase.io\n")
	fmt.Fprintf(os.Stderr, "  user      get user details\n")
	fmt.Fprintf(os.Stderr, "  apps      display user apps\n")
	fmt.Fprintf(os.Stderr, "  app       display app details\n")
	fmt.Fprintf(os.Stderr, "  create    create app\n")
	fmt.Fprintf(os.Stderr, "  delete    delete app\n")
	if imports.IsPrivate {
		fmt.Fprintf(os.Stderr, "  import    import data to appbase app\n")
	}
}

// provisionAppbaseCLI provisions the addon appbase CLI
func provisionAppbaseCLI(command string) func([]string) error {
	var run func([]string) error
	// match command
	switch strings.ToLower(command) {
	case "login":
		run = runLogin
	case "user":
		run = runUser
	case "apps":
		run = runApps
	case "app":
		run = runApp
	case "create":
		run = runCreate
	case "delete":
		run = runDelete
	default:
		usage()
		os.Exit(1)
	}
	// safe actually as we already exit when usage is called
	return run
}

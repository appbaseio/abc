package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/appbaseio/abc/imports"
)

// usageAppbase adds help on using Appbase commands
func usageAppbase() {
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "COMMANDS\n")
	fmt.Fprintf(os.Stderr, "  login     login into appbase.io\n")
	fmt.Fprintf(os.Stderr, "  user      get user details\n")
	fmt.Fprintf(os.Stderr, "  apps      display user apps\n")
	fmt.Fprintf(os.Stderr, "  app       display app details\n")
	fmt.Fprintf(os.Stderr, "  cluster   display cluster details\n")
	fmt.Fprintf(os.Stderr, "  create    create app/cluster\n")
	fmt.Fprintf(os.Stderr, "  delete    delete app/cluster\n")
	fmt.Fprintf(os.Stderr, "  logout    logout session\n")
	if imports.IsPrivate {
		fmt.Fprintf(os.Stderr, "  import    import data to appbase app\n")
	}
	fmt.Fprintf(os.Stderr, "  version   show build details\n")
	fmt.Fprintf(os.Stderr, "  license   show project license and credits\n")
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
	case "cluster":
		run = runCluster
	case "create":
		run = runCreate
	case "delete":
		run = runDelete
	case "logout":
		run = runLogout
	case "version", "--version":
		run = runVersion
	case "license":
		run = runLicense
	default:
		usage()
		os.Exit(1)
	}
	// safe actually as we already exit when usage is called
	return run
}

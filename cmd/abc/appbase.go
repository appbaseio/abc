package main

import (
	"fmt"
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
	default:
		usage()
		os.Exit(1)
	}
	// safe actually as we already exit when usage is called
	return run
}

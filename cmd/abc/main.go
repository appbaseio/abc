package main

import (
	"fmt"
	flag "github.com/ogier/pflag"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/appbaseio/abc/imports"
	_ "github.com/appbaseio/abc/imports/all"
	"github.com/appbaseio/abc/log"
)

const (
	defaultPipelineFile = "pipeline.js"
)

var version = "0.1.0"
var variant = imports.BuildName

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE\n")
	fmt.Fprintf(os.Stderr, "  %s <command> [flags]\n", os.Args[0])
	usageAppbase()
	// variant
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "VERSION\n")
	fmt.Fprintf(os.Stderr, "  %s (%s)\n", version, variant)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var run func([]string) error
	command := strings.ToLower(os.Args[1])
	if command == "import" {
		if imports.IsPrivate {
			run = runImport
		} else {
			usage()
			os.Exit(1)
		}
	} else {
		run = provisionAppbaseCLI(command)
	}

	if err := run(os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func baseFlagSet(setName string) *flag.FlagSet {
	cmdFlags := flag.NewFlagSet(setName, flag.ExitOnError)
	log.AddFlags(cmdFlags)
	return cmdFlags
}

func usageFor(fs *flag.FlagSet, short string) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "USAGE\n")
		fmt.Fprintf(os.Stderr, "  %s\n", short)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "FLAGS\n")
		w := tabwriter.NewWriter(os.Stderr, 0, 2, 2, ' ', 0)
		fs.VisitAll(func(f *flag.Flag) {
			fmt.Fprintf(w, "\t--%s=%s\t%s\n", f.Name, f.DefValue, f.Usage)
		})
		w.Flush()
		fmt.Fprintf(os.Stderr, "\n")
	}
}

func showShortHelp(short string) {
	fmt.Fprintf(os.Stderr, "USAGE\n")
	fmt.Fprintf(os.Stderr, "  %s\n\n", short)
	fmt.Println("Use --help option for more info.")
}

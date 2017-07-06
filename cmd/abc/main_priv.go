// +build !oss

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
	_ "github.com/appbaseio/abc/private/function/all"
)

const (
	defaultPipelineFile = "pipeline.js"
)

var version = "0.0.0"
var variant = imports.BuildName

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE\n")
	fmt.Fprintf(os.Stderr, "  %s <command> [flags]\n", os.Args[0])
	usageAppbase()
	// private options
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "IMPORTER\n")
	fmt.Fprintf(os.Stderr, "  run       run pipeline loaded from a file\n")
	fmt.Fprintf(os.Stderr, "  test      display the compiled nodes without starting a pipeline\n")
	fmt.Fprintf(os.Stderr, "  about     show information about available adaptors\n")
	fmt.Fprintf(os.Stderr, "  xlog      manage the commit log\n")
	fmt.Fprintf(os.Stderr, "  offset    manage the offset for sinks\n")

	// variant
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "VARIANT\n")
	fmt.Fprintf(os.Stderr, "  %s\n", variant)
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var run func([]string) error
	switch strings.ToLower(os.Args[1]) {
	case "run":
		run = runRun
	case "test":
		run = runTest
	case "about":
		run = runAbout
	case "xlog":
		run = runXlog
	case "offset":
		run = runOffset
	case "import":
		run = runImport
	default:
		run = provisionAppbaseCLI(os.Args[1])
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

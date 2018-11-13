package main

import (
	"fmt"

	"github.com/appbaseio/abc/appbase/app"
	"github.com/appbaseio/abc/appbase/cluster"
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
	basicUsage := "abc app [-c|--creds] [-m|--metrics] [--data-view] [-a| --analytics] [ID|Appname]"
	flagset.Usage = usageFor(flagset, basicUsage)
	analytics := flagset.BoolP("analytics", "a", false, "show app analytics")
	analyticsEndpoint := flagset.String("endpoint", "overview", "the analytics endpoint to be queried")
	creds := flagset.BoolP("creds", "c", false, "show app credentials")
	metrics := flagset.BoolP("metrics", "m", false, "show app metrics")
	dataView := flagset.Bool("data-view", false, "open app data view using Dejavu")
	queryView := flagset.Bool("query-view", false, "open app query view using Mirage")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	if len(args) == 1 {
		if *dataView {
			return app.OpenAppDataView(args[0])
		} else if *queryView {
			return app.OpenAppQueryView(args[0])
		} else if *analytics {
			return app.ShowAppAnalytics(args[0], *analyticsEndpoint)
		}
		return app.ShowAppDetails(args[0], *creds, *metrics)
	}
	showShortHelp(basicUsage)
	return nil
}

// runCreate runs `create` command
func runCreate(args []string) error {
	flagset := baseFlagSet("create")
	basicUsage := "abc create [--es2|--es6] [--category=category] [--cluster|-c] [--interactive|-i] [--loc] [--vmsize] [--plan] [--ssh] [--provider] [--nodes] [--version] [--volume] AppName|ClusterName"
	flagset.Usage = usageFor(flagset, basicUsage)
	// https://gobyexample.com/command-line-flags
	isEs6 := flagset.Bool("es6", false, "is app es6")
	isEs2 := flagset.Bool("es2", true, "is app es2")
	category := flagset.String("category", "generic", "category for app")

	// Cluster specific flags
	isCluster := flagset.BoolP("cluster", "c", false, "cluster mode")
	interactiveMode := flagset.BoolP("interactive", "i", false, "interactive mode for cluster creation")
	location := flagset.String("loc", "", "location of the cluster")
	vmSize := flagset.String("vmsize", "", "size of the VMs")
	pricingPlan := flagset.String("plan", "", "pricing plan")
	sshPublicKey := flagset.String("ssh", "", "SSH public key")
	provider := flagset.String("provider", "", "Accepted values are azure or gke")
	nodes := flagset.Int("nodes", 1, "number of ES nodes")
	esVersion := flagset.String("version", "", "A valid ES version")
	volumeSize := flagset.Int("volume", 1, "volume size. Valid values are from 1-500")

	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	if len(args) == 1 {
		if *isCluster {
			if *interactiveMode {
				requestBody := cluster.BuildRequestBodyInteractive()
				cluster.DeployCluster(requestBody)
			} else {
				requestBody := cluster.BuildRequestBody(args[0], *location, *vmSize, *pricingPlan, *sshPublicKey, *provider, *nodes, *esVersion, *volumeSize)
				cluster.DeployCluster(requestBody)
			}
		} else if *isEs6 {
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
	basicUsage := "abc delete [AppID|AppName|ClusterID] [--cluster=ClusterName]"
	flagset.Usage = usageFor(flagset, basicUsage)
	getCluster := flagset.Bool("cluster", false, "for deleting clusters instead of apps")
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()
	if len(args) == 1 {
		if *getCluster == true {
			return cluster.RunClusterDelete(args[0])
		}
		return app.RunAppDelete(args[0])
	}
	showShortHelp(basicUsage)
	return nil
}

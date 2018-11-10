package main

import (
	"github.com/appbaseio/abc/appbase/cluster"
)

func runCluster(args []string) error {
	flagset := baseFlagSet("cluster")
	basicUsage := "abc cluster [ClusterID]"
	flagset.Usage = usageFor(flagset, basicUsage)

	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	if len(args) == 1 {
		return cluster.ShowClusterDetails(args[0])
	}
	showShortHelp(basicUsage)
	return nil
}

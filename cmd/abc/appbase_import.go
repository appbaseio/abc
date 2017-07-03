// +build !oss

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/appbaseio/abc/imports/adaptor"
	"github.com/appbaseio/abc/log"
	"os"
	"strconv"
	"time"
)

const importInfo string = `
	abc import --type {DBType} --source {URI} [-t|--tail] [Uri|AppID|Appname]
`

// runImport runs the import command
func runImport(args []string) error {
	flagset := baseFlagSet("import")
	flagset.Usage = usageFor(flagset, importInfo)

	// custom flags
	tail := flagset.BoolP("tail", "t", false, "allow tail feature")
	srcType := flagset.String("type", "postgres", "type of source database")
	srcURL := flagset.String("source", "http://user:pass@host:port/db", "url of source database")
	typeName := flagset.String("typename", "mytype", "[csv] typeName to use")
	replicationSlot := flagset.String("replication_slot", "standby_replication_slot",
		"[postgres] replication slot to use")
	timeout := flagset.String("timeout", "10s", "source timeout")
	var destURL string

	// parse args
	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()
	if len(args) == 1 {
		destURL = args[0]
	} else {
		return errors.New("Invalid set of parameters")
	}

	// create source config
	var srcConfig = map[string]interface{}{
		"uri":              *srcURL,
		"tail":             *tail,
		"typeName":         *typeName,
		"replication_slot": *replicationSlot,
		"timeout":          *timeout,
	}

	// write config file
	file, err := writeConfigFile(*srcType, srcConfig, destURL)
	if err != nil {
		return err
	}

	log.Infof("Created temp file %s", file)
	// return nil

	// run config file
	builder, err := newBuilder(file)
	if err != nil {
		return err
	}

	return builder.run()
}

func writeConfigFile(src string, srcConfig map[string]interface{}, destURI string) (string, error) {
	fname := "pipeline_" + strconv.FormatInt(time.Now().Unix(), 10) + ".js"

	if _, err := os.Stat(fname); err == nil {
		log.Errorf("File %s exists, will be overwritten", fname)
	}
	appFileHandle, err := os.Create(fname)
	if err != nil {
		return "", err
	}
	defer appFileHandle.Close()

	args := []string{src, "elasticsearch"}
	var config = make(map[string]interface{})

	nodeName := "source"
	for _, name := range args {
		// set config
		if nodeName == "source" {
			for k, v := range srcConfig {
				config[k] = v
			}
		} else {
			config = map[string]interface{}{}
			config["uri"] = destURI
		}
		// get adaptor
		a, _ := adaptor.GetAdaptor(name, config)
		// get config json
		b, err := json.Marshal(a)
		if err != nil {
			return "", err
		}
		confJSON := string(b)
		// save to file
		appFileHandle.WriteString(fmt.Sprintf("var %s = %s(%s)\n\n", nodeName, name, confJSON))
		nodeName = "sink"
	}
	appFileHandle.WriteString(`t.Source("source", source, "/.*/").Save("sink", sink, "/.*/")`)
	appFileHandle.WriteString("\n")

	return fname, nil
}

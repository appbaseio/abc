// +build !oss

package main

import (
	"encoding/json"
	"fmt"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/importer"
	"github.com/appbaseio/abc/imports/adaptor"
	"github.com/appbaseio/abc/log"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// GLOBALS
// map from real input params to what goes in writeConfig
var srcParamMap = map[string]string{
	"src.uri":          "uri",
	"src.type":         "_name_",
	"tail":             "tail",
	"replication_slot": "replication_slot",
	"typename":         "typeName",
	// "timeout":          "timeout",
	"transform_file": "_transform_",
}

var destParamMap = map[string]string{
	"dest.uri":  "uri",
	"dest.type": "_name_",
}

const basicUsage string = `abc import --src.type {DBType} --src.uri {URI} [-t|--tail] [Uri|Appname]`

// runImport runs the import command
func runImport(args []string) error {
	flagset := baseFlagSet("import")
	flagset.Usage = usageFor(flagset, basicUsage)

	// custom flags
	tail := flagset.BoolP("tail", "t", false, "allow tail feature")
	srcType := flagset.String("src.type", "postgres", "type of source database")
	srcURL := flagset.String("src.uri", "http://user:pass@host:port/db", "url of source database")
	typeName := flagset.String("typename", "mytype", "[csv] typeName to use")
	replicationSlot := flagset.String("replication-slot", "standby_replication_slot",
		"[postgres] replication slot to use")
	// timeout := flagset.String("timeout", "10s", "source timeout")
	srcRegex := flagset.String("src.filter", ".*", "Namespace filter for source")
	test := flagset.Bool("test", false, `if set to true, only pipeline is created and sync is not started. 
		Useful for checking your configuration`)

	transformFile := flagset.String("transform-file", "", "transform file to use")

	// use external config
	config := flagset.String("config", "", "Path to external config file, if specified, only that is used")

	var destURL string
	// parse args
	if err := flagset.Parse(args); err != nil {
		return err
	}

	// use the config file
	if *config != "" {
		file, err := genPipelineFromEnv(*config)
		if err != nil {
			return err
		}
		return execBuilder(file, *test)
	}

	// use command line params
	args = flagset.Args()
	if len(args) == 1 {
		destURL = args[0]
	} else {
		showShortHelp(basicUsage)
		return nil
	}

	// create source config
	var srcConfig = map[string]interface{}{
		"_name_":           *srcType,
		"uri":              *srcURL,
		"tail":             *tail,
		"typeName":         *typeName,
		"replication_slot": *replicationSlot,
		// "timeout":          *timeout,
		"srcRegex":    *srcRegex,
		"_transform_": *transformFile,
	}

	var destConfig = map[string]interface{}{
		"uri":    destURL,
		"_name_": "elasticsearch",
	}

	// write config file
	file, err := writeConfigFile(srcConfig, destConfig)
	if err != nil {
		return err
	}

	log.Infof("Created temp file %s", file)
	// return nil

	// run config file
	return execBuilder(file, *test)
}

// execBuilder executes a pipeline file
func execBuilder(file string, isTest bool) error {
	builder, err := newBuilder(file)
	if err != nil {
		return err
	}
	if isTest {
		fmt.Println(builder)
		return nil
	}
	// delete if not a devBuild
	if !common.DevBuild {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}
	return builder.run()
}

// writeConfigFile writes config information in a pipeline file
func writeConfigFile(srcConfig map[string]interface{}, destConfig map[string]interface{}) (string, error) {
	fname := "pipeline_" + strconv.FormatInt(time.Now().Unix(), 10) + ".js"

	if _, err := os.Stat(fname); err == nil {
		log.Errorf("File %s exists, will be overwritten", fname)
	}
	appFileHandle, err := os.Create(fname)
	if err != nil {
		return "", err
	}
	defer appFileHandle.Close()

	args := []string{srcConfig["_name_"].(string), destConfig["_name_"].(string)}
	var config = make(map[string]interface{})

	// check appname as destination uri
	if !strings.Contains(destConfig["uri"].(string), "/") {
		destConfig["uri"], err = importer.GetAppURL(destConfig["uri"].(string))
		if err != nil {
			return "", err
		}
	}
	// check file path as source [json, csv]
	if common.StringInSlice(srcConfig["_name_"].(string), []string{"json", "csv"}) {
		err = common.IsFileValid(srcConfig["uri"].(string))
		if err != nil {
			return "", err
		}
	}

	nodeName := "source"
	for _, name := range args {
		// set config
		if nodeName == "source" {
			for k, v := range srcConfig {
				config[k] = v
			}
		} else {
			config = map[string]interface{}{}
			for k, v := range destConfig {
				config[k] = v
			}
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
	// custom transform file
	if srcConfig["_transform_"] != "" {
		dat, err := ioutil.ReadFile(srcConfig["_transform_"].(string))
		if err != nil {
			return "", err
		}
		appFileHandle.WriteString(string(dat))
	} else {
		// no transform file
		appFileHandle.WriteString(
			fmt.Sprintf(`t.Source("source", source, "/%s/").Save("sink", sink, "/.*/")`,
				srcConfig["srcRegex"]),
		)
	}
	appFileHandle.WriteString("\n")

	return fname, nil
}

// genPipelineFromEnv generates a pipeline file from config file
func genPipelineFromEnv(filename string) (string, error) {
	var config map[string]string
	config, err := godotenv.Read(filename)
	if err != nil {
		return "", err
	}
	// source
	src := map[string]interface{}{
		"srcRegex":    ".*", // custom param defaults
		"_transform_": "",
	}
	for k, v := range srcParamMap {
		if val, ok := config[k]; ok {
			src[v] = val
			// tail should be boolean
			if k == "tail" {
				if val == "true" {
					src[v] = true
				} else {
					src[v] = false
				}
			}
		}
	}
	// sink
	dest := map[string]interface{}{}
	for k, v := range destParamMap {
		if val, ok := config[k]; ok {
			dest[v] = val
		}
	}
	// generate file
	file, err := writeConfigFile(src, dest)
	if err != nil {
		return "", err
	}
	fmt.Printf("Writing %s...\n", file)
	return file, nil
}

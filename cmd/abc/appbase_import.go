// +build !oss

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/appbaseio/abc/appbase/app"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/imports/adaptor"
	"github.com/appbaseio/abc/log"
	"github.com/joho/godotenv"
)

// GLOBALS
// map from real input params to what goes in writeConfig
var srcParamMap = map[string]string{
	"src_uri":          "uri",
	"src_type":         "_name_",
	"tail":             "tail",
	"ssl":              "ssl",
	"replication_slot": "replication_slot",
	"typename":         "typeName",
	"src_filter":       "srcRegex",
	"sac_path":         "sacPath",
	// "timeout":          "timeout",
	"transform_file": "_transform_",
	"log_dir":        "log_dir",
}

var destParamMap = map[string]string{
	"dest_uri":      "uri",
	"dest_type":     "_name_",
	"tail":          "tail",
	"request_size":  "request_size",
	"bulk_requests": "bulk_requests",
}

const basicUsage string = `abc import --src.type {DBType} --src.uri {URI} [-t|--tail] [Uri|Appname]`

// runImport runs the import command
func runImport(args []string) error {
	flagset := baseFlagSet("import")
	flagset.Usage = usageFor(flagset, basicUsage)

	// custom flags
	tail := flagset.BoolP("tail", "t", false, "allow tail feature")
	srcType := flagset.String("src_type", "postgres", "type of source database")
	srcURL := flagset.String("src_uri", "http://user:pass@host:port/db", "url of source database")
	typeName := flagset.String("typename", "mytype", "[csv] typeName to use")
	replicationSlot := flagset.String("replication_slot", "standby_replication_slot",
		"[postgres] replication slot to use")
	// timeout := flagset.String("timeout", "10s", "source timeout")
	srcRegex := flagset.String("src_filter", ".*", "Namespace filter for source")
	test := flagset.Bool("test", false, `if set to true, only pipeline is created and sync is not started. 
		Useful for checking your configuration`)
	sacPath := flagset.String("sac_path", "./ServiceAccountKey.json", "Path to firebase service account credentials file")
	ssl := flagset.Bool("ssl", false, "Enable SSL connection to the source.")
	requestSize := flagset.Int64("request_size", 2<<19, "Http request size in bytes, specifically for bulk requests to ES.")
	bulkRequests := flagset.Int("bulk_requests", 1000, "Number of bulk requests to send during a network request to ES.")

	logDir := flagset.String("log_dir", "", "used for storing commit logs")

	transformFile := flagset.String("transform_file", "", "transform file to use")

	verify := flagset.Bool("verify", false, "verify the source and destination connections")

	// use external config
	config := flagset.String("config", "", "Path to external config file, if specified, only that is used")

	var destURL string
	// parse args
	if err := flagset.Parse(args); err != nil {
		return err
	}

	// use the config file
	if *config != "" {
		file, configuredAdaptors, err := genPipelineFromEnv(*config)

		if err != nil {
			return err
		}

		if *verify {
			return verifyConnections(configuredAdaptors)
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
		"srcRegex":         *srcRegex,
		"sacPath":          *sacPath,
		"ssl":              *ssl,
		"_transform_":      *transformFile,
		"log_dir":          *logDir,
	}

	var destConfig = map[string]interface{}{
		"uri":           destURL,
		"_name_":        "elasticsearch",
		"request_size":  *requestSize,
		"bulk_requests": *bulkRequests,
		"tail":          *tail,
	}

	// write config file
	file, configuredAdaptors, err := writeConfigFile(srcConfig, destConfig)
	if err != nil {
		return err
	}

	if *verify {
		return verifyConnections(configuredAdaptors)
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
func writeConfigFile(srcConfig map[string]interface{}, destConfig map[string]interface{}) (string, map[string]adaptor.Adaptor, error) {
	fname := "pipeline_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".js"

	if _, err := os.Stat(fname); err == nil {
		log.Errorf("File %s exists, will be overwritten", fname)
	}
	appFileHandle, err := os.Create(fname)
	if err != nil {
		return "", nil, err
	}
	defer appFileHandle.Close()

	args := []string{srcConfig["_name_"].(string), destConfig["_name_"].(string)}
	var config = make(map[string]interface{})

	// check appname as destination uri
	if !strings.Contains(destConfig["uri"].(string), "/") {
		destConfig["uri"], err = app.GetAppURL(destConfig["uri"].(string))
		if err != nil {
			return "", nil, err
		}
	}
	// check appname as source uri
	if (!strings.Contains(srcConfig["uri"].(string), "/")) && srcConfig["_name_"].(string) == "elasticsearch" {
		srcConfig["uri"], err = app.GetAppURL(srcConfig["uri"].(string))
		if err != nil {
			return "", nil, err
		}
	}
	// check file path as source [json, csv]
	if common.StringInSlice(srcConfig["_name_"].(string), []string{"json", "csv"}) {
		err = common.IsFileValid(srcConfig["uri"].(string))
		if err != nil {
			return "", nil, err
		}
	}

	configuredAdaptors := make(map[string]adaptor.Adaptor)

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
		configuredAdaptors[nodeName] = a
		// get config json
		b, err := json.Marshal(a)
		if err != nil {
			return "", nil, err
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
			return "", nil, err
		}
		appFileHandle.WriteString(string(dat))
	} else {
		// no transform file

		// set Config({log_dir})
		if srcConfig["log_dir"] != "" {
			confStr := fmt.Sprintf(`t.Config({"log_dir":"%s"}).Source("source", source, "/%s/").Save("sink", sink, "/.*/")`, srcConfig["log_dir"], srcConfig["srcRegex"])

			fmt.Println(confStr)

			appFileHandle.WriteString(confStr)
		} else {
			appFileHandle.WriteString(
				fmt.Sprintf(`t.Source("source", source, "/%s/").Save("sink", sink, "/.*/")`,
					srcConfig["srcRegex"]),
			)
		}
	}
	appFileHandle.WriteString("\n")

	return fname, configuredAdaptors, nil
}

// genPipelineFromEnv generates a pipeline file from config file
func genPipelineFromEnv(filename string) (string, map[string]adaptor.Adaptor, error) {
	var config map[string]string
	config, err := godotenv.Read(filename)
	if err != nil {
		return "", nil, err
	}
	// save keys as small
	for k := range config {
		config[strings.ToLower(k)] = config[k]
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
			// ssl should be boolean
			if k == "ssl" {
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
			if k == "tail" {
				if val == "true" {
					dest[v] = true
				} else {
					dest[v] = false
				}
			}
		}
	}
	// generate file
	file, configuredAdaptors, err := writeConfigFile(src, dest)
	if err != nil {
		return "", nil, err
	}
	fmt.Printf("Writing %s...\n", file)
	return file, configuredAdaptors, nil
}

func verifyConnections(adaptors map[string]adaptor.Adaptor) error {
	for _, ad := range adaptors {
		err := ad.Verify()
		if err != nil {
			return err
		}
	}
	return nil
}

// +build !oss

package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"

	"github.com/appbaseio/abc/imports/adaptor"
)

// GLOBALS
var srcParamMap = map[string]string{
	"src.uri":          "uri",
	"src.type":         "_name_",
	"tail":             "tail", // TODO: data type here
	"replication_slot": "replication_slot",
	"typename":         "typeName",
	"timeout":          "timeout",
}

var destParamMap = map[string]string{
	"dest.uri":  "uri",
	"dest.type": "_name_",
}

func runInit(args []string) error {
	flagset := baseFlagSet("init")
	flagset.Usage = usageFor(flagset, "abc init [source] [sink]")
	configFile := flagset.StringP("env", "e", "[FILE]", "Generate pipeline.js from an env file")

	if err := flagset.Parse(args); err != nil {
		return err
	}
	args = flagset.Args()

	// env file
	if *configFile != "" && *configFile != "[FILE]" {
		if len(args) > 0 {
			return fmt.Errorf("wrong number of arguments provided, expected 0, got %d", len(args))
		}
		_, err := genPipelineFromEnv(*configFile)
		if err != nil {
			return fmt.Errorf("There was an error %s", err)
		}
		return nil
	}

	// normal pipeline.js
	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments provided, expected 2, got %d", len(args))
	}
	if _, err := os.Stat("pipeline.js"); err == nil {
		fmt.Print("pipeline.js exists, overwrite? (y/n) ")
		var overwrite string
		fmt.Scanln(&overwrite)
		if strings.ToLower(overwrite) != "y" {
			fmt.Println("not overwriting pipeline.js, exiting...")
			return nil
		}
	}
	fmt.Println("Writing pipeline.js...")
	appFileHandle, err := os.Create(defaultPipelineFile)
	if err != nil {
		return err
	}
	defer appFileHandle.Close()
	nodeName := "source"
	for _, name := range args {
		a, _ := adaptor.GetAdaptor(name, map[string]interface{}{})
		if d, ok := a.(adaptor.Describable); ok {
			appFileHandle.WriteString(fmt.Sprintf("var %s = %s(%s)\n\n", nodeName, name, d.SampleConfig()))
			nodeName = "sink"
		} else {
			return fmt.Errorf("adaptor '%s' did not provide a sample config", name)
		}
	}
	appFileHandle.WriteString(`t.Source("source", source, "/.*/").Save("sink", sink, "/.*/")`)
	appFileHandle.WriteString("\n")
	return nil
}

func genPipelineFromEnv(filename string) (string, error) {
	var config map[string]string
	config, err := godotenv.Read(filename)
	if err != nil {
		return "", err
	}
	// source
	src := map[string]interface{}{}
	for k, v := range srcParamMap {
		if val, ok := config[k]; ok {
			src[v] = val
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
	fmt.Printf("Writing %s...", file)
	return file, nil
}

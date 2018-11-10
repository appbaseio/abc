package cluster

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1"
)

var clusterOptions = []*survey.Question{
	{
		Name:     "name",
		Prompt:   &survey.Input{Message: "Enter the name of the cluster"},
		Validate: survey.Required,
	},
	{
		Name: "location",
		Prompt: &survey.Select{
			Message: "Enter the location of the cluster",
			Help:    "The first 13 options are for azure clusters and the rest for gke",
			Options: []string{
				"eastus",
				"westeurope",
				"centralus",
				"canadacentral",
				"canadaeast",
				"australiaeast",
				"eastus2",
				"japaneast",
				"northeurope",
				"southeastasia",
				"uksouth",
				"westus2",
				"westus",
				"us-east1-b",
				"europe-west1-b",
				"us-central1-b",
				"australia-southeast1-b",
				"us-east4-b",
				"southamerica-east1-c",
				"northamerica-northeast1-b",
				"europe-north1-b",
				"asia-southeast1-b",
				"asia-east1-b",
				"asia-northeast1-a",
			},
		},
		Validate: survey.Required,
	},
	{
		Name: "vm_size",
		Prompt: &survey.Select{
			Message: "Enter the VM size of the cluster",
			Help:    "The first three options are for azure clusters and the rest for gke",
			Options: []string{
				"Standard_B2s",
				"Standard_B2ms",
				"Standard_B4ms",
				"g1-small",
				"n1-standard-1",
				"n1-standard-2",
				"n1-standard-4",
				"n1-standard-8",
			},
		},
		Validate: survey.Required,
	},
	{
		Name:     "ssh_public_key",
		Prompt:   &survey.Input{Message: "Enter the ssh public key of the cluster"},
		Validate: survey.Required,
	},
	{
		Name: "pricing_plan",
		Prompt: &survey.Select{
			Message: "Enter the pricing plan",
			Options: []string{"Growth", "Sandbox", "Hobby", "Production1", "Production2", "Production3"},
		},
		Validate: survey.Required,
	},
	{
		Name: "provider",
		Prompt: &survey.Select{
			Message: "Choose the provider:",
			Options: []string{"azure", "gke"},
		},
		Validate: survey.Required,
	},
}

func buildClusterObjectString() string {
	fmt.Println("Enter the cluster details")

	answers := make(map[string]interface{})

	err := survey.Ask(clusterOptions, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	clusterObject := "\"cluster\": {\n    "
	for key, value := range answers {
		if value != "" {
			clusterObject = clusterObject + "\"" + key + "\": " + "\"" + value.(string) + "\"," + "\n    "
		}
	}
	idx := strings.LastIndex(clusterObject, ",")
	clusterObject = clusterObject[:idx] + clusterObject[idx+1:] + "},\n"
	return clusterObject

}

var esOptions = []*survey.Question{
	{
		Name: "nodes",
		Prompt: &survey.Input{
			Message: "Enter the number of ES nodes",
			Help:    "Must be an odd number. 1 <= nodes <= 9",
		},
		Validate: survey.Required,
	},
	{
		Name: "version",
		Prompt: &survey.Input{
			Message: "Enter ES version",
			Help:    "Must be a valid Elasticsearch version in the format x.y.z",
		},
		Validate: survey.Required,
	},
	{
		Name:     "volume_size",
		Prompt:   &survey.Input{Message: "Enter the volume size"},
		Validate: survey.Required,
	},
	{
		Name: "config_url",
		Prompt: &survey.Input{
			Message: "Enter the config url",
			Help:    "Must be a URL to yaml file",
		},
	},
	{
		Name:   "heap_size",
		Prompt: &survey.Input{Message: "Enter the heap size"},
	},
	{
		Name: "plugins",
		Prompt: &survey.MultiSelect{
			Message: "Enter plugins",
			Options: []string{
				"ICU Analysis",
				"Japanese Analysis",
				"Phonetic Analysis",
				"Smart Chinese Analysis",
				"Ukrainian Analysis",
				"Stempel Polish Analysis",
				"Ingest Attachment Processor",
				"Ingest User Agent Processor",
				"Mapper Size",
				"Mapper Murmur3",
				"X-Pack",
			},
		},
	},
	{
		Name:   "backup",
		Prompt: &survey.Confirm{Message: "Do you want backup?"},
	},
	{
		Name:   "env",
		Prompt: &survey.Input{Message: "Enter env"},
	},
}

// GetESDetails ....
func buildESObjectString() string {
	fmt.Println("Enter the ES details")
	answers := make(map[string]interface{})

	err := survey.Ask(esOptions, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	esObject := "\"elasticsearch\": {\n    "

	// TODO add case for env
Next:
	for key, value := range answers {
		if value != "" {
			if key == "nodes" || key == "volume_size" {
				// Skipping integers from adding escaped quotes
			} else if key == "backup" {
				value = strconv.FormatBool(value.(bool))
			} else if key == "plugins" {
				tempStr := ""
				for _, str := range value.([]string) {
					tempStr = tempStr + "\"" + str + "\", "
				}
				value = "[" + strings.Trim(tempStr, ", ") + "]"
			} else {
				esObject = esObject + "\"" + key + "\": " + "\"" + value.(string) + "\"," + "\n    "
				continue Next
			}
			esObject = esObject + "\"" + key + "\": " + value.(string) + ",\n    "
		}
	}

	idx := strings.LastIndex(esObject, ",")
	esObject = esObject[:idx] + esObject[idx+1:] + "},\n"
	return esObject
}

var logstashKibanaOptions = []*survey.Question{
	{
		Name:   "create_node",
		Prompt: &survey.Confirm{Message: "Do you want to create node?"},
	},
	{
		Name: "version",
		Prompt: &survey.Input{
			Message: "Enter a valid ES version",
			Help:    "Must be a valid Elasticsearch version in the format x.y.z",
		},
		Validate: survey.Required,
	},
	{
		Name:   "heap_size",
		Prompt: &survey.Input{Message: "Enter the heap size"},
	},
	{
		Name:   "env",
		Prompt: &survey.Input{Message: "Enter env"},
	},
}

func buildLogstashObjectString() string {
	fmt.Println("Enter the logstash details")
	answers := make(map[string]interface{})
	err := survey.Ask(logstashKibanaOptions, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	logstashObject := "\"logstash\": {\n    "
	for key, value := range answers {
		if value != "" {
			if key == "create_node" {
				value = strconv.FormatBool(value.(bool))
				logstashObject = logstashObject + "\"" + key + "\": " + value.(string) + ",\n    "
			} else {
				logstashObject = logstashObject + "\"" + key + "\": " + "\"" + value.(string) + "\"" + ",\n    "
			}
		}
	}
	idx := strings.LastIndex(logstashObject, ",")
	logstashObject = logstashObject[:idx] + logstashObject[idx+1:] + "},\n"
	return logstashObject
}

func buildKibanaObjectString() string {
	fmt.Println("Enter the kibana details")
	answers := make(map[string]interface{})

	err := survey.Ask(logstashKibanaOptions, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	kibanaObject := "\"kibana\": {\n    "
	for key, value := range answers {
		if value != "" {
			if key == "create_node" {
				value = strconv.FormatBool(value.(bool))
				kibanaObject = kibanaObject + "\"" + key + "\": " + value.(string) + ",\n    "
			} else {
				kibanaObject = kibanaObject + "\"" + key + "\": " + "\"" + value.(string) + "\"" + ",\n    "
			}
		}
	}

	idx := strings.LastIndex(kibanaObject, ",")
	kibanaObject = kibanaObject[:idx] + kibanaObject[idx+1:] + "},\n"
	return kibanaObject
}

var addonsOptions = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Select{
			Message: "Choose an addon from the following list:",
			Help:    "In case of multiple addons be sure not to select the same addon more than once to prevent redundancy in the JSON object.",
			Options: []string{"dejavu", "elasticsearch-hq", "mirage"},
		},
		Validate: survey.Required,
	},
	{
		Name:     "image",
		Prompt:   &survey.Input{Message: "Enter image"},
		Validate: survey.Required,
	},
	{
		Name:     "exposed_port",
		Prompt:   &survey.Input{Message: "Enter the exposed port"},
		Validate: survey.Required,
	},
	{
		Name:   "env",
		Prompt: &survey.Input{Message: "Enter env"},
	},
	{
		Name:   "path",
		Prompt: &survey.Input{Message: "Enter path"},
	},
}

func buildAddonsObjectString(number int) string {
	addonsObject := "\"addons\": [\n"

	for i := 0; i < number; i++ {
		addonsObject = addonsObject + "    {\n      "
		fmt.Print("Enter the add-ons details")
		answers := make(map[string]interface{})

		err := survey.Ask(addonsOptions, &answers)
		if err != nil {
			fmt.Println(err.Error())
		}

		for key, value := range answers {
			if value != "" {
				if key == "exposed_port" {
					value = value.(string)
					addonsObject = addonsObject + "\"" + key + "\": " + value.(string) + ",\n      "
				} else {
					addonsObject = addonsObject + "\"" + key + "\": " + "\"" + value.(string) + "\",\n      "
				}
			}
		}

		ind := strings.LastIndex(addonsObject, ",")
		addonsObject = addonsObject[:ind] + addonsObject[ind+1:]
		addonsObject = addonsObject + "},\n    "
	}

	idx := strings.LastIndex(addonsObject, ",")
	addonsObject = addonsObject[:idx] + addonsObject[idx+1:] + "],\n    "
	return addonsObject
}

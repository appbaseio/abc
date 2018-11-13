package cluster

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey"
	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
)

type status struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type cluster struct {
	Name      string    `json:"name"`
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Message   string    `json:"message"`
	Provider  string    `json:"provider"`
}

type createClusterRespBody struct {
	Status  status  `json:"status"`
	Cluster cluster `json:"cluster"`
}

// BuildRequestBody creates a request body to for cluster deployment based on input from flags
func BuildRequestBody(name string, location string, vmSize string, plan string, ssh string, provider string, nodes int, esVersion string, volumeSize int) string {
	esBody := "\"elasticsearch\": {\n    \"nodes\": " + strconv.Itoa(nodes) + ",\n    \"version\": \"" + esVersion + "\",\n    \"volume_size\": " + strconv.Itoa(volumeSize) + "\n  }"
	clusterBody := "\"cluster\": {\n    \"name\": \"" + name + "\",\n    \"location\": \"" + location + "\",\n    \"vm_size\": \"" + vmSize + "\",\n    \"ssh_public_key\": \"" + ssh + "\",\n    \"pricing_plan\": \"" + plan + "\",\n    \"provider\": \"" + provider + "\"\n  }"
	return "{\n  " + esBody + ",\n  " + clusterBody + "\n}"
}

var additionalChoices = []*survey.Question{
	{
		Name:     "logstash",
		Prompt:   &survey.Confirm{Message: "Would you like to provide Logstash options to your cluster deployment?"},
		Validate: survey.Required,
	},
	{
		Name:     "kibana",
		Prompt:   &survey.Confirm{Message: "Would you like to provide Kibana options to your cluster deployment?"},
		Validate: survey.Required,
	},
	{
		Name: "addons",
		Prompt: &survey.Select{
			Message: "How many add-ons would you like to add to your cluster deployment?",
			Help:    "The following adddons are supported currently: Mirage, DejaVu, Elasticsearch-HQ",
			Options: []string{"0", "1", "2", "3"},
		},
		Validate: survey.Required,
	},
}

// BuildRequestBodyInteractive asks the user questions based on which it constructs the
// request body string to deploy the cluster.
func BuildRequestBodyInteractive() string {
	answers := make(map[string]interface{})

	err := survey.Ask(additionalChoices, &answers)
	if err != nil {
		fmt.Println(err.Error())
	}

	respBodyString := "{\n  " + buildESObjectString() + "  " + buildClusterObjectString()

	if answers["logstash"] == true {
		respBodyString = respBodyString + buildLogstashObjectString()
	}
	if answers["kibana"] == true {
		respBodyString = respBodyString + buildKibanaObjectString()
	}
	if answers["addons"] != "0" {
		num, _ := strconv.Atoi(answers["addons"].(string))
		respBodyString = respBodyString + buildAddonsObjectString(num)
	}

	idx := strings.LastIndex(respBodyString, ",")
	return respBodyString[:idx] + respBodyString[idx+1:] + "}"
}

// DeployCluster creates a cluster in interactive mode by asking the user
// for the deployment details.
func DeployCluster(body string) error {
	payload := strings.NewReader(body)
	spinner.Start()
	defer spinner.Stop()

	req, err := http.NewRequest("POST", common.AccAPIURL+"/v1/_deploy", payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := session.SendRequest(req)
	if err != nil {
		return err
	}
	spinner.Stop()

	// status code not 200
	if resp.StatusCode != 202 {
		defer resp.Body.Close()
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("There was an error %s", string(bodyBytes))
	}

	var res createClusterRespBody
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return err
	}

	// output
	fmt.Printf("ID:    %s\n", res.Cluster.ID)
	fmt.Printf("Name:  %s\n", res.Cluster.Name)
	fmt.Printf("Status:  %s\n", res.Cluster.Status)
	fmt.Printf("Provider:  %s\n", res.Cluster.Provider)
	fmt.Printf("Created at:  %s\n", res.Cluster.CreatedAt)
	fmt.Printf("Message:  %s\n", res.Cluster.Message)

	return nil
}

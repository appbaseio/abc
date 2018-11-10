package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/appbaseio/abc/appbase/common"
	"github.com/appbaseio/abc/appbase/session"
	"github.com/appbaseio/abc/appbase/spinner"
)

type clusterDetailsBody struct {
	Name              string    `json:"name"`
	ID                string    `json:"id"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	PricingPlan       string    `json:"pricing_plan"`
	Region            string    `json:"region"`
	EsVersion         string    `json:"es_version"`
	TotalNodes        int       `json:"total_nodes"`
	DashboardURL      string    `json:"dashboard_url"`
	DashboardUsername string    `json:"dashboard_username"`
	DashboardPassword string    `json:"dashboard_password"`
	DashboardHTTPS    bool      `json:"dashboard_https"`
}

type kibanaDetailsBody struct{}

type logstashDetailsBody struct{}

type esDetailsBody struct {
	Name          string `json:"name"`
	RequiredNodes int    `json:"required_nodes"`
	ReadyNodes    int    `json:"ready_nodes"`
	Status        string `json:"status"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	URL           string `json:"url"`
	HTTPS         bool   `json:"https"`
}

type deploymentDetailsBody struct {
	Elasticsearch esDetailsBody       `json:"elasticsearch"`
	Logstash      logstashDetailsBody `json:"logstash"`
	Kibana        kibanaDetailsBody   `json:"kibana"`
}

type statusDetailsBody struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type clusterRespBody struct {
	Status     statusDetailsBody     `json:"status"`
	Deployment deploymentDetailsBody `json:"deployment"`
	Cluster    clusterDetailsBody    `json:"cluster"`
}

// ShowClusterDetails shows the cluster details
func ShowClusterDetails(cluster string) error {
	spinner.StartText("Loading cluster details")
	defer spinner.Stop()

	// get cluster details
	req, err := http.NewRequest("GET", common.AccAPIURL+"/v1/_status/"+cluster, nil)
	if err != nil {
		return err
	}
	resp, err := session.SendRequest(req)
	if err != nil {
		return err
	}
	spinner.Stop()
	// decode response
	var res clusterRespBody
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&res)
	if err != nil {
		return err
	}

	if res.Status.Code != 200 {
		fmt.Println(res.Status.Code)
		return errors.New("Could not fetch cluster details")
	}

	fmt.Println("Cluster Details")
	fmt.Println("--------------------------------------------------------")
	fmt.Printf("Name:                  %s\n", res.Cluster.Name)
	fmt.Printf("ID:                    %s\n", res.Cluster.ID)
	fmt.Printf("Status:                %s\n", res.Cluster.Status)
	fmt.Printf("Created at:            %s\n", res.Cluster.CreatedAt)
	fmt.Printf("Pricing Plan:          %s\n", res.Cluster.PricingPlan)
	fmt.Printf("Region:                %s\n", res.Cluster.Region)
	fmt.Printf("ES Version:            %s\n", res.Cluster.EsVersion)
	fmt.Printf("Total Nodes:           %d\n", res.Cluster.TotalNodes)
	fmt.Printf("Dashboard URL:         %s\n", res.Cluster.DashboardURL)
	fmt.Printf("Dashboard Username:    %s\n", res.Cluster.DashboardUsername)
	fmt.Printf("Dashboard Password:    %s\n", res.Cluster.DashboardPassword)
	fmt.Printf("Dashboard HTTPS:       %s\n", strconv.FormatBool(res.Cluster.DashboardHTTPS))

	fmt.Println()

	fmt.Println("Deployment Details")
	fmt.Println("--------------------------------------------------------")
	fmt.Println("ElasticSearch:")

	fmt.Printf("ID:                %s\n", cluster)
	fmt.Printf("Name:              %s\n", res.Deployment.Elasticsearch.Name)
	fmt.Printf("Required Nodes:    %d\n", res.Deployment.Elasticsearch.RequiredNodes)
	fmt.Printf("Ready Nodes:       %d\n", res.Deployment.Elasticsearch.ReadyNodes)
	fmt.Printf("Status:            %s\n", res.Deployment.Elasticsearch.Status)
	fmt.Printf("Username:          %s\n", res.Deployment.Elasticsearch.Username)
	fmt.Printf("Password:          %s\n", res.Deployment.Elasticsearch.Password)
	fmt.Printf("URL:               %s\n", res.Deployment.Elasticsearch.URL)
	fmt.Printf("HTTPS:             %s\n", strconv.FormatBool(res.Deployment.Elasticsearch.HTTPS))

	return nil
}

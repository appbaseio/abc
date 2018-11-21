package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/appbaseio/abc/appbase/common"
	"github.com/olekukonko/tablewriter"
)

type analyticsResults struct {
	Count json.Number `json:"count"`
	Key   string      `json:"key"`
}

type analyticsVolumeResults struct {
	Count     json.Number `json:"count"`
	Key       json.Number `json:"key"`
	DateAsStr string      `json:"key_as_string"`
}

type overviewAnalyticsBody struct {
	NoResultSearches []analyticsResults       `json:"noResultSearches"`
	PopularSearches  []analyticsResults       `json:"popularSearches"`
	SearchVolume     []analyticsVolumeResults `json:"searchVolume"`
}

//ShowOverview .......
func ShowOverview(body io.ReadCloser) error {

	var res overviewAnalyticsBody
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Display the overview results
	fmt.Println("Analytics(Overview) Results:")

	// Display NoResultSearches results
	noResultTable := tablewriter.NewWriter(os.Stdout)
	noResultTable.SetHeader([]string{"Count", "Key"})

	for _, elements := range res.NoResultSearches {
		noResultTable.Append([]string{common.JSONNumberToString(elements.Count), elements.Key})
	}
	noResultTable.SetAlignment(tablewriter.ALIGN_CENTER)
	fmt.Println("No Result Searches")
	noResultTable.Render()

	// Display PopularSearches results
	popularSearchesTable := tablewriter.NewWriter(os.Stdout)
	popularSearchesTable.SetHeader([]string{"Count", "Key"})

	for _, elements := range res.PopularSearches {
		popularSearchesTable.Append([]string{common.JSONNumberToString(elements.Count), elements.Key})
	}
	popularSearchesTable.SetAlignment(tablewriter.ALIGN_CENTER)
	fmt.Println("No Result Searches")
	popularSearchesTable.Render()

	// Display SearcheVolume results
	searchVolumeTable := tablewriter.NewWriter(os.Stdout)
	searchVolumeTable.SetHeader([]string{"Count", "Key", "Date-As-Str"})
	for _, elements := range res.SearchVolume {
		searchVolumeTable.Append([]string{common.JSONNumberToString(elements.Count), common.JSONNumberToString(elements.Key), elements.DateAsStr})
	}
	searchVolumeTable.SetAlignment(tablewriter.ALIGN_CENTER)
	fmt.Println("Search Volume Results")
	searchVolumeTable.Render()

	return nil
}

type latencyResults struct {
	Count json.Number `json:"count"`
	Key   json.Number `json:"key"`
}

type latency struct {
	Latency []latencyResults `json:"latency"`
}

//ShowLatency .......
func ShowLatency(body io.ReadCloser) error {
	var res latency
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)

	if err != nil {
		fmt.Println(err)
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Count", "Key"})

	for _, elements := range res.Latency {
		table.Append([]string{common.JSONNumberToString(elements.Count), common.JSONNumberToString(elements.Key)})
	}
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	fmt.Println("Analytics(Latency) Results:")
	table.Render()

	return nil
}

type geoIP struct {
	GeoIP []analyticsResults `json:"aggrByCountry"`
}

//ShowGeoIP .......
func ShowGeoIP(body io.ReadCloser) error {
	var res geoIP
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)
	if err != nil {
		fmt.Println(err)
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Count", "Key"})

	for _, elements := range res.GeoIP {
		table.Append([]string{common.JSONNumberToString(elements.Count), elements.Key})
	}
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	fmt.Println("Analytics(GeoIP) Results:")
	table.Render()

	return nil
}

type analyticsPopularResults struct {
	Count  json.Number `json:"count"`
	Key    string      `json:"key"`
	Source string      `json:"source"`
}

type popularResults struct {
	PopularResults []analyticsPopularResults `json:"popularResults"`
}

//ShowPopularResults .......
func ShowPopularResults(body io.ReadCloser) error {
	var res popularResults
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// TODO refine output

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Count", "Key", "Source"})

	for _, elements := range res.PopularResults {
		table.Append([]string{common.JSONNumberToString(elements.Count), elements.Key, elements.Source})
	}
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	fmt.Println("Analytics(Popular) Results:")
	table.Render()

	return nil
}

type popularSearches struct {
	Results []analyticsResults `json:"popularSearches"`
}

//ShowPopularSearches .......
func ShowPopularSearches(body io.ReadCloser) error {
	var res popularSearches
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)
	if err != nil {
		fmt.Println(err)
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Count", "Key"})

	for _, elements := range res.Results {
		table.Append([]string{common.JSONNumberToString(elements.Count), elements.Key})
	}
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	fmt.Println("Analytics Popular Searches:")
	table.Render()

	return nil
}

type noResultSearches struct {
	Results []analyticsResults `json:"noResultSearches"`
}

//ShowNoResultSearches .......
func ShowNoResultSearches(body io.ReadCloser) error {
	var res noResultSearches
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)
	if err != nil {
		fmt.Println(err)
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Count", "Key"})

	for _, elements := range res.Results {
		table.Append([]string{common.JSONNumberToString(elements.Count), elements.Key})
	}
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	fmt.Println("Analytics No Result Searches:")
	table.Render()

	return nil
}

type analyticsPopularFilters struct {
	Count json.Number `json:"count"`
	Key   string      `json:"key"`
	Value string      `json:"value"`
}

type popularFilters struct {
	Results []analyticsPopularFilters `json:"popularFilters"`
}

//ShowPopularFilters .......
func ShowPopularFilters(body io.ReadCloser) error {
	var res popularFilters
	dec := json.NewDecoder(body)
	err := dec.Decode(&res)
	if err != nil {
		fmt.Println(err)
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Count", "Key", "Value"})

	for _, elements := range res.Results {
		table.Append([]string{common.JSONNumberToString(elements.Count), elements.Key, elements.Value})
	}
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	fmt.Println("Analytics Popular Filters:")
	table.Render()

	return nil
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
)

type servers struct {
	Server []server
}

type server struct {
	Key                 int
	Name                string
	Description         []string
	CPU                 string
	CPUBenchmark        int  `json:"cpu_benchmark"`
	CPUCount            int  `json:"cpu_count"`
	IsHighIO            bool `json:"is_highio"`
	IsECC               bool `json:"is_ecc"`
	Traffic             string
	Dist                []string
	Bandwitch           int
	RAM                 int
	Price               string
	PriceV              string `json:"price_v"`
	RAMHR               string `json:"ram_hr"`
	SetupPrice          string `json:"setup_price"`
	HDDSize             int    `json:"hdd_size"`
	HDDCount            int    `json:"hdd_count"`
	HDDHR               string `json:"hdd_hr"`
	FixedPrice          bool   `json:"fixed_price"`
	NextReduce          int    `json:"next_reduce"`
	NextReduceHR        string `json:"next_reduce_hr"`
	NextReduceTimestamp int    `json:"next_reduce_timestamp"`
	Datacenter          []string
	Specials            []string
	SpecialHDD          string
	FreeText            string
}

type slackRequestBody struct {
	Text string `json:"text"`
}

func handleError(err error, fatal bool) {
	if err != nil {
		if !fatal {
			log.Println(err)
		} else {
			log.Fatalln(err)
		}
	}
}

func parseData(url string) servers {
	c := &http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	handleError(err, true)

	res, err := c.Do(req)
	handleError(err, true)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	handleError(err, true)

	var servers servers
	err = json.Unmarshal(body, &servers)
	handleError(err, true)

	return servers
}

func filterData(s servers, p float64, h int) servers {
	var filteredServers servers
	for _, v := range s.Server {
		price, err := strconv.ParseFloat(v.Price, 64)
		handleError(err, true)

		if price < p {
			if (v.HDDCount*v.HDDSize) > h && strings.Contains(v.HDDHR, "TB") {
				filteredServers.Server = append(filteredServers.Server, v)
			}
		}
	}

	return filteredServers
}

func sendSlackMessage(s servers) {
	webhook := os.Getenv("SLACK_WEBHOOK")

	slackBody, _ := json.Marshal(slackRequestBody{Text: formatString(s)})

	req, err := http.NewRequest(http.MethodPost, webhook, bytes.NewBuffer(slackBody))
	handleError(err, true)

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	_, err = client.Do(req)
	handleError(err, true)
}

func formatString(s servers) string {
	var builder strings.Builder
	writer := tabwriter.NewWriter(&builder, 15, 4, 10, ' ', 0)

	fmt.Fprintf(&builder, "I found you the following %d servers that may be of interest:\n", len(s.Server))
	fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t\n", "Key", "RAM(GB)", "HDD", "CPU", "Price(€)", "Description")
	for _, v := range s.Server {
		fmt.Fprintf(writer, "%d\t%d\t%s\t%s\t%s\t%s\t\n", v.Key, v.RAM, v.HDDHR, v.CPU, v.Price, v.FreeText)
	}

	writer.Flush()
	return builder.String()
}

func main() {
	lambda.Start(handleRequest)
}

func handleRequest() {
	url := "https://www.hetzner.com/a_hz_serverboerse/live_data.json"

	servers := parseData(url)

	filteredServers := filterData(servers, 25, 2)
	sendSlackMessage(filteredServers)
}

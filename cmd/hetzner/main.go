package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/robpickerill/hetzner-auction-scraper/internal/hetzner"
	"github.com/robpickerill/hetzner-auction-scraper/internal/slack"
)

func handleError(err error, fatal bool) {
	if err != nil {
		if !fatal {
			log.Println(err)
		} else {
			log.Fatalln(err)
		}
	}
}

func filterData(s *hetzner.Servers, p float64, h int) *hetzner.Servers {
	var filteredServers hetzner.Servers
	for _, v := range s.Servers {
		price, err := strconv.ParseFloat(v.Price, 64)
		handleError(err, true)

		if price < p {
			if (v.HDDCount*v.HDDSize) > h && strings.Contains(v.HDDHR, "TB") {
				filteredServers.Servers = append(filteredServers.Servers, v)
			}
		}
	}

	return &filteredServers
}

func formatString(s *hetzner.Servers) string {
	var builder strings.Builder

	writer := tabwriter.NewWriter(&builder, 15, 4, 10, ' ', 0)

	fmt.Fprintf(&builder, ":wave: I found you the following %d servers that may be of interest:\n\n", len(s.Servers))
	fmt.Fprint(&builder, "```")
	fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t\n", "Key", "RAM(GB)", "HDD", "CPU", "Price(€)", "Description")
	fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\t%s\t\n", "---", "-------", "---", "---", "--------", "-----------")
	for _, v := range s.Servers {
		fmt.Fprintf(writer, "%d\t%d\t%s\t%s\t%s\t%s\t\n", v.Key, v.RAM, v.HDDHR, v.CPU, v.Price, v.FreeText)
	}
	writer.Flush()
	fmt.Fprint(&builder, "```")

	writer.Flush()
	return builder.String()
}

func handleRequest() {
	servers, err := hetzner.ServerAuctionResults()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Found %d servers at the auction", len(servers.Servers))

	filteredServers := filterData(servers, 25, 2)
	log.Printf("Found %d servers of interest", len(filteredServers.Servers))

	if len(filteredServers.Servers) > 0 {
		log.Printf("Sending notification to slack")

		message := slack.CreateMessage(formatString(filteredServers))

		webhook := os.Getenv("SLACK_WEBHOOK")
		err := slack.SendMessage(message, webhook)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	lambda.Start(handleRequest)
}

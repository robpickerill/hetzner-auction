package hetzner

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Servers models the servers key from the hetzner server auction api
type Servers struct {
	Servers []Server `json:"server"`
}

// Server models the server object from the hetzner server auction api
type Server struct {
	Key                 int      `json:"key"`
	Name                string   `json:"name"`
	Description         []string `json:"description"`
	CPU                 string   `json:"cpu"`
	CPUBenchmark        int      `json:"cpu_benchmark"`
	CPUCount            int      `json:"cpu_count"`
	IsHighIO            bool     `json:"is_highio"`
	IsECC               bool     `json:"is_ecc"`
	Traffic             string   `json:"traffic"`
	Dist                []string `json:"dist"`
	Bandwidth           int      `json:"bandwidth"`
	RAM                 int      `json:"ram"`
	Price               string   `json:"price"`
	PriceV              string   `json:"price_v"`
	RAMHR               string   `json:"ram_hr"`
	SetupPrice          string   `json:"setup_price"`
	HDDSize             int      `json:"hdd_size"`
	HDDCount            int      `json:"hdd_count"`
	HDDHR               string   `json:"hdd_hr"`
	FixedPrice          bool     `json:"fixed_price"`
	NextReduce          int      `json:"next_reduce"`
	NextReduceHR        string   `json:"next_reduce_hr"`
	NextReduceTimestamp int      `json:"next_reduce_timestamp"`
	Datacenter          []string `json:"datacenter"`
	Specials            []string `json:"specials"`
	SpecialHDD          string   `json:"specialHdd"`
	FreeText            string   `json:"freetext"`
}

// AuctionURL is the location of the json data for the hetzner server auction
const AuctionURL string = "https://www.hetzner.com/a_hz_serverboerse/live_data.json"

// ServerAuctionResults returns the results from the hetzner server auction as Servers{}
func ServerAuctionResults() (*Servers, error) {
	c := &http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, AuctionURL, nil)
	if err != nil {
		return &Servers{}, err
	}

	res, err := c.Do(req)
	if err != nil {
		return &Servers{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &Servers{}, err
	}

	var servers Servers
	err = json.Unmarshal(body, &servers)
	if err != nil {
		return &Servers{}, err
	}

	return &servers, nil
}

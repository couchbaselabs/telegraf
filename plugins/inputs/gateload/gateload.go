package gateload

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const kDefaultExpVarHost = "http://localhost:9876"

const kExpVarEndpoint = "debug/vars"

type Gateload struct {
	Server string
}

func (s *Gateload) Description() string {
	return "Read Gateload metrics through ExpVars"
}

var sampleConfig = `
  ## specify gateload server via a url:
  ##  [protocol://address:port]
  ##  e.g.
  ##    http://localhost:9876/
  ##
  ## If no server is specified, then localhost is used as the host.
  server = ["http://localhost:9876"]
`

func (s *Gateload) SampleConfig() string {
	return sampleConfig
}

func (s *Gateload) Gather(acc telegraf.Accumulator) error {

	var expvarServer string

	if s.Server == "" {
		expvarServer = kDefaultExpVarHost
	} else {
		expvarServer = s.Server
	}

	expvarUrl := strings.Join([]string{expvarServer, kExpVarEndpoint}, "/")

	if err := FetchExpvar(expvarUrl, acc); err != nil {
		log.Println("Error: retriving _expvar: %v", err)
		return err
	}
	return nil
}

// FetchExpvar fetches expvar by http for the given addr (host:port)
func FetchExpvar(fetchurl string, acc telegraf.Accumulator) error {

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	req, _ := http.NewRequest("GET", fetchurl, nil)

	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error: executing request: %v", err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		log.Println("Error: response StatusNotFound")
		return fmt.Errorf("ExpVars not found at provided URL %v", fetchurl)
	}

	expvar, err := ParseExpvar(resp.Body)
	if err != nil {
		log.Println("Error: parsing expvar: %v", err)
		return err
	}

	tags := map[string]string{"expvar_endpoint": fetchurl}
	acc.AddFields("gateload_expvar", expvar, tags)

	return nil
}

func init() {
	inputs.Add("gateload", func() telegraf.Input { return &Gateload{} })
}

func ParseExpvar(r io.ReadCloser) (stats map[string]interface{}, err error) {
	var expvar map[string]interface{}

	dec := json.NewDecoder(r)
	dec.UseNumber()

	if err := dec.Decode(&expvar); err != nil {
		return nil, fmt.Errorf("Error: ExpVars could not be unmarshalled: %v", err)
	}

	glstats, ok := expvar["gateload"].(map[string]interface{});

	if !ok {
		return nil, fmt.Errorf("Error: ExpVar gateload not a JSON Object")
	}

	stats, err = Flatten(glstats, "", DotStyle)

	if err != nil {
		return nil, fmt.Errorf("Error: ExpVar gateload not a nested JSON Object")
	}

	//Custom convert any json.Number values to int64 or float64
	for k, v := range stats {
		switch t := v.(type) {
		case json.Number:
			var i int64
			var err error
			i, err = t.Int64()
			if err == nil {
				stats[k] = i
			} else {
				var f float64
				f, err := t.Float64()
				if err == nil {
					stats[k] = f
				}
			}
		}
	}

	return stats, nil
}

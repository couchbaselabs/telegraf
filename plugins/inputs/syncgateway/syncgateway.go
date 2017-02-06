package syncgateway

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

const kDefaultExpVarHost = "http://localhost:4985"

const kExpVarEndpoint = "_expvar"

type SyncGateway struct {
	Server string
}

func (s *SyncGateway) Description() string {
	return "Read Sync Gateway metrics through ExpVars"
}

var sampleConfig = `
  ## specify sync gateway server via a url:
  ##  [protocol://address:port]
  ##  e.g.
  ##    http://localhost:4985/
  ##
  ## If no server is specified, then localhost is used as the host.
  server = ["http://localhost:4985"]
`

func (s *SyncGateway) SampleConfig() string {
	return sampleConfig
}

func (s *SyncGateway) Gather(acc telegraf.Accumulator) error {

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
		Timeout: 1 * time.Second, // TODO: make it configurable or left default?
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
	acc.AddFields("syncgateway_expvar", expvar, tags)

	return nil
}

func init() {
	inputs.Add("syncgateway", func() telegraf.Input { return &SyncGateway{} })
}

// ParseExpvar parses expvar data from reader.
/*
func ParseExpvar(r io.ReadCloser) (map[string]interface{}, error) {
	var expvar map[string]interface{}
	if bytes, err := ioutil.ReadAll(r); err == nil {
		if err := json.Unmarshal(bytes, &expvar); err != nil {
			return nil, fmt.Errorf("Error: ExpVars could not be unmarshalled: %v",err)
		}
	} else {
		return nil, fmt.Errorf("Error reading Sync Gateway _expvar response")
	}

	return expvar["syncGateway_stats"].(map[string]interface {}), nil
}
*/

func ParseExpvar(r io.ReadCloser) ( stats map[string]interface{}, err error) {
	var expvar map[string]interface{}

	dec := json.NewDecoder(r)
	dec.UseNumber()

	if err := dec.Decode(&expvar); err != nil {
		return nil, fmt.Errorf("Error: ExpVars could not be unmarshalled: %v", err)
	}

	if sgstats, ok := expvar["syncGateway_stats"].(map[string]interface{}); ok {
		stats = sgstats
	} else {
		return nil, fmt.Errorf("Error: ExpVar syncGateway_stats not a JSON Object")
	}

	//publish memstats.Alloc and memstats.Sys
	if memstats, ok := expvar["memstats"].(map[string]interface{}); ok {

		// Gc stats
		stats["memstatsPauseTotalNs"] = memstats["PauseTotalNs"];
		stats["memstatsNumGC"] = memstats["NumGC"];

		// Mem stats
		stats["memstatsAlloc"] = memstats["Alloc"];
		stats["memstatsSys"] = memstats["Sys"];
		stats["memstatsHeapIdle"] = memstats["HeapIdle"];
		stats["memstatsHeapInuse"] = memstats["HeapInuse"];
		stats["memstatsStackInuse"] = memstats["StackInuse"];
		stats["memstatsTotalAlloc"] = memstats["TotalAlloc"];
		stats["memstatsMSpanInuse"] = memstats["MSpanInuse"];
		

	} else {
		return nil, fmt.Errorf("Error: ExpVar memstats not a JSON Object")
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


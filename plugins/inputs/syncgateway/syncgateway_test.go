package syncgateway

import (
	"testing"
	"encoding/json"
	"log"
)

func TestExtractGoCouchbaseStats(t *testing.T) {

	jsonString := `
{
  "cb": {
    "ops": {
      "GetsRaw": {
        "count": 231,
        "sum": 169553440,
        "min": 94657,
        "max": 45945639,
        "mean": 733997.5757575758,
        "p25": 177872,
        "p50": 319769,
        "p75": 458992,
        "p90": 673833,
        "p99": 17837260
      },
      "Incr": {
        "count": 255389,
        "sum": 104447181401,
        "min": 46876,
        "max": 1509687621,
        "mean": 408972.90564981266,
        "p25": 99066,
        "p50": 127986,
        "p75": 295027,
        "p90": 884070,
        "p99": 43443195
      },
      "Write(raw)": {
        "count": 349551,
        "sum": 268038422429,
        "min": 76960,
        "max": 441616912,
        "mean": 766807.768906397,
        "p25": 185244,
        "p50": 289006,
        "p75": 645711,
        "p90": 2109329,
        "p99": 44531814
      },
      "casNext": {
        "count": 950014,
        "sum": 1294461250340,
        "min": 50514,
        "max": 2199849812,
        "mean": 1362570.7098421708,
        "p25": 119554,
        "p50": 150680,
        "p75": 353550,
        "p90": 4160701,
        "p99": 1207338334
      }
    },
    "pools": {
      "127.0.0.1:11210": {
        "count": 1555186,
        "sum": 582306655880,
        "min": -14,
        "max": 1543060294,
        "mean": 374428.94668547687,
        "p25": 664,
        "p50": 886,
        "p75": 1124,
        "p90": 2317,
        "p99": 791816816
      }
    }
  }
}
	`

	var expvars interface{}
	err := json.Unmarshal([]byte(jsonString), &expvars)
	if err != nil {
		t.Fatalf("Failed to unmarshal json")
	}

	expvarsMap, ok := expvars.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map")
	}

	stats := map[string]interface{}{}

	extractGoCouchbaseStats(expvarsMap, stats)

	log.Printf("stats: %+v", stats)

	goCouchbasePoolsP90, ok := stats["goCouchbasePoolsP90"]
	if !ok {
		t.Fatalf("Expected goCouchbasePoolsP90 stat")
	}
	goCouchbasePoolsP90Float64, ok := goCouchbasePoolsP90.(float64)
	if !ok {
		t.Fatalf("Expected goCouchbasePoolsP90 to be a float64, was: %T", goCouchbasePoolsP90)
	}

	if goCouchbasePoolsP90Float64 != 2317 {
		t.Fatalf("Did not get expected goCouchbasePoolsP90 value")
	}


}
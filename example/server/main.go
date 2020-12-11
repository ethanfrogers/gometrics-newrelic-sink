package main

import (
	"net/http"
	"os"

	nrgometrics "github.com/ethanfrogers/gometrics-newreilc-sink"

	"github.com/armon/go-metrics"

	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
)

func main() {
	harvester, err := telemetry.NewHarvester(func(c *telemetry.Config) {
		c.APIKey = os.Getenv("NEW_RELIC_INSERT_API_KEY")
		c.CommonAttributes = map[string]interface{}{
			"common": "attribute",
		}
	})
	if err != nil {
		panic(err)
	}
	sink, _ := nrgometrics.NewNewRelicSink(harvester)
	registry, err := metrics.New(nil, sink)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("world"))
	})

	registry.IncrCounterWithLabels([]string{"my.metric"}, 1, []metrics.Label{{Name: "label", Value: "value"}})
}

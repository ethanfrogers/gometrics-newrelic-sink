# gometrics-newrelic-sink

This library implements a metrics sink which is compatible with `armon/go-metrics` for New Relic. It utilizes New Relic's `newrelic-telemetry-sdk-go` library
under the hood.

## Notice

This library is in development and not intended for production use. 

## Usage

```go
// create a new harvester. see https://github.com/newrelic/newrelic-telemtry-sdk-go
// for more information about configuring the harvester.
harvester, err := telemetry.NewHarvester(func(c *telemetry.Config) {
    c.APIKey = os.Getenv("NEW_RELIC_INSERT_API_KEY")
})
if err != nil {
    panic(err)
}

// crete a new sink using the above harvester
sink, _ := nrsink.NewNewRelicSink(harvester)
registry, err := metrics.New(nil, sink)
if err != nil {
    panic(err)
}

// define a counter and increment it
registry.IncrCounterWithLabels(
	[]string{"my.metric"}, 1, 
	[]metrics.Label{{Name: "label", Value: "value"}})
```

## TODO
- [ ] Unit tests
- [ ] Test in an actual service
- [ ] Add ability to pre-define gauges/counters/summaries to reduce allocations (similar to the prom sink)

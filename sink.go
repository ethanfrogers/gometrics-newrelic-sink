package gometrics_newrelic_sink

import (
	"fmt"
	"strings"
	"sync"

	"github.com/armon/go-metrics"
	"github.com/newrelic/newrelic-telemetry-sdk-go/telemetry"
)

var DefaultSeparator string = "."

type NewRelicSink struct {
	harvester       *telemetry.Harvester
	metricSeparator string

	counters  sync.Map
	gauges    sync.Map
	summaries sync.Map
}

func NewNewRelicSink(h *telemetry.Harvester) (*NewRelicSink, error) {
	return &NewRelicSink{
		harvester:       h,
		metricSeparator: DefaultSeparator,
		counters:        sync.Map{},
		gauges:          sync.Map{},
		summaries:       sync.Map{},
	}, nil
}

func (n *NewRelicSink) SetGauge(key []string, val float32) {
	n.AddSampleWithLabels(key, val, nil)
}

func (n *NewRelicSink) SetGaugeWithLabels(key []string, val float32, labels []metrics.Label) {
	_, hash := flattenKey(key, labels)
	if gauge, ok := n.counters.Load(hash); ok {
		localGauge := *gauge.(*telemetry.AggregatedGauge)
		localGauge.Value(float64(val))
		n.counters.Store(hash, &localGauge)
		return
	}

	m := n.harvester.MetricAggregator().Gauge(
		strings.Join(key, n.metricSeparator), convertLabels(labels))
	m.Value(float64(val))
	n.counters.Store(hash, &m)
}

func (n *NewRelicSink) EmitKey(key []string, val float32) {}

func (n *NewRelicSink) IncrCounter(key []string, val float32) {
	n.IncrCounterWithLabels(key, val, nil)
}

func (n *NewRelicSink) IncrCounterWithLabels(key []string, val float32, labels []metrics.Label) {
	_, hash := flattenKey(key, labels)
	if counter, ok := n.counters.Load(hash); ok {
		localCounter := *counter.(*telemetry.AggregatedCount)
		localCounter.Increment()
		n.counters.Store(hash, &localCounter)
		return
	}

	m := n.harvester.MetricAggregator().Count(
		strings.Join(key, n.metricSeparator), convertLabels(labels))
	m.Increment()
	n.counters.Store(hash, &m)
}

func (n *NewRelicSink) AddSample(key []string, val float32) {
	n.AddSampleWithLabels(key, val, nil)
}

func (n *NewRelicSink) AddSampleWithLabels(key []string, val float32, labels []metrics.Label) {
	_, hash := flattenKey(key, labels)
	if counter, ok := n.counters.Load(hash); ok {
		localCounter := *counter.(*telemetry.AggregatedSummary)
		localCounter.Record(float64(val))
		n.counters.Store(hash, &localCounter)
		return
	}

	m := n.harvester.MetricAggregator().Summary(
		strings.Join(key, n.metricSeparator), convertLabels(labels))
	m.Record(float64(val))
	n.counters.Store(hash, &m)
}

func flattenKey(parts []string, labels []metrics.Label) (string, string) {
	key := strings.Join(parts, "_")

	hash := key
	for _, label := range labels {
		hash += fmt.Sprintf(";%s=%s", label.Name, label.Value)
	}

	return key, hash
}

func convertLabels(labels []metrics.Label) map[string]interface{} {
	var attributeMap map[string]interface{}
	for _, l := range labels {
		attributeMap[l.Name] = l.Value
	}
	return attributeMap
}

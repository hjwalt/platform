package metric

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestGetPrometheusSingleton(t *testing.T) {
	assert := assert.New(t)

	first := GetPrometheus()
	assert.NotNil(first)

	second := GetPrometheus()
	assert.NotNil(second)
	assert.Same(first, second)

	// internal collectors are initialized (same package)
	assert.NotNil(first.retry)
	assert.NotNil(first.messagesProduced)
	assert.NotNil(first.messagesProcessed)
}

func TestPrometheusAccessorsReturnSingleton(t *testing.T) {
	assert := assert.New(t)

	singleton := GetPrometheus()

	assert.Same(singleton, PrometheusRetry())
	assert.Same(singleton, PrometheusProduce())
	assert.Same(singleton, PrometheusConsume())
}

func TestPrometheusNilReceiverSafe(t *testing.T) {
	assert := assert.New(t)

	var p *Prometheus

	assert.NotPanics(func() { p.RetryCount("topicA", 1, 3) })
	assert.NotPanics(func() { p.MessagesProducedIncrement("topicA", 1, 3) })
	assert.NotPanics(func() { p.MessagesProcessedIncrement("topicA", 1, 3) })
}

func TestPrometheusRecordsMetrics(t *testing.T) {
	assert := assert.New(t)

	p := GetPrometheus()

	// counters accumulate in the process-wide default registry, so capture
	// baselines before mutating (keeps the test valid across repeated
	// executions in one process, e.g. -count=2)
	baselineFamilies, err := prometheus.DefaultGatherer.Gather()
	assert.NoError(err)
	baselineProduced, _ := metricValue(baselineFamilies, "flows_messages_produced", map[string]string{"topic": "topicA", "partition": "3"})
	baselineProcessed, _ := metricValue(baselineFamilies, "flows_messages_received", map[string]string{"topic": "topicB", "partition": "0"})

	p.RetryCount("topicA", 3, 7)
	p.RetryCount("topicA", 3, 9) // gauge semantics: last Set wins
	p.RetryCount("topicC", 1, 2)
	p.MessagesProducedIncrement("topicA", 3, 5)
	p.MessagesProducedIncrement("topicA", 3, 2) // counter semantics: accumulates
	p.MessagesProcessedIncrement("topicB", 0, 4)

	families, err := prometheus.DefaultGatherer.Gather()
	assert.NoError(err)
	assert.NotEmpty(families)

	retry, found := metricValue(families, "flows_retry_count", map[string]string{"topic": "topicA", "partition": "3"})
	assert.True(found, "expected flows_retry_count metric for topicA partition 3")
	assert.Equal(9.0, retry)

	otherRetry, found := metricValue(families, "flows_retry_count", map[string]string{"topic": "topicC", "partition": "1"})
	assert.True(found, "expected flows_retry_count metric for topicC partition 1")
	assert.Equal(2.0, otherRetry)

	produced, found := metricValue(families, "flows_messages_produced", map[string]string{"topic": "topicA", "partition": "3"})
	assert.True(found, "expected flows_messages_produced metric for topicA partition 3")
	assert.Equal(baselineProduced+7.0, produced)

	processed, found := metricValue(families, "flows_messages_received", map[string]string{"topic": "topicB", "partition": "0"})
	assert.True(found, "expected flows_messages_received metric for topicB partition 0")
	assert.Equal(baselineProcessed+4.0, processed)
}

// metricValue scans gathered families for the metric whose family name and
// labels match, returning its gauge or counter value.
func metricValue(families []*dto.MetricFamily, familyName string, wantLabels map[string]string) (float64, bool) {
	for _, family := range families {
		if family.GetName() != familyName {
			continue
		}
		for _, metric := range family.GetMetric() {
			if len(metric.GetLabel()) != len(wantLabels) {
				continue
			}
			matches := true
			for _, labelPair := range metric.GetLabel() {
				want, ok := wantLabels[labelPair.GetName()]
				if !ok || want != labelPair.GetValue() {
					matches = false
					break
				}
			}
			if !matches {
				continue
			}
			if gauge := metric.GetGauge(); gauge != nil {
				return gauge.GetValue(), true
			}
			if counter := metric.GetCounter(); counter != nil {
				return counter.GetValue(), true
			}
		}
	}
	return 0, false
}

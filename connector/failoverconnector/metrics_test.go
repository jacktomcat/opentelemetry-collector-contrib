// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package failoverconnector // import "github.com/jacktomcat/opentelemetry-collector-contrib/connector/failoverconnector"
import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/connector"
	"go.opentelemetry.io/collector/connector/connectortest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

var errMetricsConsumer = errors.New("Error from ConsumeMetrics")

func TestMetricsRegisterConsumers(t *testing.T) {
	var sinkFirst, sinkSecond, sinkThird consumertest.MetricsSink
	metricsFirst := component.NewIDWithName(component.DataTypeMetrics, "metrics/first")
	metricsSecond := component.NewIDWithName(component.DataTypeMetrics, "metrics/second")
	metricsThird := component.NewIDWithName(component.DataTypeMetrics, "metrics/third")

	cfg := &Config{
		PipelinePriority: [][]component.ID{{metricsFirst}, {metricsSecond}, {metricsThird}},
		RetryInterval:    50 * time.Millisecond,
		RetryGap:         10 * time.Millisecond,
		MaxRetries:       10000,
	}

	router := connector.NewMetricsRouter(map[component.ID]consumer.Metrics{
		metricsFirst:  &sinkFirst,
		metricsSecond: &sinkSecond,
		metricsThird:  &sinkThird,
	})

	conn, err := NewFactory().CreateMetricsToMetrics(context.Background(),
		connectortest.NewNopCreateSettings(), cfg, router.(consumer.Metrics))

	failoverConnector := conn.(*metricsFailover)
	defer func() {
		assert.NoError(t, failoverConnector.Shutdown(context.Background()))
	}()

	require.NoError(t, err)
	require.NotNil(t, conn)

	mc, _, ok := failoverConnector.failover.getCurrentConsumer()
	mc1 := failoverConnector.failover.GetConsumerAtIndex(1)
	mc2 := failoverConnector.failover.GetConsumerAtIndex(2)

	assert.True(t, ok)
	require.Implements(t, (*consumer.Metrics)(nil), mc)
	require.Implements(t, (*consumer.Metrics)(nil), mc1)
	require.Implements(t, (*consumer.Metrics)(nil), mc2)
}

func TestMetricsWithValidFailover(t *testing.T) {
	var sinkFirst, sinkSecond, sinkThird consumertest.MetricsSink
	metricsFirst := component.NewIDWithName(component.DataTypeMetrics, "metrics/first")
	metricsSecond := component.NewIDWithName(component.DataTypeMetrics, "metrics/second")
	metricsThird := component.NewIDWithName(component.DataTypeMetrics, "metrics/third")

	cfg := &Config{
		PipelinePriority: [][]component.ID{{metricsFirst}, {metricsSecond}, {metricsThird}},
		RetryInterval:    50 * time.Millisecond,
		RetryGap:         10 * time.Millisecond,
		MaxRetries:       10000,
	}

	router := connector.NewMetricsRouter(map[component.ID]consumer.Metrics{
		metricsFirst:  &sinkFirst,
		metricsSecond: &sinkSecond,
		metricsThird:  &sinkThird,
	})

	conn, err := NewFactory().CreateMetricsToMetrics(context.Background(),
		connectortest.NewNopCreateSettings(), cfg, router.(consumer.Metrics))

	require.NoError(t, err)

	failoverConnector := conn.(*metricsFailover)
	failoverConnector.failover.ModifyConsumerAtIndex(0, consumertest.NewErr(errMetricsConsumer))
	defer func() {
		assert.NoError(t, failoverConnector.Shutdown(context.Background()))
	}()

	md := sampleMetric()

	require.Eventually(t, func() bool {
		return consumeMetricsAndCheckStable(failoverConnector, 1, md)
	}, 3*time.Second, 5*time.Millisecond)
}

func TestMetricsWithFailoverError(t *testing.T) {
	var sinkFirst, sinkSecond, sinkThird consumertest.MetricsSink
	metricsFirst := component.NewIDWithName(component.DataTypeMetrics, "metrics/first")
	metricsSecond := component.NewIDWithName(component.DataTypeMetrics, "metrics/second")
	metricsThird := component.NewIDWithName(component.DataTypeMetrics, "metrics/third")

	cfg := &Config{
		PipelinePriority: [][]component.ID{{metricsFirst}, {metricsSecond}, {metricsThird}},
		RetryInterval:    50 * time.Millisecond,
		RetryGap:         10 * time.Millisecond,
		MaxRetries:       10000,
	}

	router := connector.NewMetricsRouter(map[component.ID]consumer.Metrics{
		metricsFirst:  &sinkFirst,
		metricsSecond: &sinkSecond,
		metricsThird:  &sinkThird,
	})

	conn, err := NewFactory().CreateMetricsToMetrics(context.Background(),
		connectortest.NewNopCreateSettings(), cfg, router.(consumer.Metrics))

	require.NoError(t, err)

	failoverConnector := conn.(*metricsFailover)
	failoverConnector.failover.ModifyConsumerAtIndex(0, consumertest.NewErr(errMetricsConsumer))
	failoverConnector.failover.ModifyConsumerAtIndex(1, consumertest.NewErr(errMetricsConsumer))
	failoverConnector.failover.ModifyConsumerAtIndex(2, consumertest.NewErr(errMetricsConsumer))
	defer func() {
		assert.NoError(t, failoverConnector.Shutdown(context.Background()))
	}()

	md := sampleMetric()

	assert.EqualError(t, conn.ConsumeMetrics(context.Background(), md), "All provided pipelines return errors")
}

func consumeMetricsAndCheckStable(conn *metricsFailover, idx int, mr pmetric.Metrics) bool {
	_ = conn.ConsumeMetrics(context.Background(), mr)
	stableIndex := conn.failover.pS.TestStableIndex()
	return stableIndex == idx
}

func sampleMetric() pmetric.Metrics {
	m := pmetric.NewMetrics()
	rm := m.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutInt("sample", 1)
	metric := rm.ScopeMetrics().AppendEmpty().Metrics().AppendEmpty()
	metric.SetEmptySum()
	metric.SetName("test")
	return m
}

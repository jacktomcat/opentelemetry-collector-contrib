// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rabbitmqexporter // import "github.com/jacktomcat/opentelemetry-collector-contrib/exporter/rabbitmqexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/jacktomcat/opentelemetry-collector-contrib/exporter/rabbitmqexporter/internal/metadata"
)

const (
	defaultEncoding = "otlp_proto"
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.LogsStability),
		exporter.WithMetrics(createMetricsExporter, metadata.TracesStability),
		exporter.WithTraces(createTracesExporter, metadata.LogsStability),
	)
}

func createDefaultConfig() component.Config {
	retrySettings := configretry.BackOffConfig{
		Enabled: false,
	}
	return &Config{
		MessageBodyEncoding: defaultEncoding,
		Durable:             true,
		RetrySettings:       retrySettings,
	}
}

func createTracesExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	config := cfg.(*Config)
	r := newRabbitmqExporter(config, set.TelemetrySettings)

	return exporterhelper.NewTracesExporter(
		ctx,
		set,
		cfg,
		r.pushTraces,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(r.start),
		exporterhelper.WithShutdown(r.shutdown),
		exporterhelper.WithRetry(config.RetrySettings),
	)
}

func createMetricsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Metrics, error) {
	config := (cfg.(*Config))
	r := newRabbitmqExporter(config, set.TelemetrySettings)

	return exporterhelper.NewMetricsExporter(
		ctx,
		set,
		cfg,
		r.pushMetrics,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(r.start),
		exporterhelper.WithShutdown(r.shutdown),
		exporterhelper.WithRetry(config.RetrySettings),
	)
}

func createLogsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Logs, error) {
	config := (cfg.(*Config))
	r := newRabbitmqExporter(config, set.TelemetrySettings)

	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,
		r.pushLogs,
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: false}),
		exporterhelper.WithStart(r.start),
		exporterhelper.WithShutdown(r.shutdown),
		exporterhelper.WithRetry(config.RetrySettings),
	)
}

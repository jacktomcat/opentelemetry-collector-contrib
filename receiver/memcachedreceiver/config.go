// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package memcachedreceiver // import "github.com/jacktomcat/opentelemetry-collector-contrib/receiver/memcachedreceiver"

import (
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/jacktomcat/opentelemetry-collector-contrib/receiver/memcachedreceiver/internal/metadata"
)

type Config struct {
	scraperhelper.ControllerConfig `mapstructure:",squash"`
	confignet.AddrConfig           `mapstructure:",squash"`

	// MetricsBuilderConfig allows customizing scraped metrics/attributes representation.
	metadata.MetricsBuilderConfig `mapstructure:",squash"`
}

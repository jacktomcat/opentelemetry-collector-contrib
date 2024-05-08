// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package mongodbreceiver // import "github.com/jacktomcat/opentelemetry-collector-contrib/receiver/mongodbreceiver"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/jacktomcat/opentelemetry-collector-contrib/receiver/mongodbreceiver/internal/metadata"
)

func TestType(t *testing.T) {
	factory := NewFactory()
	ft := factory.Type()
	require.EqualValues(t, metadata.Type, ft)
}

func TestValidConfig(t *testing.T) {
	factory := NewFactory()
	require.NoError(t, component.ValidateConfig(factory.CreateDefaultConfig()))
}

func TestCreateMetricsReceiver(t *testing.T) {
	factory := NewFactory()
	_, err := factory.CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		&Config{
			ControllerConfig: scraperhelper.ControllerConfig{
				CollectionInterval: 10 * time.Second,
				InitialDelay:       time.Second,
			},
		},
		consumertest.NewNop(),
	)
	require.NoError(t, err)
}

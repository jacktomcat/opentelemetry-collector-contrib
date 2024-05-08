// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package apachereceiver

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/jacktomcat/opentelemetry-collector-contrib/receiver/apachereceiver/internal/metadata"
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
	metricsReceiver, err := factory.CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		&Config{
			ControllerConfig: scraperhelper.ControllerConfig{
				CollectionInterval: 10 * time.Second,
			},
		},
		consumertest.NewNop(),
	)
	require.NoError(t, err)
	require.NotNil(t, metricsReceiver)
}

func TestPortValidate(t *testing.T) {
	testCases := []struct {
		desc         string
		endpoint     string
		expectedPort string
	}{
		{
			desc:         "http with specified port",
			endpoint:     "http://localhost:8080/server-status?auto",
			expectedPort: "8080",
		},
		{
			desc:         "http without specified port",
			endpoint:     "http://localhost/server-status?auto",
			expectedPort: "80",
		},
		{
			desc:         "https with specified port",
			endpoint:     "https://localhost:8080/server-status?auto",
			expectedPort: "8080",
		},
		{
			desc:         "https without specified port",
			endpoint:     "https://localhost/server-status?auto",
			expectedPort: "443",
		},
		{
			desc:         "unknown protocol with specified port",
			endpoint:     "abc://localhost:8080/server-status?auto",
			expectedPort: "8080",
		},
		{
			desc:         "port unresolvable",
			endpoint:     "abc://localhost/server-status?auto",
			expectedPort: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := NewFactory().CreateDefaultConfig().(*Config)
			cfg.Endpoint = tc.endpoint
			_, port, err := parseResourceAttributes(tc.endpoint)

			require.NoError(t, err)
			require.Equal(t, tc.expectedPort, port)
		})
	}
}

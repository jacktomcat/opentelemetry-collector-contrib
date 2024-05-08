// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package skywalkingexporter

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/jacktomcat/opentelemetry-collector-contrib/exporter/skywalkingexporter/internal/metadata"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	defaultCfg := createDefaultConfig().(*Config)
	defaultCfg.Endpoint = "1.2.3.4:11800"

	tests := []struct {
		id       component.ID
		expected component.Config
	}{
		{
			id:       component.NewIDWithName(metadata.Type, ""),
			expected: defaultCfg,
		},
		{
			id: component.NewIDWithName(metadata.Type, "2"),
			expected: &Config{
				BackOffConfig: configretry.BackOffConfig{
					Enabled:             true,
					InitialInterval:     10 * time.Second,
					MaxInterval:         1 * time.Minute,
					MaxElapsedTime:      10 * time.Minute,
					RandomizationFactor: backoff.DefaultRandomizationFactor,
					Multiplier:          backoff.DefaultMultiplier,
				},
				QueueSettings: exporterhelper.QueueSettings{
					Enabled:      true,
					NumConsumers: 2,
					QueueSize:    10,
				},
				TimeoutSettings: exporterhelper.TimeoutSettings{
					Timeout: 10 * time.Second,
				},
				ClientConfig: configgrpc.ClientConfig{
					Headers: map[string]configopaque.String{
						"can you have a . here?": "F0000000-0000-0000-0000-000000000000",
						"header1":                "234",
						"another":                "somevalue",
					},
					Endpoint:    "1.2.3.4:11800",
					Compression: "gzip",
					TLSSetting: configtls.ClientConfig{
						Config: configtls.Config{
							CAFile: "/var/lib/mycert.pem",
						},
						Insecure: false,
					},
					Keepalive: &configgrpc.KeepaliveClientConfig{
						Time:                20,
						PermitWithoutStream: true,
						Timeout:             30,
					},
					WriteBufferSize: 512 * 1024,
					BalancerName:    "round_robin",
				},
				NumStreams: 233,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestValidate(t *testing.T) {
	c1 := &Config{
		ClientConfig: configgrpc.ClientConfig{
			Endpoint: "",
		},
		NumStreams: 3,
	}
	err := c1.Validate()
	assert.Error(t, err)
	c2 := &Config{
		ClientConfig: configgrpc.ClientConfig{
			Endpoint: "",
		},
		NumStreams: 0,
	}
	err2 := c2.Validate()
	assert.Error(t, err2)
}

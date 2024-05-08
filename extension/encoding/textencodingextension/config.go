// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package textencodingextension // import "github.com/jacktomcat/opentelemetry-collector-contrib/extension/encoding/textencodingextension"

type Config struct {
	Encoding string `mapstructure:"encoding"`
}

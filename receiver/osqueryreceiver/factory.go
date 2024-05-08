// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package osqueryreceiver // import "github.com/jacktomcat/opentelemetry-collector-contrib/receiver/osqueryreceiver"

import (
	"go.opentelemetry.io/collector/receiver"

	"github.com/jacktomcat/opentelemetry-collector-contrib/receiver/osqueryreceiver/internal/metadata"
)

func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
	)
}

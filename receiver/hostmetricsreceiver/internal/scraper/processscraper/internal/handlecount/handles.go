// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package handlecount // import "github.com/jacktomcat/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/processscraper/internal/handlecount"

type Manager interface {
	Refresh() error
	GetProcessHandleCount(pid int64) (uint32, error)
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows
// +build windows

package tracker // import "github.com/jacktomcat/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/tracker"

import (
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/fileset"
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/stanza/fileconsumer/internal/reader"
)

// On windows, we close files immediately after reading because they cannot be moved while open.
func (t *Tracker) EndConsume() {
	// t.currentPollFiles -> t.previousPollFiles
	t.previousPollFiles = t.currentPollFiles
	t.ClosePreviousFiles()
	t.currentPollFiles = fileset.New[*reader.Reader](t.maxBatchFiles)
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package filestorage // import "github.com/jacktomcat/opentelemetry-collector-contrib/extension/storage/filestorage"

import (
	"os"
	"path/filepath"
)

func getDefaultDirectory() string {
	return filepath.Join(os.Getenv("ProgramData"), "Otelcol", "FileStorage")
}

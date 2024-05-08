// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package logs // import "github.com/jacktomcat/opentelemetry-collector-contrib/processor/transformprocessor/internal/logs"

import (
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
)

func LogFunctions() map[string]ottl.Factory[ottllog.TransformContext] {
	// No logs-only functions yet.
	return ottlfuncs.StandardFuncs[ottllog.TransformContext]()
}

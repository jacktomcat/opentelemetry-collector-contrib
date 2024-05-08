// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package traces // import "github.com/jacktomcat/opentelemetry-collector-contrib/processor/transformprocessor/internal/traces"

import (
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspanevent"
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
)

func SpanFunctions() map[string]ottl.Factory[ottlspan.TransformContext] {
	// No trace-only functions yet.
	return ottlfuncs.StandardFuncs[ottlspan.TransformContext]()
}

func SpanEventFunctions() map[string]ottl.Factory[ottlspanevent.TransformContext] {
	// No trace-only functions yet.
	return ottlfuncs.StandardFuncs[ottlspanevent.TransformContext]()
}

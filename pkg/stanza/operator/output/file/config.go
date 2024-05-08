// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package file // import "github.com/jacktomcat/opentelemetry-collector-contrib/pkg/stanza/operator/output/file"

import (
	"fmt"
	"html/template"

	"go.uber.org/zap"

	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/jacktomcat/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
)

const operatorType = "file_output"

func init() {
	operator.Register(operatorType, func() operator.Builder { return NewConfig("") })
}

// NewConfig creates a new file output config with default values
func NewConfig(operatorID string) *Config {
	return &Config{
		OutputConfig: helper.NewOutputConfig(operatorID, operatorType),
	}
}

// Config is the configuration of a file output operatorn.
type Config struct {
	helper.OutputConfig `mapstructure:",squash"`

	Path   string `mapstructure:"path"`
	Format string `mapstructure:"format"`
}

// Build will build a file output operator.
func (c Config) Build(logger *zap.SugaredLogger) (operator.Operator, error) {
	outputOperator, err := c.OutputConfig.Build(logger)
	if err != nil {
		return nil, err
	}

	var tmpl *template.Template
	if c.Format != "" {
		tmpl, err = template.New("file").Parse(c.Format)
		if err != nil {
			return nil, err
		}
	}

	if c.Path == "" {
		return nil, fmt.Errorf("must provide a path to output to")
	}

	return &Output{
		OutputOperator: outputOperator,
		path:           c.Path,
		tmpl:           tmpl,
	}, nil
}

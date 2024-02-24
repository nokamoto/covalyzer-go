package config

import (
	"encoding/json"
	"fmt"
	"os"

	v1 "github.com/nokamoto/covalyzer-go/pkg/api/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"gopkg.in/yaml.v3"
)

// NewConfig returns a new Config.
// It reads a YAML file from the given path and unmarshals it into a Config.
func NewConfig(path string) (*v1.Config, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var out map[string]interface{}
	if err := yaml.Unmarshal(bytes, &out); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	bytes, err = json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	var res v1.Config
	if err := protojson.Unmarshal(bytes, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protojson: %w", err)
	}

	return &res, nil
}

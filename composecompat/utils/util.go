// Package utils provides the small Compose conversion helpers required by the CLI.
package utils

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// ConvertByJSON copies compatible data between structures through JSON.
func ConvertByJSON(source, target interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		logrus.Errorf("Failed to unmarshal JSON: %v\n%s", err, string(data))
	}
	return err
}

// Convert copies compatible data between structures through YAML.
func Convert(source, target interface{}) error {
	data, err := yaml.Marshal(source)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, target)
	if err != nil {
		logrus.Errorf("Failed to unmarshal YAML: %v\n%s", err, string(data))
	}
	return err
}

// CopySlice returns an independent copy of a string slice.
func CopySlice(source []string) []string {
	if source == nil {
		return nil
	}
	result := make([]string, len(source))
	copy(result, source)
	return result
}

// CopyMap returns an independent copy of a string map.
func CopyMap(source map[string]string) map[string]string {
	if source == nil {
		return nil
	}
	result := make(map[string]string, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

// Contains reports whether item exists in list.
func Contains(list []string, item string) bool {
	for _, value := range list {
		if value == item {
			return true
		}
	}
	return false
}

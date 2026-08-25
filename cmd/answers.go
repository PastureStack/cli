package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cli "github.com/PastureStack/cli/internal/legacycli"
	"gopkg.in/yaml.v3"
)

func parseAnswers(ctx *cli.Context) (map[string]interface{}, error) {
	file := ctx.String("answers")
	if file == "" {
		return map[string]interface{}{}, nil
	}
	if extension := strings.ToLower(filepath.Ext(file)); extension == ".yaml" || extension == ".yml" || extension == ".json" {
		return parseStructuredAnswers(file, extension)
	}
	return parseEnvironmentAnswers(file)
}

func parseStructuredAnswers(file, extension string) (map[string]interface{}, error) {
	contents, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	values := map[string]interface{}{}
	if extension == ".json" {
		err = json.Unmarshal(contents, &values)
	} else {
		err = yaml.Unmarshal(contents, &values)
	}
	if err != nil {
		return nil, fmt.Errorf("parse answers file %s: %w", file, err)
	}
	for key, value := range values {
		switch value.(type) {
		case nil:
			values[key] = ""
		case string, bool, int, int64, uint64, float64:
			values[key] = fmt.Sprint(value)
		default:
			return nil, fmt.Errorf("answer %q must be a scalar, got %T", key, value)
		}
	}
	return values, nil
}

func parseEnvironmentAnswers(file string) (map[string]interface{}, error) {
	input, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	values := map[string]interface{}{}
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, fmt.Errorf("answers file contains an empty key")
		}
		value := ""
		if len(parts) == 2 {
			value = parts[1]
		}
		values[key] = value
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return values, nil
}

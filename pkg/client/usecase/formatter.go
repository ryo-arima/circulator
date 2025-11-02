package usecase

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format formats data according to the specified output format
func Format(outputFormat string, data interface{}) string {
	switch strings.ToLower(outputFormat) {
	case "json":
		return formatJSON(data)
	case "yaml":
		return formatYAML(data)
	case "table":
		return formatTable(data)
	default:
		return formatTable(data)
	}
}

func formatJSON(data interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error formatting JSON: %v", err)
	}
	return string(jsonBytes) + "\n"
}

func formatYAML(data interface{}) string {
	yamlBytes, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Sprintf("Error formatting YAML: %v", err)
	}
	return string(yamlBytes)
}

func formatTable(data interface{}) string {
	// Simple table formatting - in a real implementation this would be more sophisticated
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		var result strings.Builder
		typ := val.Type()

		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			value := val.Field(i)

			// Skip unexported fields
			if !field.IsExported() {
				continue
			}

			result.WriteString(fmt.Sprintf("%-15s: %v\n", field.Name, value.Interface()))
		}

		return result.String()
	}

	// Fallback to simple string representation
	return fmt.Sprintf("%+v\n", data)
}

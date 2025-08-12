package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type SwaggerSpec struct {
	Swagger     string              `json:"swagger" yaml:"swagger"`
	Paths       map[string]PathItem `json:"paths" yaml:"paths"`
	Definitions map[string]Schema   `json:"definitions" yaml:"definitions"`
}

type PathItem struct {
	Get    *Operation `json:"get" yaml:"get"`
	Post   *Operation `json:"post" yaml:"post"`
	Put    *Operation `json:"put" yaml:"put"`
	Delete *Operation `json:"delete" yaml:"delete"`
}

type Operation struct {
	Tags        []string            `json:"tags" yaml:"tags"`
	OperationID string              `json:"operationId" yaml:"operationId"`
	Parameters  []Parameter         `json:"parameters" yaml:"parameters"`
	Responses   map[string]Response `json:"responses" yaml:"responses"`
	Consumes    []string            `json:"consumes" yaml:"consumes"`
}

type Parameter struct {
	Name     string  `json:"name" yaml:"name"`
	In       string  `json:"in" yaml:"in"`
	Required bool    `json:"required" yaml:"required"`
	Type     string  `json:"type" yaml:"type"`
	Schema   *Schema `json:"schema" yaml:"schema"`
}

type Response struct {
	Description string  `json:"description" yaml:"description"`
	Schema      *Schema `json:"schema" yaml:"schema"`
}

type Schema struct {
	Ref         string            `json:"$ref" yaml:"$ref"`
	Type        string            `json:"type" yaml:"type"`
	Properties  map[string]Schema `json:"properties" yaml:"properties"`
	Items       *Schema           `json:"items" yaml:"items"`
	Description string            `json:"description" yaml:"description"`
	Required    []string          `json:"required" yaml:"required"`
}

func writeSchema(output *strings.Builder, schema *Schema, indent int, definitions map[string]Schema) {
	indentation := strings.Repeat(" ", indent)

	if schema.Type == "array" && schema.Items != nil {
		if schema.Items.Type != "" {
			output.WriteString(fmt.Sprintf("%s- items: %s[]\n", indentation, schema.Items.Type))
		} else if schema.Items.Ref != "" {
			refType := strings.TrimPrefix(schema.Items.Ref, "#/definitions/")
			output.WriteString(fmt.Sprintf("%s- items: %s[]\n", indentation, refType))

			if refSchema, ok := definitions[refType]; ok {
				writeSchema(output, &refSchema, indent+2, definitions)
			}
		}
		return
	}

	if schema.Properties != nil {
		for name, prop := range schema.Properties {
			required := ""
			if contains(schema.Required, name) {
				required = " (required)"
			}

			if prop.Type == "array" && prop.Items != nil {
				if prop.Items.Type != "" {
					output.WriteString(fmt.Sprintf("%s- %s: %s[]%s\n", indentation, name, prop.Items.Type, required))
				} else if prop.Items.Ref != "" {
					refType := strings.TrimPrefix(prop.Items.Ref, "#/definitions/")
					output.WriteString(fmt.Sprintf("%s- %s: %s[]%s\n", indentation, name, refType, required))

					if refSchema, ok := definitions[refType]; ok {
						writeSchema(output, &refSchema, indent+2, definitions)
					}
				}
			} else if prop.Type != "" {
				example := getExampleForType(prop.Type)
				output.WriteString(fmt.Sprintf("%s- %s: %s%s | %s\n", indentation, name, prop.Type, required, example))
			} else if prop.Ref != "" {
				refType := strings.TrimPrefix(prop.Ref, "#/definitions/")
				output.WriteString(fmt.Sprintf("%s- %s: %s%s\n", indentation, name, refType, required))

				if refSchema, ok := definitions[refType]; ok {
					writeSchema(output, &refSchema, indent+2, definitions)
				}
			}
		}
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func getExampleForType(paramType string) string {
	switch paramType {
	case "string":
		return "\"example\""
	case "number":
		return "1"
	case "integer":
		return "1"
	case "boolean":
		return "true"
	case "file":
		return "file"
	default:
		return "\"example\""
	}
}

func main() {
	data, err := os.ReadFile("./generate/swagger/swagger.yaml")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	var spc SwaggerSpec
	err = yaml.Unmarshal(data, &spc)
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		return
	}

	var output strings.Builder

	for path, pathItem := range spc.Paths {
		if pathItem.Get != nil {
			output.WriteString(formatOperation(spc, pathItem.Get, "get", path))
		}
		if pathItem.Post != nil {
			output.WriteString(formatOperation(spc, pathItem.Post, "post", path))
		}
		if pathItem.Put != nil {
			output.WriteString(formatOperation(spc, pathItem.Put, "put", path))
		}
		if pathItem.Delete != nil {
			output.WriteString(formatOperation(spc, pathItem.Delete, "delete", path))
		}
	}

	err = os.WriteFile("./generate/swagger/swagger.md", []byte(output.String()), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		return
	}
}

func formatOperation(spc SwaggerSpec, op *Operation, method string, path string) string {
	var output strings.Builder

	var serviceName string
	if len(op.Tags) > 0 && op.OperationID != "" {
		serviceName = fmt.Sprintf("backend.%s.%s", op.Tags[0], op.OperationID)
	} else {
		serviceName = "backend.unknown.unknown"
	}

	output.WriteString("# ")
	output.WriteString(fmt.Sprintf("%s\n\n", serviceName))

	var queryParams []string
	var bodyParam string
	var formParams []string

	for _, param := range op.Parameters {
		switch param.In {
		case "query":
			paramType := param.Type
			if paramType == "integer" {
				paramType = "number"
			}
			example := getExampleForType(paramType)
			queryParams = append(queryParams, fmt.Sprintf("   - %s: %s | %s", param.Name, paramType, example))

		case "body":
			if param.Schema != nil && param.Schema.Ref != "" {
				bodyParam = strings.TrimPrefix(param.Schema.Ref, "#/definitions/")
			} else if param.Schema != nil && param.Schema.Type == "array" && param.Schema.Items != nil && param.Schema.Items.Ref != "" {
				bodyParam = strings.TrimPrefix(param.Schema.Items.Ref, "#/definitions/") + "[]"
			}

		case "formData":
			required := ""
			if param.Required {
				required = " (required)"
			}
			example := getExampleForType(param.Type)
			formParams = append(formParams, fmt.Sprintf("   - %s: %s%s | %s", param.Name, param.Type, required, example))
		}
	}

	// * add query parameters
	if len(queryParams) > 0 {
		output.WriteString(" - query\n")
		output.WriteString(strings.Join(queryParams, "\n"))
		output.WriteString("\n")
	}

	// * add form parameters (for multipart form data)
	if len(formParams) > 0 {
		output.WriteString(" - form\n")
		output.WriteString(strings.Join(formParams, "\n"))
		output.WriteString("\n")
	}

	// * add body parameters (for JSON)
	if bodyParam != "" {
		output.WriteString(" - body\n")
		if strings.HasSuffix(bodyParam, "[]") {
			arrayType := strings.TrimSuffix(bodyParam, "[]")
			if schema, ok := spc.Definitions[arrayType]; ok {
				output.WriteString(fmt.Sprintf("   - items: %s[]\n", arrayType))
				writeSchema(&output, &schema, 5, spc.Definitions)
			}
		} else if schema, ok := spc.Definitions[bodyParam]; ok {
			writeSchema(&output, &schema, 3, spc.Definitions)
		}
	}

	// * add response
	if resp, ok := op.Responses["200"]; ok && resp.Schema != nil {
		output.WriteString(" - response\n")
		if resp.Schema.Ref != "" {
			responseType := strings.TrimPrefix(resp.Schema.Ref, "#/definitions/")
			if schema, ok := spc.Definitions[responseType]; ok {
				writeSchema(&output, &schema, 3, spc.Definitions)
			}
		} else if resp.Schema.Type == "array" && resp.Schema.Items != nil {
			if resp.Schema.Items.Ref != "" {
				responseType := strings.TrimPrefix(resp.Schema.Items.Ref, "#/definitions/")
				output.WriteString(fmt.Sprintf("   - items: %s[]\n", responseType))
				if schema, ok := spc.Definitions[responseType]; ok {
					writeSchema(&output, &schema, 5, spc.Definitions)
				}
			} else if resp.Schema.Items.Type != "" {
				output.WriteString(fmt.Sprintf("   - items: %s[]\n", resp.Schema.Items.Type))
			}
		} else {
			writeSchema(&output, resp.Schema, 3, spc.Definitions)
		}
	}

	output.WriteString("\n")
	return output.String()
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lithammer/dedent"
)

type EndpointInfo struct {
	Name       string
	Method     string
	Path       string
	Tag        string
	ReturnType string
	ErrorType  string
	QueryType  string
	BodyType   string
	FormFields []*FormField
}

type FormField struct {
	Name     string
	Type     string
	Required bool
	IsFile   bool
}

type RouteInfo struct {
	Method      string
	Path        string
	HandlerName string
}

func main() {
	_ = os.MkdirAll("./generate/swagger", 0755)

	outputFile, err := os.Create("./generate/swagger/declaration.go")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = outputFile.Close()
	}()

	_, _ = outputFile.WriteString("package swagger\n\nimport (\n\t_ \"backend/type/payload\"\n\t_ \"backend/type/response\"\n)\n\n")

	endpointContent, err := os.ReadFile("./endpoint/endpoint.go")
	if err != nil {
		panic(err)
	}

	routes := extractRoutes(string(endpointContent))

	err = filepath.Walk("./endpoint", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(info.Name(), ".go") ||
			strings.HasSuffix(info.Name(), "_test.go") ||
			info.Name() == "endpoint.go" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		endpoints := extractEndpointInfo(string(content), routes)

		for _, endpoint := range endpoints {
			swagger := generateSwaggerComment(endpoint)
			_, _ = outputFile.WriteString(swagger)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}

func extractRoutes(content string) []RouteInfo {
	var routes []RouteInfo

	routeRegex := regexp.MustCompile(`(\w+)\.(Get|Post|Put|Delete)\("([^"]+)",\s*\w+\.(\w+)\)`)

	matches := routeRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		groupName := match[1]
		method := match[2]
		path := match[3]
		handlerName := match[4]

		if groupName != "app" {
			path = "/" + strings.ToLower(groupName) + path
		}

		routes = append(routes, RouteInfo{
			Method:      strings.ToLower(method),
			Path:        path,
			HandlerName: handlerName,
		})
	}

	return routes
}

func extractEndpointInfo(content string, routes []RouteInfo) []EndpointInfo {
	var endpoints []EndpointInfo

	handlerRegex := regexp.MustCompile(`func\s+\([^)]+\)\s+(Handle[^\s(]+)`)
	handlerMatches := handlerRegex.FindAllStringSubmatch(content, -1)

	for _, handler := range handlerMatches {
		handlerName := handler[1]

		var matchedRoute *RouteInfo
		for _, route := range routes {
			if route.HandlerName == handlerName {
				matchedRoute = &route
				break
			}
		}

		if matchedRoute == nil {
			continue
		}

		endpoint := EndpointInfo{
			Name:       handlerName,
			Method:     matchedRoute.Method,
			Path:       matchedRoute.Path,
			ErrorType:  "response.ErrorResponse",
			ReturnType: "response.SuccessResponse", // * default return type
		}

		// * extract tag from path
		pathParts := strings.Split(strings.Trim(endpoint.Path, "/"), "/")
		if len(pathParts) > 0 {
			endpoint.Tag = pathParts[0]
		}

		// * find handler content using function name - use (?s) for multiline matching
		handlerContentRegex := regexp.MustCompile(fmt.Sprintf(`(?s)func\s+\([^)]+\)\s+%s[^{]*{(.*?)\n}`, regexp.QuoteMeta(handlerName)))
		handlerContent := handlerContentRegex.FindStringSubmatch(content)

		if len(handlerContent) > 1 {
			functionBody := handlerContent[1]

			// * check for query parser
			queryRegex := regexp.MustCompile(`(\w+)\s*:=\s*new\(([^)]+)\)[\s\n]*if\s+err\s*:=\s*c\.QueryParser\(`)
			queryMatch := queryRegex.FindStringSubmatch(functionBody)
			if len(queryMatch) > 2 {
				endpoint.QueryType = queryMatch[2]
			}

			// * check for body parser
			bodyRegex := regexp.MustCompile(`(\w+)\s*:=\s*new\(([^)]+)\)[\s\n]*if\s+err\s*:=\s*c\.BodyParser\(`)
			bodyMatch := bodyRegex.FindStringSubmatch(functionBody)
			if len(bodyMatch) > 2 {
				endpoint.BodyType = bodyMatch[2]
			}

			// * check for form fields
			if endpoint.BodyType == "" {
				formFields := extractFormFields(functionBody)
				if len(formFields) > 0 {
					endpoint.FormFields = formFields
				}
			}

			// * track variable declarations and their types
			varTypes := make(map[string]string)
			varIsArray := make(map[string]bool)

			// * case 1: var varName *Type - explicit variable declarations
			varDeclRegex := regexp.MustCompile(`var\s+(\w+)\s+\*([^\s]+)`)
			for _, varMatch := range varDeclRegex.FindAllStringSubmatch(functionBody, -1) {
				if len(varMatch) > 2 {
					varName := varMatch[1]
					varType := varMatch[2]
					varTypes[varName] = varType
					varIsArray[varName] = false
				}
			}

			// * case 2: varName := &Type{...} - variable assignment with type initialization
			assignRegex := regexp.MustCompile(`(\w+)\s*:=\s*&([^\s{]+)`)
			for _, assignMatch := range assignRegex.FindAllStringSubmatch(functionBody, -1) {
				if len(assignMatch) > 2 {
					varName := assignMatch[1]
					varType := strings.TrimSpace(assignMatch[2])
					varTypes[varName] = varType
					varIsArray[varName] = false
				}
			}

			// * case 3: var varName []*Type or var varName []Type - array variable declarations
			arrayDeclRegex := regexp.MustCompile(`var\s+(\w+)\s+\[\]\*?([^\s]+)`)
			for _, arrayMatch := range arrayDeclRegex.FindAllStringSubmatch(functionBody, -1) {
				if len(arrayMatch) > 2 {
					varName := arrayMatch[1]
					varType := arrayMatch[2]
					varTypes[varName] = varType
					varIsArray[varName] = true
				}
			}

			// * case 4: varName := make([]*Type, 0) or varName := []*Type{} - array initialization
			arrayInitRegex := regexp.MustCompile(`(\w+)\s*:=\s*(make\(\[\]\*?([^\s,]+)|(\[\]\*?([^\s{]+)))`)
			for _, arrayInitMatch := range arrayInitRegex.FindAllStringSubmatch(functionBody, -1) {
				if len(arrayInitMatch) > 2 {
					varName := arrayInitMatch[1]
					// * handle both make([]*Type, 0) and []*Type{} cases
					var varType string
					if arrayInitMatch[3] != "" {
						varType = arrayInitMatch[3]
					} else if arrayInitMatch[5] != "" {
						varType = arrayInitMatch[5]
					}
					if varType != "" {
						varTypes[varName] = varType
						varIsArray[varName] = true
					}
				}
			}

			// * find specific array declaration pattern: var varName []*structType
			specificArrayRegex := regexp.MustCompile(`var\s+(\w+)\s+\[]\*([^\s]+)`)
			specificArrayMatches := specificArrayRegex.FindAllStringSubmatch(functionBody, -1)
			for _, match := range specificArrayMatches {
				if len(match) > 2 {
					varName := match[1]
					varType := match[2]
					varTypes[varName] = varType
					varIsArray[varName] = true
				}
			}

			// * find array initializations like: varName := make([]*Type, 0)
			makeArrayRegex := regexp.MustCompile(`(\w+)\s*:=\s*make\(\[]\*([^,]+)`)
			makeArrayMatches := makeArrayRegex.FindAllStringSubmatch(functionBody, -1)
			for _, match := range makeArrayMatches {
				if len(match) > 2 {
					varName := match[1]
					varType := match[2]
					varTypes[varName] = varType
					varIsArray[varName] = true
				}
			}

			// * find the return statement with response.Success
			returnRegex := regexp.MustCompile(`return\s+c\.JSON\(response\.Success\(c, ([^)]+)\)\)`)
			returnMatch := returnRegex.FindStringSubmatch(functionBody)

			if len(returnMatch) > 1 {
				successArg := strings.TrimSpace(returnMatch[1])

				// * check if argument is a tracked variable
				if varType, exists := varTypes[successArg]; exists {
					if isArray, ok := varIsArray[successArg]; ok && isArray {
						// * it's an array type
						endpoint.ReturnType = fmt.Sprintf("response.GenericResponse[[]%s]", varType)
					} else {
						// * it's a single object type
						endpoint.ReturnType = fmt.Sprintf("response.GenericResponse[%s]", varType)
					}
				} else if strings.HasPrefix(successArg, "&") {
					// * handle inline struct creation &Type{...}
					typeRegex := regexp.MustCompile(`&([^\s{]+)`)
					typeMatch := typeRegex.FindStringSubmatch(successArg)

					if len(typeMatch) > 1 {
						typeName := strings.TrimSpace(typeMatch[1])
						endpoint.ReturnType = fmt.Sprintf("response.GenericResponse[%s]", typeName)
					}
				}
			}
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}

func extractFormFields(functionBody string) []*FormField {
	var formFields []*FormField

	// * detect form value
	formValueRegex := regexp.MustCompile(`(\w+)\s*:=\s*c\.FormValue\("([^"]+)"\)`)
	formValueMatches := formValueRegex.FindAllStringSubmatch(functionBody, -1)

	for _, match := range formValueMatches {
		if len(match) > 2 {
			fieldName := match[2]

			formFields = append(formFields, &FormField{
				Name:     fieldName,
				Type:     "string",
				Required: true,
				IsFile:   false,
			})
		}
	}

	// * detect form file
	formFileRegex := regexp.MustCompile(`[^,\s]+,\s*err\s*:=\s*c\.FormFile\("([^"]+)"\)`)
	formFileMatches := formFileRegex.FindAllStringSubmatch(functionBody, -1)

	for _, match := range formFileMatches {
		if len(match) > 1 {
			fieldName := match[1]

			formFields = append(formFields, &FormField{
				Name:     fieldName,
				Type:     "file",
				Required: true,
				IsFile:   true,
			})
		}
	}

	return formFields
}

func generateSwaggerComment(endpoint EndpointInfo) string {
	// * generate id from endpoint name
	name := strings.TrimPrefix(endpoint.Name, "Handle")
	id := strings.ToLower(name[:1]) + name[1:]

	// * build swagger parameters
	params := ""

	// * add query parameters
	if endpoint.QueryType != "" {
		params += fmt.Sprintf("\n// @Param query query %s true \"Query\"", endpoint.QueryType)
	}

	// * add body parameters
	if endpoint.BodyType != "" {
		params += fmt.Sprintf("\n// @Param body body %s true \"Body\"", endpoint.BodyType)
	}

	// * add form parameters
	if len(endpoint.FormFields) > 0 {
		for _, field := range endpoint.FormFields {
			required := "true"
			if !field.Required {
				required = "false"
			}

			paramType := "formData"
			fieldType := field.Type

			params += fmt.Sprintf("\n// @Param %s %s %s %s \"%s\"",
				field.Name, paramType, fieldType, required, field.Name)
		}

		// * add consumes annotation for multipart form
		params += "\n// @Accept multipart/form-data"
	}

	return fmt.Sprintf(dedent.Dedent(`
			// %s
			// @ID %s
			// @Tags %s%s
			// @Success 200 {object} %s
			// @Failure 400 {object} %s
			// @Router %s [%s]
			func %s() {
				_ = 0
			}
		`),
		endpoint.Name,
		id,
		endpoint.Tag,
		params,
		endpoint.ReturnType,
		endpoint.ErrorType,
		endpoint.Path,
		endpoint.Method,
		endpoint.Name,
	)
}

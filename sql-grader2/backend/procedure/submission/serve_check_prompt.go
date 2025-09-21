package submissionProcedure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bsthun/gut"
	"github.com/lithammer/dedent"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

type PromptCheckResult struct {
	Passed      bool   `json:"passed"`
	Description string `json:"description"`
}

func (r *Service) ServeCheckPrompt(ctx context.Context, studentAnswer string, checkPrompt string) (bool, string, *gut.ErrorInstance) {
	systemPrompt := dedent.Dedent(`You are a SQL validation expert. Analyze the given SQL statement against the provided requirements and determine if it meets all the criteria.

		Rules:
		1. Check for specific SQL features mentioned in the requirements (e.g., JOIN types, subqueries, specific clauses)
		2. Validate SQL syntax and structure
		3. Ensure the query follows best practices if specified
		4. If the requirement mentions specific JOIN types (LEFT JOIN, RIGHT JOIN, INNER JOIN, etc.), verify they are used correctly
		5. Check for required clauses like ORDER BY, GROUP BY, HAVING, LIMIT, etc. if specified
		6. Validate proper use of aggregate functions, window functions, etc. if required
		
		Respond with:
		- passed: true/false based on whether ALL requirements are met
		- description: specific reasons why it passed or failed, explanation of the analysis, short 2 sentences or less when passed`)

	userPrompt := fmt.Sprintf(dedent.Dedent(`Requirements to check:
		%s
		
		SQL statement to validate:
		%s
		
		Analyze if this SQL statement meets all the requirements specified above.`),
		checkPrompt,
		studentAnswer)

	chatParams := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model:       *r.config.OpenaiModel,
		MaxTokens:   openai.Int(128),
		Temperature: openai.Float(0.01),
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
				Type: "json_schema",
				JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:        "prompt_validation",
					Description: openai.String("SQL prompt validation result"),
					Schema: map[string]any{
						"type": "object",
						"properties": map[string]any{
							"passed": map[string]any{
								"type":        "boolean",
								"description": "Whether the SQL statement meets all requirements",
							},
							"description": map[string]any{
								"type":        "string",
								"description": "Detailed explanation of why it passed or failed",
							},
						},
						"required": []string{"passed", "description"},
					},
					Strict: openai.Bool(true),
				},
			},
		},
	}

	chatCompletion, err := r.openaiClient.Chat.Completions.New(ctx, chatParams)
	if err != nil {
		return false, "", gut.Err(false, "unable to get chat completion from OpenAI", err)
	}

	responseContent := chatCompletion.Choices[0].Message.Content

	var result *PromptCheckResult
	err = json.Unmarshal([]byte(responseContent), &result)
	if err != nil {
		return false, "", gut.Err(false, "unable to parse OpenAI response", err)
	}

	return result.Passed, result.Description, nil
}

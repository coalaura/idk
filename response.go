package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/coalaura/openingrouter"
)

type CommandAnswer struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

var answerResponseFormat = openingrouter.ChatResponseFormat{
	Type: openingrouter.ChatResponseFormatTypeJSONSchema,
	JSONSchema: &openingrouter.ChatJSONSchema{
		Name:   "answer",
		Strict: new(true),
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command":     map[string]any{"type": "string"},
				"description": map[string]any{"type": "string"},
			},
			"required":             []string{"command", "description"},
			"additionalProperties": false,
		},
	},
}

func parseAnswer(response *openingrouter.ChatCompletionResponse) (CommandAnswer, error) {
	if len(response.Choices) == 0 {
		return CommandAnswer{}, errors.New("no answer returned")
	}

	choice := &response.Choices[0]
	if choice.Message.Refusal != "" {
		return CommandAnswer{}, fmt.Errorf("model refused: %s", choice.Message.Refusal)
	}

	if choice.FinishReason != openingrouter.ChatFinishReasonStop {
		return CommandAnswer{}, fmt.Errorf("incomplete answer: finish reason %q", choice.FinishReason)
	}

	decoder := json.NewDecoder(strings.NewReader(choice.Message.Content.String()))
	decoder.DisallowUnknownFields()

	var result CommandAnswer

	err := decoder.Decode(&result)
	if err != nil {
		return CommandAnswer{}, fmt.Errorf("invalid answer: %w", err)
	}

	err = decoder.Decode(new(any))
	if !errors.Is(err, io.EOF) {
		return CommandAnswer{}, errors.New("invalid answer: trailing JSON data")
	}

	if strings.TrimSpace(result.Command) == "" || strings.TrimSpace(result.Description) == "" {
		return CommandAnswer{}, errors.New("invalid answer: command and description must be nonempty")
	}

	return result, nil
}

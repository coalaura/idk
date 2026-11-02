package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/coalaura/openingrouter"
)

var (
	//go:embed prompt.txt
	Prompt string
)

func answer(config *Config, question string) (*exec.Cmd, string, error) {
	shell, err := resolveShell()
	if err != nil {
		return nil, "", err
	}

	available := availableCommands()

	client := openingrouter.NewClient(config.OpenRouter.ApiKey)

	prompt := fmt.Sprintf(Prompt, runtime.GOOS, filepath.Base(shell), strings.Join(available, ","))

	request := openingrouter.ChatCompletionRequest{
		Model: config.OpenRouter.Model,
		Messages: []openingrouter.ChatMessage{
			openingrouter.SystemMessage(prompt),
			openingrouter.UserMessage(question),
		},
		ResponseFormat: &answerResponseFormat,
		Provider: &openingrouter.ProviderPreferences{
			Only:              config.OpenRouter.Providers,
			Order:             config.OpenRouter.Providers,
			RequireParameters: new(true),
		},
	}

	if len(config.OpenRouter.Fallbacks) > 0 {
		request.Model = ""

		request.Models = make([]string, 0, 1+len(config.OpenRouter.Fallbacks))
		request.Models = append(request.Models, config.OpenRouter.Model)
		request.Models = append(request.Models, config.OpenRouter.Fallbacks...)
	}

	response, err := client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		return nil, "", err
	}

	result, err := parseAnswer(response)
	if err != nil {
		return nil, "", err
	}

	err = result.render(os.Stdout)
	if err != nil {
		return nil, "", err
	}

	var flag string

	if runtime.GOOS == "windows" {
		flag = "/C"
	} else {
		flag = "-c"
	}

	return exec.Command(shell, flag, result.Command), result.Command, nil
}

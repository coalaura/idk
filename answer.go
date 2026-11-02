package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
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

	prompt := fmt.Sprintf(
		Prompt,
		runtime.GOOS,
		filepath.Base(shell),
		available,
	)

	request := openingrouter.ChatCompletionRequest{
		Model: config.OpenRouter.Model,
		Messages: []openingrouter.ChatMessage{
			openingrouter.SystemMessage(prompt),
			openingrouter.UserMessage(question),
		},
	}

	if len(config.OpenRouter.Fallbacks) > 0 {
		request.Model = ""

		request.Models = make([]string, 0, 1+len(config.OpenRouter.Fallbacks))
		request.Models = append(request.Models, config.OpenRouter.Model)
		request.Models = append(request.Models, config.OpenRouter.Fallbacks...)
	}

	if len(config.OpenRouter.Providers) > 0 {
		request.Provider = &openingrouter.ProviderPreferences{
			Only:  config.OpenRouter.Providers,
			Order: config.OpenRouter.Providers,
		}
	}

	stream, err := client.CreateChatCompletionStream(context.Background(), request)
	if err != nil {
		return nil, "", err
	}

	defer stream.Close()

	defer fmt.Println()

	var response strings.Builder

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, "", err
		}

		if len(chunk.Choices) == 0 {
			continue
		}

		delta := chunk.Choices[0].Delta.Content

		if delta != "" {
			response.WriteString(delta)

			fmt.Print(delta)
		}
	}

	rgx := regexp.MustCompile(`(?m)^\$(.+?)$`)
	command := rgx.FindString(response.String())

	command = strings.Trim(command, "$ \r\n")

	var flag string

	if runtime.GOOS == "windows" {
		flag = "/C"
	} else {
		flag = "-c"
	}

	return exec.Command(shell, flag, command), command, nil
}

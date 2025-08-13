package main

import (
	"context"
	_ "embed"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/revrost/go-openrouter"
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

	client := openrouter.NewClient(config.OpenRouter.ApiKey)

	prompt := fmt.Sprintf(
		Prompt,
		runtime.GOOS,
		filepath.Base(shell),
		available,
	)

	request := openrouter.ChatCompletionRequest{
		Model: config.OpenRouter.Model,
		Messages: []openrouter.ChatCompletionMessage{
			openrouter.SystemMessage(prompt),
			openrouter.UserMessage(question),
		},
		Stream: true,
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
			if err == io.EOF {
				break
			}

			return nil, "", err
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

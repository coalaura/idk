package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/coalaura/getch"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: idk <question>")

		os.Exit(1)
	}

	if os.Args[1] == "version" {
		fmt.Printf("idk version %s\n", version)

		return
	}

	config, err := loadConfig()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		os.Exit(1)
	}

	question := strings.Join(os.Args[1:], " ")

	cmd, command, err := answer(config, question)
	if err != nil {
		fmt.Printf("error: %v\n", err)

		os.Exit(1)
	}

	if command == "" || !config.Commands.AskToRun {
		return
	}

	fmt.Print("\n> run command? [y/n]: ")

	ch, err := getch.GetChar()
	if err != nil {
		fmt.Printf("error: %v\n", err)

		os.Exit(1)
	}

	if ch == 0x03 || ch == 'n' || ch == 'N' {
		fmt.Println("no")

		return
	}

	fmt.Println("yes")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

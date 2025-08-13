package main

import (
	"os"
	"os/exec"
	"runtime"
	"sort"

	"github.com/shirou/gopsutil/v4/process"
)

func resolveShell() (string, error) {
	ppid := os.Getppid()

	proc, err := process.NewProcess(int32(ppid))
	if err != nil {
		return "", err
	}

	return proc.Exe()
}

func availableCommands() []string {
	check := []string{"ls", "cat", "rm", "cp", "mv", "ln", "grep", "curl", "wget", "tar", "zip", "unzip"}

	switch runtime.GOOS {
	case "linux":
		check = append(check, "yay", "paru", "apt", "pacman", "dnf", "zypper")
	case "windows":
		check = append(check, "winget", "choco", "scoop")
	}

	available := make([]string, 0, len(check))

	for _, command := range check {
		if _, err := exec.LookPath(command); err != nil {
			continue
		}

		available = append(available, command)
	}

	sort.Strings(available)

	return available
}

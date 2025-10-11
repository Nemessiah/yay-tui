package yay

import (
	"bufio"
	"fmt"
	"os/exec"
)

// Search runs `yay -Ss <query>` and returns its output as a slice of lines.
func Inspect(searchQuery string) ([]string, error) {
	cmd := exec.Command("yay", "-Si", searchQuery)
	output, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to capture output: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to run yay: %v", err)
	}

	var results []string
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("yay failed: %v", err)
	}

	return results, nil
}

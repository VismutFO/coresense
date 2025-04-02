package validation

import (
	"bytes"
	"os/exec"
	"strings"
)

// RunValidationScript executes a Lua script with input and returns validation result.
func RunValidationScript(scriptCode, input string) (bool, error) {
	cmd := exec.Command("lua", "-e", scriptCode)
	cmd.Stdin = bytes.NewBufferString(input)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(out.String()) == "true", nil
}

// RunStatisticsScript executes a Lua script with input and returns statistics.
func RunStatisticsScript(scriptCode string, input []string) (string, error) {
	cmd := exec.Command("lua", "-e", scriptCode)
	cmd.Stdin = bytes.NewBufferString(strings.Join(input, "\n"))

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

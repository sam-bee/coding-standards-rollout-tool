package csfixing

import (
	"os/exec"
	"strings"
)

type systemCallerInterface interface {
	doSystemCall(command string, args []string) ([]string, error)
}

type SystemCaller struct{}

func (sc *SystemCaller) doSystemCall(command string, args []string) ([]string, error) {
    cmd := exec.Command(command, args...)

    out, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(out), "\n")
    return lines, nil
}

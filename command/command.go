package command

import (
	"os/exec"
	"bytes"
)

type Runner interface {
	Run(string, ...string) ([]byte, error)
}

type LocalRunner struct{}

func (r LocalRunner) Run(cmd string, args ...string) ([]byte, error) {
	c := exec.Command(cmd, args...)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := c.Start(); err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	return buf.Bytes(), nil
}


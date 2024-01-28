package ssh

import (
	"bytes"
	"context"
	"io"
	"log"
	"os/exec"

	"github.com/thin-edge/c8y-tedge/pkg/runner"
)

type ExecResult struct {
	exitCode  int
	outBuffer *bytes.Buffer
}

func (r ExecResult) ExitCode() int {
	return r.exitCode
}

func (r ExecResult) Stdout() []byte {
	return r.outBuffer.Bytes()
}

type Runner struct {
	Target  string
	Command string
}

func (r *Runner) Execute(cmd ...string) (runner.Responder, error) {
	fullCommand := append([]string{"-n", r.Target}, cmd...)
	ctx := context.Background()
	execCmd := exec.CommandContext(ctx, "ssh", fullCommand...)

	stdout, _ := execCmd.StdoutPipe()
	// stderr, _ := cmd.StderrPipe()
	err := execCmd.Start()

	if err != nil {
		log.Fatal(err)
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, stdout)
	if err != nil {
		return ExecResult{}, err
	}

	_ = execCmd.Wait()

	result := ExecResult{
		exitCode:  execCmd.ProcessState.ExitCode(),
		outBuffer: buf,
	}

	return result, nil
}

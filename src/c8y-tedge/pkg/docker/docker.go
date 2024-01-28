package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/thin-edge/c8y-tedge/pkg/runner"
)

type DockerRunner struct {
	client    *client.Client
	container string
}

func NewDockerRunner(target ...string) (*DockerRunner, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}
	runner := &DockerRunner{
		client: cli,
	}
	if len(target) > 0 {
		runner.SetTarget(target[0])
	}
	return runner, nil
}

// func (r *DockerRunner) ListComposeProjects
// TODO: Use labels to look for compose projects (this is specified by the spec)
// https://github.com/compose-spec/compose-spec/blob/master/spec.md
// com.docker.compose.project
// com.docker.compose.service

func (r *DockerRunner) ListContainers(ctx context.Context, filterArgs ...filters.KeyValuePair) ([]string, error) {

	filters := filters.NewArgs(filterArgs...)
	// for _, filter := range filterArgs {
	// 	filters.Add()
	// }
	// filters.Add("label", "label1")

	containers, err := r.client.ContainerList(ctx, container.ListOptions{
		Filters: filters,
		// Filters: filters.NewArgs(
		// 	filters.KeyValuePair{
		// 		Key:   "Label",
		// 		Value: "",
		// 	},
		// ),
	})

	if err != nil {
		return nil, err
	}
	names := []string{}
	for _, container := range containers {
		if len(container.Names) > 0 {
			names = append(names, strings.TrimLeft(container.Names[0], "/"))
		}
	}
	return names, nil
}

func (r *DockerRunner) SetTarget(target string) error {
	r.container = target
	return nil
}

func (r *DockerRunner) Validate() error {
	containers, err := r.client.ContainerList(context.Background(), container.ListOptions{
		Limit: 1,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: r.container,
		}),
	})
	if err != nil {
		fmt.Printf("failed to get list of containers")
		return err
	}
	if len(containers) == 0 {
		return fmt.Errorf("container with given name not found")
	}
	// TODO: Should the id be used over the name?
	// r.container = containers[0].ID
	fmt.Printf("failed to get list of containers")
	return nil
}

// ExecResult represents a result returned from Exec()
type ExecResult struct {
	exitCode  int
	outBuffer *bytes.Buffer
	// errBuffer *bytes.Buffer
}

func (r ExecResult) ExitCode() int {
	return r.exitCode
}

func (r ExecResult) Stdout() []byte {
	return r.outBuffer.Bytes()
}

// Exec executes a command inside a container, returning the result
// containing stdout, stderr, and exit code. Note:
//   - this is a synchronous operation;
//   - cmd stdin is closed.
func Exec(ctx context.Context, apiClient client.APIClient, id string, cmd []string, ops ...func(*types.ExecConfig)) (ExecResult, error) {
	// prepare exec
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: false,
		Cmd:          cmd,
	}

	for _, op := range ops {
		op(&execConfig)
	}

	cresp, err := apiClient.ContainerExecCreate(ctx, id, execConfig)
	if err != nil {
		return ExecResult{}, err
	}
	execID := cresp.ID

	// run it, with stdout attached
	aresp, err := apiClient.ContainerExecAttach(ctx, execID, types.ExecStartCheck{})
	if err != nil {
		return ExecResult{}, err
	}

	// buf := new(strings.Builder)
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, aresp.Reader)
	if err != nil {
		return ExecResult{}, err
	}

	// get the exit code
	iresp, err := apiClient.ContainerExecInspect(ctx, execID)
	if err != nil {
		return ExecResult{}, err
	}

	return ExecResult{exitCode: iresp.ExitCode, outBuffer: buf}, nil
}

func (r *DockerRunner) Execute(cmd ...string) (runner.Responder, error) {

	result, err := Exec(context.Background(), r.client, r.container, cmd)
	if err != nil {
		return nil, err
	}
	return result, err
}

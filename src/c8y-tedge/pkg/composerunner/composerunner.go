package composerunner

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/compose-spec/compose-go/cli"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/thin-edge/c8y-tedge/pkg/runner"
)

// reference: https://github.com/mgale/examples/blob/main/docker-compose/main.go

func listContainers() {
	out, err := exec.Command("docker", "compose", "ps").Output()
	if err != nil {
		log.Fatal(err)
	}
}

// func WithCommand() func() error {
// 	return nil
// }

func createDockerProject(ctx context.Context) *types.Project {
	// cli.WithDefaultConfigPath()

	env := make(map[string]string)
	for _, k := range os.Environ() {
		env[k] = os.Getenv(k)
	}

	cwd, _ := os.Getwd()
	configDetails := types.ConfigDetails{
		WorkingDir:  cwd,
		ConfigFiles: types.ToConfigFiles(cli.DefaultFileNames),
		Environment: env,
	}

	p, err := loader.LoadWithContext(ctx, configDetails, func(options *loader.Options) {
		// options.SetProjectName(projectName, true)
	})

	// p.GetServices()

	if err != nil {
		log.Fatalln("error load:", err)
	}
	// addServiceLabels(p)
	return p
}

// createDockerService creates a docker service which can be
// used to interact with docker-compose.
func createDockerService() (api.Service, error) {
	var srv api.Service
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return srv, err
	}

	dockerContext := "default"

	//Magic line to fix error:
	//Failed to initialize: unable to resolve docker endpoint: no context store initialized
	myOpts := &flags.ClientOptions{Context: dockerContext, LogLevel: "error"}
	err = dockerCli.Initialize(myOpts)
	if err != nil {
		return srv, err
	}

	srv = compose.NewComposeService(dockerCli)

	return srv, nil
}

func lookupContainer(ctx context.Context, srv api.Service, name string) {
	// stacks, err := srv.(ctx, api.ListOptions{

	// })

	// for _, stack := range stacks {
	// 	stack.ID
	// }
}

func myExec(ctx context.Context, srv api.Service, p *types.Project, service string, cmd ...string) {
	result, err := srv.Exec(ctx, p.Name, api.RunOptions{
		Service: service,
		Command: cmd,
		// WorkingDir:  "/",
		Tty: false,
		// Environment: []string{"TONE=test1"},
	})
	log.Println("Command result:", result, " and err:", err)
}

type ComposeRunner struct {
	client    *client.Client
	container string
}

func NewComposeRunner(target string) (*ComposeRunner, error) {

	ctx := context.TODO()

	p := createDockerProject(ctx)

	srv, err := createDockerService()

	if err != nil {
		log.Fatalln("error create docker service:", err)
	}
	myExec(ctx, srv, p, "")

	// cli, err := client.NewClientWithOpts(
	// 	client.FromEnv,
	// 	client.WithAPIVersionNegotiation(),
	// )
	if err != nil {
		panic(err)
	}
	runner := &ComposeRunner{
		client: cli,
	}
	runner.SetTarget(target)
	return runner, nil
}

func (r *ComposeRunner) SetTarget(target string) error {
	r.container = target
	return nil
}

func (r *ComposeRunner) Validate() error {
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

func (r *ComposeRunner) Execute(cmd ...string) (runner.Responder, error) {

	result, err := Exec(context.Background(), r.client, r.container, cmd)
	if err != nil {
		return nil, err
	}
	return result, err
}

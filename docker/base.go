package docker

import (
	"fmt"
	"github.com/mgutz/logxi/v1"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/network"
	"golang.org/x/net/context"
	/*
	"github.com/Sirupsen/logrus"
	*/
	"strconv"
	"os"
)

const scapeVersion = 1
const dockerVersion = "v1.22"
const dockerDaemon = "unix:///var/run/docker.sock"
const imageName = "dmp42/scape-go:testing"
const userAgent = "scape-go/1.0"

var cli *client.Client
// var log *logrus.Logger
var logger log.Logger
var ctx context.Context

func init() {
	logger = log.NewLogger(log.NewConcurrentWriter(os.Stdout), "docker")
	/*
	log = logrus.New()
	logger.Formatter = &logrus.TextFormatter{}
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	*/

	defaultHeaders := map[string]string{"User-Agent": userAgent}
	var err error
	cli, err = client.NewClient(dockerDaemon, dockerVersion, nil, defaultHeaders)
	if err != nil {
		logger.Fatal("Failed to setup docker!", "err", err)
	}
}

// Selector is a high level construct to easily select scape containers
type Selector struct {
	Name string
	URL  string
	Path string
	//	Tags string
}

// Select returns a list of scape containers matching the passed Selector
func Select(selector Selector, running bool) []types.Container {
	logger.Debug("Entering Select with selector ", selector)
	args := filters.NewArgs()

	// Filter just our containers
	args.Add("label", "com.dmp42.scape")

	// Any optional argument passed, add them
	if selector.Name != "" {
		args.Add("name", selector.Name)
	}

	/*
		if selector.Tags != "" {
			args.Add("label", fmt.Sprintf("com.dmp42.scape.tags=%s", selector.Tags))
		}
	*/

	if selector.URL != "" {
		args.Add("label", fmt.Sprintf("com.dmp42.scape.url=%s", selector.URL))
	}

	if selector.Path != "" {
		args.Add("label", fmt.Sprintf("com.dmp42.scape.path=%s", selector.Path))
	}

	options := types.ContainerListOptions{All: !running, Filter: args}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		log.Fatal("Failed to run docker list command!", "err", err)
	}
	logger.Debug("Exiting Select with containers ", containers)
	return containers
}

// Stop lets you stop a scape container
func Stop(container types.Container) {
	logger.Debug("Entering Stop container ", container)
	cli.ContainerStop(context.Background(), container.ID, 0)
	logger.Debug("Exiting Stop")
}

// Remove lets you delete a scape container
func Remove(container types.Container) {
	logger.Debug("Entering Remove with container ", container)
	cli.ContainerRemove(context.Background(), types.ContainerRemoveOptions{
		ContainerID: container.ID,
		/*
			RemoveVolumes: true,
			RemoveLinks:   true,
		*/
		Force: true,
	})
	logger.Debug("Exiting Remove")
}

// Run starts a new scape container
func Run(selector Selector) {
	logger.Debug("Entering Run with selector ", selector)
	cc := container.Config{
		Image: imageName,
		Labels: map[string]string{
			"com.dmp42.scape": strconv.Itoa(scapeVersion),
			//			"com.dmp42.scape.tags": selector.Tags,
			"com.dmp42.scape.url":  selector.URL,
			"com.dmp42.scape.path": selector.Path,
		},
		Entrypoint: []string{"bash"},
		Env: []string{},
		Tty: true,
	}
	chc := container.HostConfig{
		Binds: []string{selector.Path + ":/src"},
	}
	response, err := cli.ContainerCreate(context.Background(), &cc, &chc, &network.NetworkingConfig{}, selector.Name)
	if err != nil {
		panic(err)
	}
	err = cli.ContainerStart(context.Background(), response.ID)
	if err != nil {
		panic(err)
	}
	logger.Debug("Exiting Run")
}

// Exec runs a command inside a running scape container
func Exec(selector Selector, cmd []string) {
	logger.Debug("Entering Exec with selector ", selector, " cmd ", cmd)
	ec := types.ExecConfig{
		Container:    selector.Name,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		Cmd:          cmd,
	}

	/*		User         string   // User that will run the command
			Privileged   bool     // Is the container in privileged mode
			Tty          bool     // Attach standard streams to a tty.
			Detach       bool     // Execute in detach mode
			DetachKeys   string   // Escape keys for detach
			Cmd          []string // Execution commands and args*/

	_, err := cli.ContainerExecCreate(context.Background(), ec)
	if err != nil {
		panic(err)
	}
	logger.Debug("Exiting Exec")
}

/*	err = cli.ContainerExecStart(context.Background(), resp.ID, config types.ExecStartCheck)
	if err != nil {
		panic(err)
	}*/

/*	resp, err = cli.ContainerExecAttach(context.Background(), resp.ID, ec)
	if err != nil {
		panic(err)
	}*/

/*
func FromName(name string) types.Container{
	logger.Debug("Entering FromName with container ", name)
	c := Select(&Selector{Name: name}, true)
	if len(c) > 0 {
		return c[0]
	}
	return nil
}
*/

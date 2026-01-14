package runner

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
)

type DockerRunner struct {
	cli *client.Client
}

func NewDockerRunner() (*DockerRunner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerRunner{cli: cli}, nil
}

func (r *DockerRunner) RunStep(ctx context.Context, pipelineImage string, step domain.Step, workspace string, logWriter io.Writer) (int, error) {
	// 1. Pull Image (silently if already exists, but for now we pull)
	// In a real CI, we'd check if image exists or use a local cache
	reader, err := r.cli.ImagePull(ctx, pipelineImage, image.PullOptions{})
	if err != nil {
		return 0, err
	}
	io.Copy(io.Discard, reader) // Pull silently for now
	reader.Close()

	// 2. Prepare commands
	// Join all commands with && so they run in sequence and stop on failure
	fullCmd := ""
	for i, c := range step.Commands {
		if i > 0 {
			fullCmd += " && "
		}
		fullCmd += c
	}

	// 3. Create Container
	resp, err := r.cli.ContainerCreate(ctx, &container.Config{
		Image:      pipelineImage,
		Cmd:        []string{"sh", "-c", fullCmd},
		Env:        flattenEnv(step.Env),
		WorkingDir: "/workspace",
		Tty:        true,
	}, &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s:/workspace", workspace)},
	}, nil, nil, "")
	if err != nil {
		return 0, err
	}

	// 4. Start Container
	if err := r.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return 0, err
	}

	// 5. Stream Logs
	out, err := r.cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	if err == nil {
		go func() {
			io.Copy(logWriter, out)
			out.Close()
		}()
	}

	// 6. Wait for completion
	statusCh, errCh := r.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return 0, err
		}
	case status := <-statusCh:
		// Cleanup container
		_ = r.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return int(status.StatusCode), nil
	}

	return 0, nil
}

func flattenEnv(env map[string]string) []string {
	var res []string
	for k, v := range env {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return res
}

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

func (r *DockerRunner) RunStep(ctx context.Context, pipelineImage string, step domain.Step, workspace string) (int, error) {
	// 1. Pull Image
	reader, err := r.cli.ImagePull(ctx, pipelineImage, image.PullOptions{})
	if err != nil {
		return 0, err
	}
	io.Copy(os.Stdout, reader) // For now, log pull progress to worker stdout
	reader.Close()

	// 2. Create Container
	resp, err := r.cli.ContainerCreate(ctx, &container.Config{
		Image: pipelineImage,
		Cmd:   []string{"sh", "-c", step.Commands[0]}, // Simplified: just running the first command for now
		Env:   flattenEnv(step.Env),
		WorkingDir: "/workspace",
	}, &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s:/workspace", workspace)},
	}, nil, nil, "")
	if err != nil {
		return 0, err
	}

	// 3. Start Container
	if err := r.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return 0, err
	}

	// 4. Wait for completion
	statusCh, errCh := r.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return 0, err
		}
	case status := <-statusCh:
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

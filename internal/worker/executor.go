package worker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/princetheprogrammerbtw/nanoci/internal/domain"
	"github.com/princetheprogrammerbtw/nanoci/internal/runner"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Executor struct {
	buildRepo   domain.BuildRepository
	projectRepo domain.ProjectRepository
	runner      *runner.DockerRunner
}

func NewExecutor(br domain.BuildRepository, pr domain.ProjectRepository, r *runner.DockerRunner) *Executor {
	return &Executor{
		buildRepo:   br,
		projectRepo: pr,
		runner:      r,
	}
}

func (e *Executor) Execute(ctx context.Context, buildID string) error {
	id, err := uuid.Parse(buildID)
	if err != nil {
		return err
	}

	build, err := e.buildRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if build == nil {
		return fmt.Errorf("build not found: %s", buildID)
	}

	project, err := e.projectRepo.GetByID(ctx, build.ProjectID)
	if err != nil {
		return err
	}

	// Update build status to RUNNING
	now := time.Now()
	build.Status = domain.BuildStatusRunning
	build.StartedAt = &now
	if err := e.buildRepo.Update(ctx, build); err != nil {
		return err
	}

	// 1. Create Workspace
	workspace, err := os.MkdirTemp("", "nanoci-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workspace)

	// 2. Clone Repo
	zap.L().Info("cloning repository", zap.String("url", project.RepoURL), zap.String("commit", build.CommitHash))
	cmd := exec.CommandContext(ctx, "git", "clone", "--depth", "1", project.RepoURL, workspace)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repo: %w", err)
	}

	// 3. Parse .nanoci.yml
	pipelineFile := filepath.Join(workspace, ".nanoci.yml")
	data, err := os.ReadFile(pipelineFile)
	if err != nil {
		return fmt.Errorf("failed to read .nanoci.yml: %w", err)
	}

	var pipeline domain.Pipeline
	if err := yaml.Unmarshal(data, &pipeline); err != nil {
		return fmt.Errorf("failed to parse .nanoci.yml: %w", err)
	}

	// 4. Run Steps
	for _, step := range pipeline.Steps {
		zap.L().Info("running step", zap.String("name", step.Name))
		exitCode, err := e.runner.RunStep(ctx, pipeline.Image, step, workspace, os.Stdout)
		if err != nil {
			return e.markFailed(ctx, build, err)
		}
		if exitCode != 0 {
			return e.markFailed(ctx, build, fmt.Errorf("step %s failed with exit code %d", step.Name, exitCode))
		}
	}

	// 5. Success
	finishTime := time.Now()
	build.Status = domain.BuildStatusSuccess
	build.FinishedAt = &finishTime
	return e.buildRepo.Update(ctx, build)
}

func (e *Executor) markFailed(ctx context.Context, build *domain.Build, err error) error {
	zap.L().Error("build failed", zap.String("id", build.ID.String()), zap.Error(err))
	finishTime := time.Now()
	build.Status = domain.BuildStatusFailed
	build.FinishedAt = &finishTime
	_ = e.buildRepo.Update(ctx, build)
	return err
}

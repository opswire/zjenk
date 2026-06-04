package tool

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-jenkins/internal/client"
	"mcp-jenkins/internal/config"
	"mcp-jenkins/internal/dto"
)

type JobTools struct {
	jenkins *client.Jenkins
	cfg     *config.Config
}

func NewJobTools(j *client.Jenkins, cfg *config.Config) *JobTools {
	return &JobTools{jenkins: j, cfg: cfg}
}

// ListJobs — Out is `any` ([]dto.Job); see comment in build.go on why not []dto.Job directly.
func (t *JobTools) ListJobs(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input ListJobsInput,
) (*mcp.CallToolResult, any, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	jobs, err := t.jenkins.ListJobs(ctx, input.ProjectName)
	if err != nil {
		return toolError(fmt.Sprintf("failed to list jobs: %v", err)), nil, nil
	}
	return nil, jobs, nil
}

func (t *JobTools) GetJob(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input GetJobInput,
) (*mcp.CallToolResult, *dto.Job, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	job, err := t.jenkins.GetJob(ctx, input.JobPath)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get job: %v", err)), nil, nil
	}
	return nil, job, nil
}

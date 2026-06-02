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

// ListJobs returns all Jenkins jobs.
// Out is `any` because the MCP SDK panics when output schema type is not "object";
// slices produce type "array". The serialised value is still []dto.Job.
func (t *JobTools) ListJobs(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ struct{},
) (*mcp.CallToolResult, any, error) {
	jobs, err := t.jenkins.ListJobs(ctx)
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
	job, err := t.jenkins.GetJob(ctx, input.JobURL)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get job: %v", err)), nil, nil
	}
	return nil, job, nil
}

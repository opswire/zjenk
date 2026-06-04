package tool

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-jenkins/internal/client"
	"mcp-jenkins/internal/config"
	"mcp-jenkins/internal/dto"
)

type BuildTools struct {
	jenkins *client.Jenkins
	cfg     *config.Config
}

func NewBuildTools(j *client.Jenkins, cfg *config.Config) *BuildTools {
	return &BuildTools{jenkins: j, cfg: cfg}
}

func (t *BuildTools) GetLastBuild(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input GetLastBuildInput,
) (*mcp.CallToolResult, *dto.Build, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	build, err := t.jenkins.GetLastBuild(ctx, input.JobPath)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get last build: %v", err)), nil, nil
	}
	return nil, build, nil
}

// ListBuilds — Out is `any` ([]dto.Build); see ListJobs for the reason.
func (t *BuildTools) ListBuilds(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input ListBuildsInput,
) (*mcp.CallToolResult, any, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	builds, err := t.jenkins.ListBuilds(ctx, input.JobPath)
	if err != nil {
		return toolError(fmt.Sprintf("failed to list builds: %v", err)), nil, nil
	}
	return nil, builds, nil
}

func (t *BuildTools) GetBuild(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input BuildRefInput,
) (*mcp.CallToolResult, *dto.Build, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	build, err := t.jenkins.GetBuild(ctx, input.JobPath, input.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build: %v", err)), nil, nil
	}
	return nil, build, nil
}

func (t *BuildTools) GetBuildLog(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input GetBuildLogInput,
) (*mcp.CallToolResult, *dto.BuildLogPage, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	log, err := t.jenkins.GetBuildLog(ctx, input.JobPath, input.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build log: %v", err)), nil, nil
	}
	content, pagination := paginateText(log, input.Page, input.CharsPerPage, t.cfg.Pagination.MaxCharsPerPage)
	return nil, &dto.BuildLogPage{Content: content, Pagination: pagination}, nil
}

func (t *BuildTools) TriggerBuild(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input TriggerBuildInput,
) (*mcp.CallToolResult, *dto.TriggerResult, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	queueID, err := t.jenkins.TriggerBuild(ctx, input.JobPath, input.Params)
	if err != nil {
		return toolError(fmt.Sprintf("failed to trigger build: %v", err)), nil, nil
	}
	return nil, &dto.TriggerResult{JobPath: input.JobPath, QueueID: queueID}, nil
}

func (t *BuildTools) StopBuild(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input BuildRefInput,
) (*mcp.CallToolResult, any, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	if err := t.jenkins.StopBuild(ctx, input.JobPath, input.BuildNumber); err != nil {
		return toolError(fmt.Sprintf("failed to stop build: %v", err)), nil, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Build #%d stopped", input.BuildNumber),
			},
		},
	}, nil, nil
}

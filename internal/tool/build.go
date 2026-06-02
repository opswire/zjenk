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

// ListBuilds — Out is `any` ([]dto.Build); see ListJobs for the reason.
func (t *BuildTools) ListBuilds(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input ListBuildsInput,
) (*mcp.CallToolResult, any, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	builds, err := t.jenkins.ListBuilds(ctx, input.JobURL)
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
	build, err := t.jenkins.GetBuild(ctx, input.JobURL, input.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build: %v", err)), nil, nil
	}
	return nil, build, nil
}

// GetBuildLog returns raw console text — no structured output.
func (t *BuildTools) GetBuildLog(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input BuildRefInput,
) (*mcp.CallToolResult, any, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	log, err := t.jenkins.GetBuildLog(ctx, input.JobURL, input.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build log: %v", err)), nil, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: log}},
	}, nil, nil
}

func (t *BuildTools) TriggerBuild(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input TriggerBuildInput,
) (*mcp.CallToolResult, *dto.TriggerResult, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	queueID, err := t.jenkins.TriggerBuild(ctx, input.JobURL, input.Params)
	if err != nil {
		return toolError(fmt.Sprintf("failed to trigger build: %v", err)), nil, nil
	}
	return nil, &dto.TriggerResult{JobURL: input.JobURL, QueueID: queueID}, nil
}

func (t *BuildTools) StopBuild(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input BuildRefInput,
) (*mcp.CallToolResult, any, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	if err := t.jenkins.StopBuild(ctx, input.JobURL, input.BuildNumber); err != nil {
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

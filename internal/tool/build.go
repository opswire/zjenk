package tool

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-jenkins/internal/client"
	"mcp-jenkins/internal/config"
)

type BuildTools struct {
	jenkins *client.Jenkins
	cfg     *config.Config
}

func NewBuildTools(j *client.Jenkins, cfg *config.Config) *BuildTools {
	return &BuildTools{jenkins: j, cfg: cfg}
}

func (t *BuildTools) Register(s *mcp.Server) {
	if tc := t.cfg.Tool("jenkins_list_builds"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.listBuilds)
	}
	if tc := t.cfg.Tool("jenkins_get_build"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.getBuild)
	}
	if tc := t.cfg.Tool("jenkins_get_build_log"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.getBuildLog)
	}
	if tc := t.cfg.Tool("jenkins_trigger_build"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.triggerBuild)
	}
	if tc := t.cfg.Tool("jenkins_stop_build"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.stopBuild)
	}
}

type jobNameInput struct {
	JobName string `json:"job_name" jsonschema:"description=Jenkins job name"`
}

type buildRefInput struct {
	JobName     string `json:"job_name"     jsonschema:"description=Jenkins job name"`
	BuildNumber int64  `json:"build_number" jsonschema:"description=Build number"`
}

type triggerBuildInput struct {
	JobName string            `json:"job_name"         jsonschema:"description=Jenkins job name"`
	Params  map[string]string `json:"params,omitempty" jsonschema:"description=Optional build parameters as key-value pairs"`
}

func (t *BuildTools) listBuilds(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input jobNameInput,
) (*mcp.CallToolResult, struct{}, error) {
	if input.JobName == "" {
		return toolError("job_name is required"), struct{}{}, nil
	}
	builds, err := t.jenkins.ListBuilds(ctx, input.JobName)
	if err != nil {
		return toolError(fmt.Sprintf("failed to list builds: %v", err)), struct{}{}, nil
	}
	result, err := toolJSON(builds)
	return result, struct{}{}, err
}

func (t *BuildTools) getBuild(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input buildRefInput,
) (*mcp.CallToolResult, struct{}, error) {
	if err := validateBuildRef(input.JobName, input.BuildNumber); err != nil {
		return toolError(err.Error()), struct{}{}, nil
	}
	build, err := t.jenkins.GetBuild(ctx, input.JobName, input.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build: %v", err)), struct{}{}, nil
	}
	result, err := toolJSON(build)
	return result, struct{}{}, err
}

func (t *BuildTools) getBuildLog(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input buildRefInput,
) (*mcp.CallToolResult, struct{}, error) {
	if err := validateBuildRef(input.JobName, input.BuildNumber); err != nil {
		return toolError(err.Error()), struct{}{}, nil
	}
	log, err := t.jenkins.GetBuildLog(ctx, input.JobName, input.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build log: %v", err)), struct{}{}, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: log}},
	}, struct{}{}, nil
}

func (t *BuildTools) triggerBuild(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input triggerBuildInput,
) (*mcp.CallToolResult, struct{}, error) {
	if input.JobName == "" {
		return toolError("job_name is required"), struct{}{}, nil
	}
	queueID, err := t.jenkins.TriggerBuild(ctx, input.JobName, input.Params)
	if err != nil {
		return toolError(fmt.Sprintf("failed to trigger build: %v", err)), struct{}{}, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Build triggered for job %q. Queue ID: %d", input.JobName, queueID),
			},
		},
	}, struct{}{}, nil
}

func (t *BuildTools) stopBuild(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input buildRefInput,
) (*mcp.CallToolResult, struct{}, error) {
	if err := validateBuildRef(input.JobName, input.BuildNumber); err != nil {
		return toolError(err.Error()), struct{}{}, nil
	}
	if err := t.jenkins.StopBuild(ctx, input.JobName, input.BuildNumber); err != nil {
		return toolError(fmt.Sprintf("failed to stop build: %v", err)), struct{}{}, nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Build #%d for job %q has been stopped", input.BuildNumber, input.JobName),
			},
		},
	}, struct{}{}, nil
}

func validateBuildRef(jobName string, buildNumber int64) error {
	if jobName == "" {
		return fmt.Errorf("job_name is required")
	}
	if buildNumber <= 0 {
		return fmt.Errorf("build_number must be a positive integer")
	}
	return nil
}

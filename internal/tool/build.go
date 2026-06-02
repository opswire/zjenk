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
	JobName string `json:"job_name" jsonschema:"Jenkins job name"`
}

type buildRefInput struct {
	JobName     string `json:"job_name"     jsonschema:"Jenkins job name"`
	BuildNumber int64  `json:"build_number" jsonschema:"Build number"`
}

type triggerBuildInput struct {
	JobName string            `json:"job_name"         jsonschema:"Jenkins job name"`
	Params  map[string]string `json:"params,omitempty" jsonschema:"Optional build parameters as key-value pairs"`
}

func (t *BuildTools) listBuilds(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[jobNameInput],
) (*mcp.CallToolResultFor[struct{}], error) {
	if params.Arguments.JobName == "" {
		return toolError("job_name is required"), nil
	}

	builds, err := t.jenkins.ListBuilds(ctx, params.Arguments.JobName)
	if err != nil {
		return toolError(fmt.Sprintf("failed to list builds: %v", err)), nil
	}
	return toolJSON(builds)
}

func (t *BuildTools) getBuild(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[buildRefInput],
) (*mcp.CallToolResultFor[struct{}], error) {
	if err := validateBuildRef(params.Arguments.JobName, params.Arguments.BuildNumber); err != nil {
		return toolError(err.Error()), nil
	}

	build, err := t.jenkins.GetBuild(ctx, params.Arguments.JobName, params.Arguments.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build: %v", err)), nil
	}
	return toolJSON(build)
}

func (t *BuildTools) getBuildLog(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[buildRefInput],
) (*mcp.CallToolResultFor[struct{}], error) {
	if err := validateBuildRef(params.Arguments.JobName, params.Arguments.BuildNumber); err != nil {
		return toolError(err.Error()), nil
	}

	log, err := t.jenkins.GetBuildLog(ctx, params.Arguments.JobName, params.Arguments.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build log: %v", err)), nil
	}
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{&mcp.TextContent{Text: log}},
	}, nil
}

func (t *BuildTools) triggerBuild(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[triggerBuildInput],
) (*mcp.CallToolResultFor[struct{}], error) {
	if params.Arguments.JobName == "" {
		return toolError("job_name is required"), nil
	}

	queueID, err := t.jenkins.TriggerBuild(ctx, params.Arguments.JobName, params.Arguments.Params)
	if err != nil {
		return toolError(fmt.Sprintf("failed to trigger build: %v", err)), nil
	}
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Build triggered for job %q. Queue ID: %d", params.Arguments.JobName, queueID),
			},
		},
	}, nil
}

func (t *BuildTools) stopBuild(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[buildRefInput],
) (*mcp.CallToolResultFor[struct{}], error) {
	if err := validateBuildRef(params.Arguments.JobName, params.Arguments.BuildNumber); err != nil {
		return toolError(err.Error()), nil
	}

	if err := t.jenkins.StopBuild(ctx, params.Arguments.JobName, params.Arguments.BuildNumber); err != nil {
		return toolError(fmt.Sprintf("failed to stop build: %v", err)), nil
	}
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: fmt.Sprintf("Build #%d for job %q has been stopped", params.Arguments.BuildNumber, params.Arguments.JobName),
			},
		},
	}, nil
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

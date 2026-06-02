package tool

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-jenkins/internal/client"
	"mcp-jenkins/internal/config"
)

type JobTools struct {
	jenkins *client.Jenkins
	cfg     *config.Config
}

func NewJobTools(j *client.Jenkins, cfg *config.Config) *JobTools {
	return &JobTools{jenkins: j, cfg: cfg}
}

func (t *JobTools) Register(s *mcp.Server) {
	if tc := t.cfg.Tool("jenkins_list_jobs"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.listJobs)
	}
	if tc := t.cfg.Tool("jenkins_get_job"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.getJob)
	}
}

type getJobInput struct {
	JobName string `json:"job_name" jsonschema:"Jenkins job name"`
}

func (t *JobTools) listJobs(
	ctx context.Context,
	_ *mcp.ServerSession,
	_ *mcp.CallToolParamsFor[struct{}],
) (*mcp.CallToolResultFor[struct{}], error) {
	jobs, err := t.jenkins.ListJobs(ctx)
	if err != nil {
		return toolError(fmt.Sprintf("failed to list jobs: %v", err)), nil
	}
	return toolJSON(jobs)
}

func (t *JobTools) getJob(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[getJobInput],
) (*mcp.CallToolResultFor[struct{}], error) {
	if params.Arguments.JobName == "" {
		return toolError("job_name is required"), nil
	}

	job, err := t.jenkins.GetJob(ctx, params.Arguments.JobName)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get job %q: %v", params.Arguments.JobName, err)), nil
	}
	return toolJSON(job)
}

func toolJSON(v any) (*mcp.CallToolResultFor[struct{}], error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil
}

func toolError(msg string) *mcp.CallToolResultFor[struct{}] {
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}
}

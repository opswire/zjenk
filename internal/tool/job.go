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
	JobName string `json:"job_name" jsonschema:"description=Jenkins job name"`
}

func (t *JobTools) listJobs(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ struct{},
) (*mcp.CallToolResult, struct{}, error) {
	jobs, err := t.jenkins.ListJobs(ctx)
	if err != nil {
		return toolError(fmt.Sprintf("failed to list jobs: %v", err)), struct{}{}, nil
	}
	result, err := toolJSON(jobs)
	return result, struct{}{}, err
}

func (t *JobTools) getJob(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input getJobInput,
) (*mcp.CallToolResult, struct{}, error) {
	if input.JobName == "" {
		return toolError("job_name is required"), struct{}{}, nil
	}

	job, err := t.jenkins.GetJob(ctx, input.JobName)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get job %q: %v", input.JobName, err)), struct{}{}, nil
	}
	result, err := toolJSON(job)
	return result, struct{}{}, err
}

func toolJSON(v any) (*mcp.CallToolResult, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal result: %w", err)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(data)}},
	}, nil
}

func toolError(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}
}

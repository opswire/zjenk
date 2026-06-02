package tool

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-jenkins/internal/client"
	"mcp-jenkins/internal/config"
)

type SearchTools struct {
	jenkins *client.Jenkins
	cfg     *config.Config
}

func NewSearchTools(j *client.Jenkins, cfg *config.Config) *SearchTools {
	return &SearchTools{jenkins: j, cfg: cfg}
}

func (t *SearchTools) Register(s *mcp.Server) {
	if tc := t.cfg.Tool("jenkins_search_log"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.searchLog)
	}
}

type searchLogInput struct {
	JobName     string `json:"job_name"     jsonschema:"description=Jenkins job name"`
	BuildNumber int64  `json:"build_number" jsonschema:"description=Build number"`
	Pattern     string `json:"pattern"      jsonschema:"description=Regex pattern or substring to search for in the build log"`
}

type searchMatch struct {
	Line    int    `json:"line"`
	Content string `json:"content"`
}

func (t *SearchTools) searchLog(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input searchLogInput,
) (*mcp.CallToolResult, struct{}, error) {
	if err := validateBuildRef(input.JobName, input.BuildNumber); err != nil {
		return toolError(err.Error()), struct{}{}, nil
	}
	if input.Pattern == "" {
		return toolError("pattern is required"), struct{}{}, nil
	}

	log, err := t.jenkins.GetBuildLog(ctx, input.JobName, input.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build log: %v", err)), struct{}{}, nil
	}

	re, reErr := regexp.Compile(input.Pattern)

	var matches []searchMatch
	for i, line := range strings.Split(log, "\n") {
		var matched bool
		if reErr != nil {
			matched = strings.Contains(line, input.Pattern)
		} else {
			matched = re.MatchString(line)
		}
		if matched {
			matches = append(matches, searchMatch{Line: i + 1, Content: line})
		}
	}

	result, err := toolJSON(matches)
	return result, struct{}{}, err
}

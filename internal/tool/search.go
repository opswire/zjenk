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
	JobName     string `json:"job_name"     jsonschema:"Jenkins job name"`
	BuildNumber int64  `json:"build_number" jsonschema:"Build number"`
	Pattern     string `json:"pattern"      jsonschema:"Regex pattern or substring to search for in the build log"`
}

type searchMatch struct {
	Line    int    `json:"line"`
	Content string `json:"content"`
}

func (t *SearchTools) searchLog(
	ctx context.Context,
	_ *mcp.ServerSession,
	params *mcp.CallToolParamsFor[searchLogInput],
) (*mcp.CallToolResultFor[struct{}], error) {
	if err := validateBuildRef(params.Arguments.JobName, params.Arguments.BuildNumber); err != nil {
		return toolError(err.Error()), nil
	}
	if params.Arguments.Pattern == "" {
		return toolError("pattern is required"), nil
	}

	log, err := t.jenkins.GetBuildLog(ctx, params.Arguments.JobName, params.Arguments.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build log: %v", err)), nil
	}

	re, reErr := regexp.Compile(params.Arguments.Pattern)

	var matches []searchMatch
	for i, line := range strings.Split(log, "\n") {
		var matched bool
		if reErr != nil {
			matched = strings.Contains(line, params.Arguments.Pattern)
		} else {
			matched = re.MatchString(line)
		}
		if matched {
			matches = append(matches, searchMatch{Line: i + 1, Content: line})
		}
	}

	return toolJSON(matches)
}

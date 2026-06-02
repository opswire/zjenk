package tool

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-jenkins/internal/client"
	"mcp-jenkins/internal/config"
	"mcp-jenkins/internal/dto"
)

type SearchTools struct {
	jenkins *client.Jenkins
	cfg     *config.Config
}

func NewSearchTools(j *client.Jenkins, cfg *config.Config) *SearchTools {
	return &SearchTools{jenkins: j, cfg: cfg}
}

// SearchLog — Out is `any` ([]dto.SearchMatch); see ListJobs for the reason.
func (t *SearchTools) SearchLog(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input SearchLogInput,
) (*mcp.CallToolResult, any, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}

	log, err := t.jenkins.GetBuildLog(ctx, input.JobURL, input.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build log: %v", err)), nil, nil
	}

	re, reErr := regexp.Compile(input.Pattern)

	var matches []dto.SearchMatch
	for i, line := range strings.Split(log, "\n") {
		var matched bool
		if reErr != nil {
			matched = strings.Contains(line, input.Pattern)
		} else {
			matched = re.MatchString(line)
		}
		if matched {
			matches = append(matches, dto.SearchMatch{Line: i + 1, Content: line})
		}
	}

	return nil, matches, nil
}

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

func (t *SearchTools) SearchLog(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input SearchLogInput,
) (*mcp.CallToolResult, *dto.SearchPage, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}

	rawLog, err := t.jenkins.GetBuildLog(ctx, input.JobPath, input.BuildNumber)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get build log: %v", err)), nil, nil
	}

	re, reErr := regexp.Compile(input.Pattern)
	lines := strings.Split(rawLog, "\n")

	var matchIndices []int
	for i, line := range lines {
		var matched bool
		if reErr != nil {
			matched = strings.Contains(line, input.Pattern)
		} else {
			matched = re.MatchString(line)
		}
		if matched {
			matchIndices = append(matchIndices, i)
		}
	}

	matches := applyContext(lines, matchIndices, input.ContextLines, t.cfg.Pagination.MaxContextLines)
	page, pagination := paginateMatches(matches, input.Page, t.cfg.Pagination.MaxSearchResultsPerPage)

	return nil, &dto.SearchPage{Matches: page, Pagination: pagination}, nil
}

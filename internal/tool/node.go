package tool

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-jenkins/internal/client"
)

type NodeTools struct {
	jenkins *client.Jenkins
}

func NewNodeTools(j *client.Jenkins) *NodeTools {
	return &NodeTools{jenkins: j}
}

// ListNodes — Out is `any` ([]dto.Node); see ListJobs for the reason.
func (t *NodeTools) ListNodes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input ListNodesInput,
) (*mcp.CallToolResult, any, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	nodes, err := t.jenkins.ListNodes(ctx, input.ProjectName)
	if err != nil {
		return toolError(fmt.Sprintf("failed to list nodes: %v", err)), nil, nil
	}
	return nil, nodes, nil
}

// GetQueue — Out is `any` ([]dto.QueueItem); see ListJobs for the reason.
func (t *NodeTools) GetQueue(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input GetQueueInput,
) (*mcp.CallToolResult, any, error) {
	if err := input.Validate(); err != nil {
		return toolError(err.Error()), nil, nil
	}
	items, err := t.jenkins.GetQueue(ctx, input.ProjectName)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get queue: %v", err)), nil, nil
	}
	return nil, items, nil
}

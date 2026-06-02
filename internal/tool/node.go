package tool

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"mcp-jenkins/internal/client"
	"mcp-jenkins/internal/config"
)

type NodeTools struct {
	jenkins *client.Jenkins
	cfg     *config.Config
}

func NewNodeTools(j *client.Jenkins, cfg *config.Config) *NodeTools {
	return &NodeTools{jenkins: j, cfg: cfg}
}

func (t *NodeTools) Register(s *mcp.Server) {
	if tc := t.cfg.Tool("jenkins_list_nodes"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.listNodes)
	}
	if tc := t.cfg.Tool("jenkins_get_queue"); tc.IsEnabled {
		mcp.AddTool(s, &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}, t.getQueue)
	}
}

func (t *NodeTools) listNodes(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ struct{},
) (*mcp.CallToolResult, struct{}, error) {
	nodes, err := t.jenkins.ListNodes(ctx)
	if err != nil {
		return toolError(fmt.Sprintf("failed to list nodes: %v", err)), struct{}{}, nil
	}
	result, err := toolJSON(nodes)
	return result, struct{}{}, err
}

func (t *NodeTools) getQueue(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	_ struct{},
) (*mcp.CallToolResult, struct{}, error) {
	items, err := t.jenkins.GetQueue(ctx)
	if err != nil {
		return toolError(fmt.Sprintf("failed to get queue: %v", err)), struct{}{}, nil
	}
	result, err := toolJSON(items)
	return result, struct{}{}, err
}

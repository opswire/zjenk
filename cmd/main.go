package main

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"mcp-jenkins/internal/client"
	"mcp-jenkins/internal/config"
	"mcp-jenkins/internal/tool"
	"mcp-jenkins/server"
)

func main() {
	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),
		fx.Provide(
			newLogger,
			config.Load,
			client.New,
			tool.NewJobTools,
			tool.NewBuildTools,
			tool.NewNodeTools,
			tool.NewSearchTools,
			server.New,
		),
		fx.Invoke(registerTools),
		fx.Invoke(server.Run),
	).Run()
}

func registerTools(
	s *mcp.Server,
	cfg *config.Config,
	jobs *tool.JobTools,
	builds *tool.BuildTools,
	nodes *tool.NodeTools,
	search *tool.SearchTools,
) {
	type entry struct {
		key string
		fn  func()
	}

	for _, e := range []entry{
		{"list_jobs", func() { mcp.AddTool(s, toolDef(cfg, "list_jobs"), jobs.ListJobs) }},
		{"get_job", func() { mcp.AddTool(s, toolDef(cfg, "get_job"), jobs.GetJob) }},
		{"get_last_build", func() { mcp.AddTool(s, toolDef(cfg, "get_last_build"), builds.GetLastBuild) }},
		{"list_builds", func() { mcp.AddTool(s, toolDef(cfg, "list_builds"), builds.ListBuilds) }},
		{"get_build", func() { mcp.AddTool(s, toolDef(cfg, "get_build"), builds.GetBuild) }},
		{"get_build_log", func() { mcp.AddTool(s, toolDef(cfg, "get_build_log"), builds.GetBuildLog) }},
		{"trigger_build", func() { mcp.AddTool(s, toolDef(cfg, "trigger_build"), builds.TriggerBuild) }},
		{"stop_build", func() { mcp.AddTool(s, toolDef(cfg, "stop_build"), builds.StopBuild) }},
		{"list_nodes", func() { mcp.AddTool(s, toolDef(cfg, "list_nodes"), nodes.ListNodes) }},
		{"get_queue", func() { mcp.AddTool(s, toolDef(cfg, "get_queue"), nodes.GetQueue) }},
		{"search_log", func() { mcp.AddTool(s, toolDef(cfg, "search_log"), search.SearchLog) }},
	} {
		if cfg.Tool(e.key).IsEnabled {
			e.fn()
		}
	}
}

// toolDef builds a *mcp.Tool from config for the given key.
func toolDef(cfg *config.Config, key string) *mcp.Tool {
	tc := cfg.Tool(key)
	return &mcp.Tool{Name: tc.Name, Description: tc.Description.EN}
}

func newLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

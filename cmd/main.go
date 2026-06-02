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
		{"jenkins_list_jobs", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_list_jobs"), jobs.ListJobs) }},
		{"jenkins_get_job", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_get_job"), jobs.GetJob) }},
		{"jenkins_list_builds", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_list_builds"), builds.ListBuilds) }},
		{"jenkins_get_build", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_get_build"), builds.GetBuild) }},
		{"jenkins_get_build_log", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_get_build_log"), builds.GetBuildLog) }},
		{"jenkins_trigger_build", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_trigger_build"), builds.TriggerBuild) }},
		{"jenkins_stop_build", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_stop_build"), builds.StopBuild) }},
		{"jenkins_list_nodes", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_list_nodes"), nodes.ListNodes) }},
		{"jenkins_get_queue", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_get_queue"), nodes.GetQueue) }},
		{"jenkins_search_log", func() { mcp.AddTool(s, toolDef(cfg, "jenkins_search_log"), search.SearchLog) }},
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

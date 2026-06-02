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
		fx.Invoke(func(
			s *mcp.Server,
			jobs *tool.JobTools,
			builds *tool.BuildTools,
			nodes *tool.NodeTools,
			search *tool.SearchTools,
		) {
			jobs.Register(s)
			builds.Register(s)
			nodes.Register(s)
			search.Register(s)
		}),
		fx.Invoke(server.Run),
	).Run()
}

func newLogger() (*zap.Logger, error) {
	return zap.NewProduction()
}

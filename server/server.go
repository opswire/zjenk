package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"mcp-jenkins/internal/config"
)

func New(cfg *config.Config, log *zap.Logger) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name:    cfg.MCP.Name,
		Version: cfg.MCP.Version,
	}, nil)

	log.Info("MCP server configured",
		zap.String("name", cfg.MCP.Name),
		zap.String("version", cfg.MCP.Version),
		zap.String("jenkins_url", cfg.Jenkins.URL),
	)

	return s
}

func Run(lc fx.Lifecycle, s *mcp.Server, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Info("starting MCP server on stdio")
				if err := s.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
					log.Error("MCP server stopped with error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("MCP server shutting down")
			return nil
		},
	})
}

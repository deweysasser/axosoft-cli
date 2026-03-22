package program

import (
	"github.com/deweysasser/axosoft-cli/tools"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

// ServeCmd starts the MCP server via stdio transport
type ServeCmd struct{}

func (cmd *ServeCmd) Run(options *Options) error {
	client, err := newClient(options)
	if err != nil {
		return err
	}

	log.Info().Str("base_url", options.BaseURL).Msg("Starting MCP server")

	s := server.NewMCPServer(
		"axosoft",
		Version,
		server.WithToolCapabilities(false),
	)

	tools.RegisterAll(s, client)

	return server.ServeStdio(s)
}

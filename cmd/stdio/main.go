package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"tutorial/mcp"

	"github.com/mark3labs/mcp-go/server"
)

const version = "1.0.0"

func main() {
	// Create logger that outputs to stderr (so it doesn't interfere with stdio transport)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	mcpServer := server.NewMCPServer(
		"tutorial-mcp-server",
		version,
		server.WithLogging(),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	mcpServer.AddTools(
		mcp.CalculatorTool(),
		mcp.SystemInfoTool(),
	)

	mcpServer.AddPrompts(
		mcp.MathTutorPrompt(),
		mcp.CodeReviewPrompt(),
	)

	mcpServer.AddResources(
		mcp.SystemStatusResource(),
		mcp.MathConstantsResource(),
	)

	stdioServer := server.NewStdioServer(mcpServer)

	errChan := make(chan error, 1)
	go func() {
		errChan <- stdioServer.Listen(ctx, os.Stdin, os.Stdout)
	}()

	logger.Info("Tutorial MCP Server started", "version", version, "transport", "stdio")

	select {
	case <-ctx.Done():
		logger.Info("Tutorial MCP Server stopped")
	case err := <-errChan:
		if err != nil {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}
}

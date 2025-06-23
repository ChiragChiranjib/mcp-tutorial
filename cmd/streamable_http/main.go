package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"tutorial/mcp"

	"github.com/mark3labs/mcp-go/server"
)

const version = "1.0.0"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	port := 8081
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	mcpServer := server.NewMCPServer(
		"tutorial-mcp-server",
		version,
		server.WithLogging(),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
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

	httpServer := server.NewStreamableHTTPServer(
		mcpServer,
		server.WithStateLess(true),
	)

	errChan := make(chan error, 1)
	go func() {
		errChan <- httpServer.Start(fmt.Sprintf(":%d", port))
	}()

	logger.Info("Tutorial MCP Server started", "version", version, "transport", "streamable_http", "port", port)

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

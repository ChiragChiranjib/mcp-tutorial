package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	port := 8080
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

	sseServer := server.NewSSEServer(
		mcpServer,
		server.WithKeepAlive(true),
		server.WithKeepAliveInterval(10*time.Second),
	)

	errChan := make(chan error, 1)
	go func() {
		errChan <- sseServer.Start(fmt.Sprintf(":%d", port))
	}()

	logger.Info("Tutorial MCP Server started", "version", version, "transport", "sse", "port", port)

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

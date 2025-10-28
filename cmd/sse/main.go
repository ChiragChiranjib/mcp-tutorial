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

	"reconsaas/mcp"

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
		"reconsaas-mcp-server",
		version,
		server.WithLogging(),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)

	mcpServer.AddTools(
		mcp.CalculatorTool(),
		mcp.SystemInfoTool(),
		mcp.ReconFileAnalysisTool(),
		mcp.ReconMasterSourceTool(),
		mcp.ReconMerchantSourceTool(),
		mcp.ReconStateRuleTool(),
		mcp.ReconProcessSetupTool(),
		mcp.ReconAggregationTool(),
	)

	mcpServer.AddPrompts(
		mcp.MathTutorPrompt(),
		mcp.CodeReviewPrompt(),
		mcp.ReconFileAnalysisPrompt(),
		mcp.ReconMasterSourcePrompt(),
		mcp.ReconMerchantSourcePrompt(),
		mcp.ReconStateRulePrompt(),
		mcp.ReconProcessSetupPrompt(),
		mcp.ReconAggregationPrompt(),
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

	logger.Info("ReconSaas MCP Server started", "version", version, "transport", "sse", "port", port)

	select {
	case <-ctx.Done():
		logger.Info("ReconSaas MCP Server stopped")
	case err := <-errChan:
		if err != nil {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}
}

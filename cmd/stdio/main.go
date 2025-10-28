package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"reconsaas/mcp"

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
		// Recon-SaaS tools
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
		// Recon-SaaS prompts
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

	stdioServer := server.NewStdioServer(mcpServer)

	errChan := make(chan error, 1)
	go func() {
		errChan <- stdioServer.Listen(ctx, os.Stdin, os.Stdout)
	}()

	logger.Info("ReconSaas MCP Server started", "version", version, "transport", "stdio")

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

package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// SystemStatusResource System status resource providing server status information
func SystemStatusResource() server.ServerResource {
	resource := mcp.NewResource(
		"system://status",
		"System Status",
		mcp.WithResourceDescription("Provides current system status and server information"),
		mcp.WithMIMEType("application/json"),
	)

	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		now := time.Now()

		status := map[string]interface{}{
			"timestamp":   now.Format(time.RFC3339),
			"server_name": "Tutorial MCP Server",
			"version":     "1.0.0",
			"status":      "operational",
			"uptime_info": map[string]interface{}{
				"current_time": now.Format("2006-01-02 15:04:05 MST"),
				"unix_time":    now.Unix(),
			},
			"capabilities": []string{
				"calculator",
				"system_info",
			},
		}

		content, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal status: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(content),
			},
		}, nil
	}

	return server.ServerResource{
		Resource: resource,
		Handler:  handler,
	}
}

// MathConstantsResource Math constants resource with common mathematical constants
func MathConstantsResource() server.ServerResource {
	resource := mcp.NewResource(
		"math://constants",
		"Mathematical Constants",
		mcp.WithResourceDescription("Common mathematical constants and their values"),
		mcp.WithMIMEType("application/json"),
	)

	handler := func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		constants := map[string]interface{}{
			"pi": map[string]interface{}{
				"symbol":      "π",
				"value":       3.141592653589793,
				"description": "The ratio of a circle's circumference to its diameter",
			},
			"e": map[string]interface{}{
				"symbol":      "e",
				"value":       2.718281828459045,
				"description": "Euler's number, the base of natural logarithm",
			},
			"phi": map[string]interface{}{
				"symbol":      "φ",
				"value":       1.618033988749895,
				"description": "The golden ratio",
			},
			"sqrt2": map[string]interface{}{
				"symbol":      "√2",
				"value":       1.4142135623730951,
				"description": "The square root of 2",
			},
		}

		content, err := json.MarshalIndent(constants, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal constants: %w", err)
		}

		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      request.Params.URI,
				MIMEType: "application/json",
				Text:     string(content),
			},
		}, nil
	}

	return server.ServerResource{
		Resource: resource,
		Handler:  handler,
	}
}

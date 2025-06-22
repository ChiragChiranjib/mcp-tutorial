package mcp

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// CalculatorTool Calculator tool for basic math operations
func CalculatorTool() server.ServerTool {
	tool := mcp.NewTool("calculator",
		mcp.WithDescription("Perform basic mathematical calculations"),
		mcp.WithString("operation",
			mcp.Description("The mathematical operation to perform"),
			mcp.Required(),
			mcp.Enum("add", "subtract", "multiply", "divide", "power", "sqrt"),
		),
		mcp.WithNumber("first_number",
			mcp.Description("The first number for the operation"),
			mcp.Required(),
		),
		mcp.WithNumber("second_number",
			mcp.Description("The second number (not required for sqrt)"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		operation, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		firstNum, err := request.RequireFloat("first_number")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var result float64

		switch operation {
		case "add":
			secondNum, err := request.RequireFloat("second_number")
			if err != nil {
				return mcp.NewToolResultError("second_number is required for addition"), nil
			}
			result = firstNum + secondNum
		case "subtract":
			secondNum, err := request.RequireFloat("second_number")
			if err != nil {
				return mcp.NewToolResultError("second_number is required for subtraction"), nil
			}
			result = firstNum - secondNum
		case "multiply":
			secondNum, err := request.RequireFloat("second_number")
			if err != nil {
				return mcp.NewToolResultError("second_number is required for multiplication"), nil
			}
			result = firstNum * secondNum
		case "divide":
			secondNum, err := request.RequireFloat("second_number")
			if err != nil {
				return mcp.NewToolResultError("second_number is required for division"), nil
			}
			if secondNum == 0 {
				return mcp.NewToolResultError("cannot divide by zero"), nil
			}
			result = firstNum / secondNum
		case "power":
			secondNum, err := request.RequireFloat("second_number")
			if err != nil {
				return mcp.NewToolResultError("second_number is required for power operation"), nil
			}
			result = math.Pow(firstNum, secondNum)
		case "sqrt":
			if firstNum < 0 {
				return mcp.NewToolResultError("cannot calculate square root of negative number"), nil
			}
			result = math.Sqrt(firstNum)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown operation: %s", operation)), nil
		}

		// Format the result
		var resultStr string
		if operation == "sqrt" {
			resultStr = fmt.Sprintf("√%.2f = %.6f", firstNum, result)
		} else {
			secondNum, _ := request.RequireFloat("second_number")
			var operatorSymbol string
			switch operation {
			case "add":
				operatorSymbol = "+"
			case "subtract":
				operatorSymbol = "-"
			case "multiply":
				operatorSymbol = "×"
			case "divide":
				operatorSymbol = "÷"
			case "power":
				operatorSymbol = "^"
			}
			resultStr = fmt.Sprintf("%.2f %s %.2f = %.6f", firstNum, operatorSymbol, secondNum, result)
		}

		return mcp.NewToolResultText(resultStr), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// SystemInfoTool System info tool for time and date information
func SystemInfoTool() server.ServerTool {
	tool := mcp.NewTool("system_info",
		mcp.WithDescription("Get system information like current time and date"),
		mcp.WithString("info_type",
			mcp.Description("Type of system information to retrieve"),
			mcp.Required(),
			mcp.Enum("time", "date", "datetime"),
		),
		mcp.WithString("format",
			mcp.Description("Format for the output"),
			mcp.Enum("iso", "rfc3339", "unix", "human"),
			mcp.DefaultString("human"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		infoType, err := request.RequireString("info_type")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		format := request.GetString("format", "human")

		now := time.Now()
		var result string

		switch infoType {
		case "time":
			switch format {
			case "iso":
				result = now.Format("15:04:05")
			case "rfc3339":
				result = now.Format(time.RFC3339)
			case "unix":
				result = strconv.FormatInt(now.Unix(), 10)
			case "human":
				result = now.Format("3:04:05 PM MST")
			}
		case "date":
			switch format {
			case "iso":
				result = now.Format("2006-01-02")
			case "rfc3339":
				result = now.Format(time.RFC3339)
			case "unix":
				result = strconv.FormatInt(now.Unix(), 10)
			case "human":
				result = now.Format("Monday, January 2, 2006")
			}
		case "datetime":
			switch format {
			case "iso":
				result = now.Format("2006-01-02T15:04:05")
			case "rfc3339":
				result = now.Format(time.RFC3339)
			case "unix":
				result = strconv.FormatInt(now.Unix(), 10)
			case "human":
				result = now.Format("Monday, January 2, 2006 at 3:04:05 PM MST")
			}
		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown info_type: %s", infoType)), nil
		}

		return mcp.NewToolResultText(result), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

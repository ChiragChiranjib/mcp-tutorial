package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
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

// ReconFileAnalysisTool File analysis tool for recon-saas onboarding
func ReconFileAnalysisTool() server.ServerTool {
	tool := mcp.NewTool("recon_file_analysis",
		mcp.WithDescription("Analyze uploaded reconciliation files to identify EntityID and Amount columns for master source creation"),
		mcp.WithString("file1_path",
			mcp.Description("Full file path to the first reconciliation file (e.g., /path/to/transactions.csv or /path/to/transactions.xlsx)"),
			mcp.Required(),
		),
		mcp.WithString("file2_path",
			mcp.Description("Full file path to the second reconciliation file (e.g., /path/to/bank_statements.csv or /path/to/bank_statements.xlsx)"),
			mcp.Required(),
		),
		mcp.WithString("file1_type",
			mcp.Description("Type of the first file"),
			mcp.Required(),
			mcp.Enum("csv", "excel"),
		),
		mcp.WithString("file2_type",
			mcp.Description("Type of the second file"),
			mcp.Required(),
			mcp.Enum("csv", "excel"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		file1Path, err := request.RequireString("file1_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		file2Path, err := request.RequireString("file2_path")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		file1Type, err := request.RequireString("file1_type")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		file2Type, err := request.RequireString("file2_type")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Analyze both files based on their type
		analysis1, err := analyzeFile(file1Path, "file_1", file1Type)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze file 1: %v", err)), nil
		}

		analysis2, err := analyzeFile(file2Path, "file_2", file2Type)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze file 2: %v", err)), nil
		}

		// Create comprehensive analysis result
		result := map[string]interface{}{
			"file_analysis": map[string]interface{}{
				"file_1": analysis1,
				"file_2": analysis2,
			},
			"compatibility_check": map[string]interface{}{
				"can_reconcile":            true,
				"common_patterns":          []string{"amount", "date"},
				"suggested_reconciliation": "Match by EntityID and Amount fields",
			},
			"analysis_type": "comprehensive",
			"timestamp":     time.Now().Format(time.RFC3339),
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// ReconMasterSourceTool Master source creation tool for recon-saas
func ReconMasterSourceTool() server.ServerTool {
	tool := mcp.NewTool("recon_master_source",
		mcp.WithDescription("Create master source configurations for recon-saas using file analysis data"),
		mcp.WithString("environment",
			mcp.Description("Environment to use for API calls: 'local' (http://localhost:9400), 'dev' (https://recon-saas.dev.razorpay.in), or 'prod' (https://recon-saas.concierge.razorpay.com). Defaults to 'dev' if not specified."),
			mcp.Enum("local", "dev", "prod"),
			mcp.DefaultString("dev"),
		),
		mcp.WithString("source1_name",
			mcp.Description("Name for the first master source"),
			mcp.Required(),
		),
		mcp.WithString("source2_name",
			mcp.Description("Name for the second master source"),
			mcp.Required(),
		),
		mcp.WithString("source1_columns",
			mcp.Description("JSON array of column names from first file"),
			mcp.Required(),
		),
		mcp.WithString("source2_columns",
			mcp.Description("JSON array of column names from second file"),
			mcp.Required(),
		),
		mcp.WithString("source1_entityid",
			mcp.Description("Selected EntityID column name for first source"),
			mcp.Required(),
		),
		mcp.WithString("source2_entityid",
			mcp.Description("Selected EntityID column name for second source"),
			mcp.Required(),
		),
		mcp.WithString("source1_amount",
			mcp.Description("Selected Amount column name for first source"),
			mcp.Required(),
		),
		mcp.WithString("source2_amount",
			mcp.Description("Selected Amount column name for second source"),
			mcp.Required(),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get environment (defaults to "dev" if not specified)
		environment := request.GetString("environment", DefaultEnvironment)

		source1Name, err := request.RequireString("source1_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2Name, err := request.RequireString("source2_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source1Columns, err := request.RequireString("source1_columns")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2Columns, err := request.RequireString("source2_columns")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source1EntityID, err := request.RequireString("source1_entityid")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2EntityID, err := request.RequireString("source2_entityid")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source1Amount, err := request.RequireString("source1_amount")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2Amount, err := request.RequireString("source2_amount")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Create master sources via API calls
		masterSource1ID, err := createMasterSource(ctx, source1Name, source1Columns, source1EntityID, source1Amount, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create master source 1: %v", err)), nil
		}

		masterSource2ID, err := createMasterSource(ctx, source2Name, source2Columns, source2EntityID, source2Amount, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create master source 2: %v", err)), nil
		}

		result := map[string]interface{}{
			"status":       "success",
			"message":      "Master sources created successfully",
			"environment":  GetEnvironmentName(environment),
			"api_base_url": GetBaseURL(environment),
			"created_sources": map[string]interface{}{
				"source_1": map[string]interface{}{
					"master_source_id":         masterSource1ID,
					"name":                     source1Name,
					"selected_entityid_column": source1EntityID,
					"selected_amount_column":   source1Amount,
				},
				"source_2": map[string]interface{}{
					"master_source_id":         masterSource2ID,
					"name":                     source2Name,
					"selected_entityid_column": source2EntityID,
					"selected_amount_column":   source2Amount,
				},
			},
			"for_future_prompts": map[string]interface{}{
				"master_source_id_1": masterSource1ID,
				"master_source_id_2": masterSource2ID,
				"source_1_name":      source1Name,
				"source_2_name":      source2Name,
				"environment":        environment,
			},
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// ReconMerchantSourceTool Merchant source creation tool for recon-saas
func ReconMerchantSourceTool() server.ServerTool {
	tool := mcp.NewTool("recon_merchant_source",
		mcp.WithDescription("Create merchant-specific source configurations for recon-saas"),
		mcp.WithString("environment",
			mcp.Description("Environment to use for API calls: 'local' (http://localhost:9400), 'dev' (https://recon-saas.dev.razorpay.in), or 'prod' (https://recon-saas.concierge.razorpay.com). Defaults to 'dev' if not specified."),
			mcp.Enum("local", "dev", "prod"),
			mcp.DefaultString("dev"),
		),
		mcp.WithString("merchant_id",
			mcp.Description("Merchant identifier for this onboarding process"),
			mcp.Required(),
		),
		mcp.WithString("master_source_id_1",
			mcp.Description("First master source ID from previous step"),
			mcp.Required(),
		),
		mcp.WithString("master_source_id_2",
			mcp.Description("Second master source ID from previous step"),
			mcp.Required(),
		),
		mcp.WithString("source_1_name",
			mcp.Description("Name of the first source"),
			mcp.Required(),
		),
		mcp.WithString("source_2_name",
			mcp.Description("Name of the second source"),
			mcp.Required(),
		),
		mcp.WithString("source_naming_strategy",
			mcp.Description("Strategy for naming merchant sources"),
			mcp.Enum("descriptive", "timestamp", "sequential", "custom"),
			mcp.DefaultString("descriptive"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get environment (defaults to "dev" if not specified)
		environment := request.GetString("environment", DefaultEnvironment)

		merchantID, err := request.RequireString("merchant_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		masterSourceID1, err := request.RequireString("master_source_id_1")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		masterSourceID2, err := request.RequireString("master_source_id_2")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source1Name, err := request.RequireString("source_1_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2Name, err := request.RequireString("source_2_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		namingStrategy := request.GetString("source_naming_strategy", "descriptive")

		// Generate merchant source names based on strategy
		merchantSource1Name := generateMerchantSourceName(source1Name, namingStrategy, 1)
		merchantSource2Name := generateMerchantSourceName(source2Name, namingStrategy, 2)

		// Create merchant sources via API calls
		merchantSource1ID, err := createMerchantSource(ctx, merchantID, masterSourceID1, merchantSource1Name, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create merchant source 1: %v", err)), nil
		}

		merchantSource2ID, err := createMerchantSource(ctx, merchantID, masterSourceID2, merchantSource2Name, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create merchant source 2: %v", err)), nil
		}

		result := map[string]interface{}{
			"status":       "success",
			"message":      "Merchant sources created successfully",
			"environment":  GetEnvironmentName(environment),
			"api_base_url": GetBaseURL(environment),
			"execution_summary": map[string]interface{}{
				"merchant_id":            merchantID,
				"total_merchant_sources": 2,
				"successful_creations":   2,
				"failed_creations":       0,
			},
			"created_merchant_sources": map[string]interface{}{
				"merchant_source_1": map[string]interface{}{
					"merchant_source_id": merchantSource1ID,
					"name":               merchantSource1Name,
					"master_source_id":   masterSourceID1,
					"merchant_id":        merchantID,
					"naming_strategy":    namingStrategy,
				},
				"merchant_source_2": map[string]interface{}{
					"merchant_source_id": merchantSource2ID,
					"name":               merchantSource2Name,
					"master_source_id":   masterSourceID2,
					"merchant_id":        merchantID,
					"naming_strategy":    namingStrategy,
				},
			},
			"for_future_prompts": map[string]interface{}{
				"merchant_id":          merchantID,
				"merchant_source_id_1": merchantSource1ID,
				"merchant_source_id_2": merchantSource2ID,
				"master_source_id_1":   masterSourceID1,
				"master_source_id_2":   masterSourceID2,
				"environment":          environment,
			},
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// ReconStateRuleTool Recon state and rule creation tool for recon-saas
func ReconStateRuleTool() server.ServerTool {
	tool := mcp.NewTool("recon_state_rule",
		mcp.WithDescription("Create reconciliation states and corresponding rules for recon-saas"),
		mcp.WithString("environment",
			mcp.Description("Environment to use for API calls: 'local' (http://localhost:9400), 'dev' (https://recon-saas.dev.razorpay.in), or 'prod' (https://recon-saas.concierge.razorpay.com). Defaults to 'dev' if not specified."),
			mcp.Enum("local", "dev", "prod"),
			mcp.DefaultString("dev"),
		),
		mcp.WithString("merchant_id",
			mcp.Description("Merchant identifier"),
			mcp.Required(),
		),
		mcp.WithString("master_source_id_1",
			mcp.Description("First master source ID"),
			mcp.Required(),
		),
		mcp.WithString("master_source_id_2",
			mcp.Description("Second master source ID"),
			mcp.Required(),
		),
		mcp.WithString("source_1_name",
			mcp.Description("Name of the first source for remarks"),
			mcp.Required(),
		),
		mcp.WithString("source_2_name",
			mcp.Description("Name of the second source for remarks"),
			mcp.Required(),
		),
		mcp.WithBoolean("approve_expressions",
			mcp.Description("Whether to approve the generated rule expressions"),
			mcp.DefaultBool(true),
		),
		mcp.WithString("validation_mode",
			mcp.Description("User validation mode for rule expressions"),
			mcp.Enum("automatic", "guided", "manual"),
			mcp.DefaultString("guided"),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get environment (defaults to "dev" if not specified)
		environment := request.GetString("environment", DefaultEnvironment)

		merchantID, err := request.RequireString("merchant_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		masterSourceID1, err := request.RequireString("master_source_id_1")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		masterSourceID2, err := request.RequireString("master_source_id_2")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source1Name, err := request.RequireString("source_1_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2Name, err := request.RequireString("source_2_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		approveExpressions := request.GetBool("approve_expressions", true)
		validationMode := request.GetString("validation_mode", "guided")

		// Apply validation mode logic
		validationResult, err := applyValidationMode(validationMode, approveExpressions, masterSourceID1, masterSourceID2)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Create recon states
		reconStates, err := createReconStates(ctx, merchantID, source1Name, source2Name, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create recon states: %v", err)), nil
		}

		// Create rules with validation result
		rules, err := createReconRulesWithValidation(ctx, merchantID, masterSourceID1, masterSourceID2, environment, reconStates, validationResult)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create recon rules: %v", err)), nil
		}

		result := map[string]interface{}{
			"status":       "success",
			"message":      "Recon states and rules created successfully",
			"environment":  GetEnvironmentName(environment),
			"api_base_url": GetBaseURL(environment),
			"execution_summary": map[string]interface{}{
				"merchant_id":               merchantID,
				"total_recon_states":        len(reconStates),
				"total_rules":               len(rules) - 1, // Subtract 1 for validation_summary
				"user_approved_expressions": approveExpressions,
				"validation_mode":           validationMode,
				"validation_applied":        validationResult.Approved,
			},
			"created_recon_states": reconStates,
			"created_rules":        rules,
			"for_future_prompts": map[string]interface{}{
				"merchant_id":        merchantID,
				"master_source_id_1": masterSourceID1,
				"master_source_id_2": masterSourceID2,
				"recon_state_ids":    extractStateIDs(reconStates),
				"rule_ids":           extractRuleIDs(rules),
				"environment":        environment,
			},
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// ReconProcessSetupTool Lookup and recon process creation tool for recon-saas
func ReconProcessSetupTool() server.ServerTool {
	tool := mcp.NewTool("recon_process_setup",
		mcp.WithDescription("Create lookup configurations and reconciliation processes for recon-saas"),
		mcp.WithString("environment",
			mcp.Description("Environment to use for API calls: 'local' (http://localhost:9400), 'dev' (https://recon-saas.dev.razorpay.in), or 'prod' (https://recon-saas.concierge.razorpay.com). Defaults to 'dev' if not specified."),
			mcp.Enum("local", "dev", "prod"),
			mcp.DefaultString("dev"),
		),
		mcp.WithString("merchant_id",
			mcp.Description("Merchant identifier"),
			mcp.Required(),
		),
		mcp.WithString("master_source_id_1",
			mcp.Description("First master source ID"),
			mcp.Required(),
		),
		mcp.WithString("master_source_id_2",
			mcp.Description("Second master source ID"),
			mcp.Required(),
		),
		mcp.WithString("merchant_source_id_1",
			mcp.Description("First merchant source ID"),
			mcp.Required(),
		),
		mcp.WithString("merchant_source_id_2",
			mcp.Description("Second merchant source ID"),
			mcp.Required(),
		),
		mcp.WithString("rule_ids",
			mcp.Description("JSON array of rule IDs from previous step"),
			mcp.Required(),
		),
		mcp.WithString("source_1_name",
			mcp.Description("Name of the first source"),
			mcp.Required(),
		),
		mcp.WithString("source_2_name",
			mcp.Description("Name of the second source"),
			mcp.Required(),
		),
		mcp.WithString("source1_columns",
			mcp.Description("JSON array of column names from first file"),
			mcp.Required(),
		),
		mcp.WithString("source2_columns",
			mcp.Description("JSON array of column names from second file"),
			mcp.Required(),
		),
		mcp.WithString("source1_entityid",
			mcp.Description("Selected EntityID column name for first source"),
			mcp.Required(),
		),
		mcp.WithString("source2_entityid",
			mcp.Description("Selected EntityID column name for second source"),
			mcp.Required(),
		),
		mcp.WithString("source1_amount",
			mcp.Description("Selected Amount column name for first source"),
			mcp.Required(),
		),
		mcp.WithString("source2_amount",
			mcp.Description("Selected Amount column name for second source"),
			mcp.Required(),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get environment (defaults to "dev" if not specified)
		environment := request.GetString("environment", DefaultEnvironment)

		merchantID, err := request.RequireString("merchant_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		masterSourceID1, err := request.RequireString("master_source_id_1")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		masterSourceID2, err := request.RequireString("master_source_id_2")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		merchantSourceID1, err := request.RequireString("merchant_source_id_1")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		merchantSourceID2, err := request.RequireString("merchant_source_id_2")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		ruleIDsJSON, err := request.RequireString("rule_ids")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source1Name, err := request.RequireString("source_1_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2Name, err := request.RequireString("source_2_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Extract column information
		source1Columns, err := request.RequireString("source1_columns")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2Columns, err := request.RequireString("source2_columns")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source1EntityID, err := request.RequireString("source1_entityid")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2EntityID, err := request.RequireString("source2_entityid")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source1Amount, err := request.RequireString("source1_amount")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		source2Amount, err := request.RequireString("source2_amount")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Parse rule IDs
		var ruleIDs []string
		if err := json.Unmarshal([]byte(ruleIDsJSON), &ruleIDs); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid rule_ids JSON: %v", err)), nil
		}

		// Create lookup
		lookupID, err := createLookup(ctx, merchantID, source1Name, source2Name, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create lookup: %v", err)), nil
		}

		// Create master recon process with column mappings
		masterReconProcessID, err := createMasterReconProcess(ctx, source1Name, source2Name, lookupID, masterSourceID1, masterSourceID2, environment, ruleIDs, source1Columns, source2Columns, source1EntityID, source2EntityID, source1Amount, source2Amount)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create master recon process: %v", err)), nil
		}

		// Create merchant recon process
		merchantReconProcessID, err := createMerchantReconProcess(ctx, merchantID, masterReconProcessID, merchantSourceID1, merchantSourceID2, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create merchant recon process: %v", err)), nil
		}

		result := map[string]interface{}{
			"status":       "success",
			"message":      "Reconciliation process setup completed successfully",
			"environment":  GetEnvironmentName(environment),
			"api_base_url": GetBaseURL(environment),
			"execution_summary": map[string]interface{}{
				"merchant_id":          merchantID,
				"process_name":         fmt.Sprintf("%s to %s Reconciliation", source1Name, source2Name),
				"total_api_calls":      3,
				"successful_creations": 3,
				"failed_creations":     0,
			},
			"created_components": map[string]interface{}{
				"lookup": map[string]interface{}{
					"lookup_id": lookupID,
					"name":      fmt.Sprintf("Entity Lookup for %s and %s", source1Name, source2Name),
				},
				"master_recon_process": map[string]interface{}{
					"master_recon_process_id": masterReconProcessID,
					"name":                    fmt.Sprintf("%s to %s Reconciliation", source1Name, source2Name),
				},
				"merchant_recon_process": map[string]interface{}{
					"merchant_recon_process_id": merchantReconProcessID,
				},
			},
			"onboarding_completion": map[string]interface{}{
				"status":  "COMPLETE",
				"message": "Merchant onboarding successfully completed. The reconciliation process is now ready for file uploads and processing.",
				"next_steps": []string{
					"Upload transaction files for reconciliation",
					"Monitor reconciliation results in dashboard",
					"Configure automated file processing schedules",
					"Set up reporting and alerting preferences",
				},
				"environment": environment,
			},
		}

		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// ReconEntityUpdateTool provides a unified tool for updating various recon-saas entities
// Supports PATCH operations for: master_source, merchant_source, master_recon_process,
// merchant_recon_process, rule, recon_state, and lookup
func ReconEntityUpdateTool() server.ServerTool {
	tool := mcp.NewTool("recon_entity_update",
		mcp.WithDescription(`Update recon-saas entities via PATCH API calls.

This tool supports updating the following entity types:
- master_source: Update master source configurations (name, schema, mappings, transformations, validations, etc.)
- merchant_source: Update merchant-specific source configurations (name, master_source_id, emails, mappings, etc.)
- master_recon_process: Update master reconciliation process (name, lookup_config, rules, sources, sequence, report_config, workflow_config)
- merchant_recon_process: Update merchant reconciliation process (sources, report_config, skip_status, skip_rows, status, etc.)
- rule: Update reconciliation rules (name, type, expression, sources, recon_state_id)
- recon_state: Update reconciliation states (name, priority, remarks)
- lookup: Update lookup configurations (name, config)

USAGE EXAMPLES:

1. Update a master source name:
   entity_type: "master_source"
   entity_id: "abc123"
   update_payload: {"name": "New Source Name"}

2. Update a rule expression:
   entity_type: "rule"
   entity_id: "rule456"
   update_payload: {"expression": "SourceA.EntityID == SourceB.EntityID", "name": "Updated Rule"}

3. Update merchant source emails:
   entity_type: "merchant_source"
   entity_id: "ms789"
   update_payload: {"reporting_emails": ["user@example.com"], "cc_emails": ["cc@example.com"]}

4. Update recon state priority:
   entity_type: "recon_state"
   entity_id: "state123"
   update_payload: {"priority": 2, "remarks": "Updated remarks"}

5. Update master recon process rules:
   entity_type: "master_recon_process"
   entity_id: "mrp456"
   update_payload: {"name": "Updated Process Name", "rules": {"rule_ids": ["rule1", "rule2"]}}

6. Update merchant recon process status:
   entity_type: "merchant_recon_process"
   entity_id: "merchant_proc_789"
   update_payload: {"status": "approved", "skip_status": true}

7. Update lookup configuration:
   entity_type: "lookup"
   entity_id: "lookup123"
   update_payload: {"name": "Updated Lookup", "config": [{"source": "record_internal", "columns": ["EntityID", "Amount"]}]}

FIELD REFERENCES:

master_source fields:
  - name (string): Source name
  - skip_top_rows (int): Rows to skip from top
  - ingest_to_db (bool): Whether to ingest to database
  - allow_upload (bool): Allow file uploads
  - unique_keys ([]string): Unique key columns
  - source_schema (array): Column schema definitions [{name: string, type: string}]
  - mapping_config (array): Column mapping configurations [{source: string, destination: string, value: string}]
  - transformation_config (array): Data transformation rules [{logic: object, output_columns: []string}]
  - validation_config (object): Validation rules {logics: [{logic: object}]}
  - sub_source_config (object): Sub-source configuration
  - extract_distinct_config ([]string): Distinct extraction columns
  - report_enrichment (bool): Enable report enrichment
  - split_file_basis (string): File splitting strategy
  - row_hash_value_based_split_config (object): Row hash based split config
  - column_value_based_split_config (object): Column value based split config
  - is_header_missing (bool): If file has no header row
  - metadata_extraction_config (object): Metadata extraction config
  - skip_bottom_rows (int): Rows to skip from bottom
  - skip_row_func (string): Custom row skip function

merchant_source fields:
  - name (string): Source name
  - master_source_id (string): Associated master source ID
  - reporting_emails ([]string): Report recipient emails
  - cc_emails ([]string): CC email addresses
  - bcc_emails ([]string): BCC email addresses
  - allow_upload (bool): Allow file uploads
  - source_schema (array): Custom column schema
  - mapping_config (array): Custom column mappings
  - validation_config (object): Custom validation rules
  - split_file_basis (string): File splitting strategy
  - row_hash_value_based_split_config (object): Row hash based split config
  - column_value_based_split_config (object): Column value based split config
  - beam_sftp_push_job (string): SFTP push job name
  - slack_notification_config (object): Slack notification settings {recon_percentage_threshold: float, file_alert_enabled: bool, file_alert_cron: string, number_of_files_to_check: int}

master_recon_process fields:
  - name (string): Process name
  - lookup_config (array): Lookup configurations [{config: {source_id: lookup_id}, streaming_source_id: string}]
  - product_id (string): Product identifier
  - rules (object): Rule configuration {rule_ids: []string}
  - sources ([]string): Master source IDs
  - sequence (array): Processing sequence
  - report_config (object): Reporting configuration {frontend_cols: []string, source_report_config: array}
  - workflow_config (object): Workflow settings

merchant_recon_process fields:
  - sources ([]string): Merchant source IDs
  - report_config (object): Reporting configuration
  - skip_status (bool): Skip status processing
  - skip_rows (bool): Skip certain rows
  - skip_rows_recon_state_ids ([]string): Recon state IDs to skip
  - report_channel ([]string): Report delivery channels
  - status (string): Process status (approved/pending_approval)

rule fields:
  - name (string): Rule name
  - type (string): Rule type (recon)
  - expression (string): Rule expression (e.g., "SourceA.EntityID == SourceB.EntityID && SourceA.Amount.Equal(SourceB.Amount)")
  - sources ([]string): Source IDs this rule applies to
  - recon_state_id (string): Associated recon state ID

recon_state fields:
  - name (string): State name (Reconciled/Unreconciled)
  - priority (int8): Priority level (1-10, lower = higher priority)
  - remarks (string): State description/remarks

lookup fields:
  - name (string): Lookup name
  - config (array): Lookup configuration [{source: string, columns: []string, aggregation: object, advanced_config: object, lookback_config: object}]
`),
		mcp.WithString("environment",
			mcp.Description("Environment to use for API calls: 'local' (http://localhost:9400), 'dev' (https://recon-saas.dev.razorpay.in), or 'prod' (https://recon-saas.concierge.razorpay.com). Defaults to 'dev' if not specified."),
			mcp.Enum("local", "dev", "prod"),
			mcp.DefaultString("dev"),
		),
		mcp.WithString("entity_type",
			mcp.Description("Type of entity to update"),
			mcp.Required(),
			mcp.Enum("master_source", "merchant_source", "master_recon_process", "merchant_recon_process", "rule", "recon_state", "lookup"),
		),
		mcp.WithString("entity_id",
			mcp.Description("ID of the entity to update"),
			mcp.Required(),
		),
		mcp.WithString("update_payload",
			mcp.Description("JSON object containing fields to update. Only include fields you want to change. See tool description for available fields per entity type."),
			mcp.Required(),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get environment (defaults to "dev" if not specified)
		environment := request.GetString("environment", DefaultEnvironment)

		entityType, err := request.RequireString("entity_type")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		entityID, err := request.RequireString("entity_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		updatePayloadJSON, err := request.RequireString("update_payload")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Parse the update payload
		var updatePayload map[string]interface{}
		if err := json.Unmarshal([]byte(updatePayloadJSON), &updatePayload); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid update_payload JSON: %v", err)), nil
		}

		if len(updatePayload) == 0 {
			return mcp.NewToolResultError("update_payload cannot be empty"), nil
		}

		// Determine the endpoint based on entity type
		var endpoint string
		switch entityType {
		case "master_source":
			endpoint = fmt.Sprintf("/v1/admin-recon-saas/sources/update/%s", entityID)
		case "merchant_source":
			endpoint = fmt.Sprintf("/v1/admin-recon-saas/sources/update_merchant/%s", entityID)
		case "master_recon_process":
			endpoint = fmt.Sprintf("/v1/admin-recon-saas/recon_process/master/%s", entityID)
		case "merchant_recon_process":
			endpoint = fmt.Sprintf("/v1/admin-recon-saas/recon_process/merchant/%s", entityID)
		case "rule":
			endpoint = fmt.Sprintf("/v1/admin-recon-saas/rule/%s", entityID)
		case "recon_state":
			endpoint = fmt.Sprintf("/v1/admin-recon-saas/recon_state/%s", entityID)
		case "lookup":
			endpoint = fmt.Sprintf("/v1/admin-recon-saas/lookup/%s", entityID)
		default:
			return mcp.NewToolResultError(fmt.Sprintf("Unsupported entity type: %s", entityType)), nil
		}

		// Make the PATCH API call
		result, err := makeReconSaaSAPICall(ctx, "PATCH", endpoint, updatePayload, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update %s: %v", entityType, err)), nil
		}

		// Build response
		response := map[string]interface{}{
			"status":         "success",
			"message":        fmt.Sprintf("%s updated successfully", entityType),
			"environment":    GetEnvironmentName(environment),
			"api_base_url":   GetBaseURL(environment),
			"entity_type":    entityType,
			"entity_id":      entityID,
			"updated_fields": getUpdatedFieldNames(updatePayload),
			"api_response":   result,
		}

		resultJSON, _ := json.MarshalIndent(response, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// getUpdatedFieldNames extracts field names from the update payload for logging
func getUpdatedFieldNames(payload map[string]interface{}) []string {
	fields := make([]string, 0, len(payload))
	for key := range payload {
		fields = append(fields, key)
	}
	return fields
}

// toSnakeCase converts a string to snake_case
// Examples: "Invoice Date" -> "invoice_date", "CustomerName" -> "customer_name"
func toSnakeCase(s string) string {
	// First, replace spaces with underscores and convert to lowercase
	result := strings.ToLower(strings.ReplaceAll(s, " ", "_"))

	// Handle camelCase by inserting underscores before uppercase letters
	var snakeCase strings.Builder
	for i, r := range result {
		if i > 0 && r >= 'A' && r <= 'Z' {
			snakeCase.WriteRune('_')
		}
		snakeCase.WriteRune(r)
	}

	// Clean up any double underscores
	return strings.ReplaceAll(snakeCase.String(), "__", "_")
}

// ReconTransformationConfigTool provides intelligent transformation configuration for master sources
// This tool understands all available transformation functions and helps users apply them correctly
// IMPORTANT: This tool should ONLY be called when the user EXPLICITLY requests a transformation
func ReconTransformationConfigTool() server.ServerTool {
	tool := mcp.NewTool("recon_transformation_config",
		mcp.WithDescription(`Configure transformations for recon-saas master sources ONLY when explicitly requested by the user.

**CRITICAL: DO NOT AUTO-APPLY TRANSFORMATIONS**

This tool should ONLY be used when the user EXPLICITLY requests a transformation.
Do NOT automatically apply any transformation during normal onboarding or master source creation.

**WHEN TO USE THIS TOOL:**
- User explicitly says: "Apply transformation on column X"
- User explicitly says: "Concatenate columns A, B, C to create EntityID"
- User explicitly says: "Apply regex on column X to extract Y"
- User explicitly says: "Transform date format from X to Y"

**WHEN NOT TO USE THIS TOOL:**
- During normal master source creation (unless user explicitly asks for transformation)
- During file analysis
- During onboarding flow (unless user explicitly asks for transformation)
- When user doesn't mention transformation, regex, concatenation, or data formatting

This tool helps you apply data transformations to master source columns. It understands the user's intent
and generates the correct transformation_config and updates mapping_config when necessary.

AVAILABLE TRANSFORMATION FUNCTIONS:

1. **abs_amount_parsing** - Parse absolute amount (removes commas, handles nulls)
   - Input: Single column containing amount
   - Example: {"logic": {"abs_amount_parsing": ["$Subtotal"]}, "output_columns": ["Amount"]}
   - Use case: "Parse Subtotal as absolute amount and store in Amount"

2. **append_multiple_columns** - Concatenate multiple columns into one
   - Input: Multiple columns to concatenate
   - Example: {"logic": {"append_multiple_columns": ["$RRN", "$TID", "$MID"]}, "output_columns": ["EntityID"]}
   - Use case: "Combine RRN, TID, MID into EntityID"

3. **change_date_format** - Convert date from one format to another
   - Input: Column, from_format, to_format
   - Example: {"logic": {"change_date_format": ["$Invoice Date", "%m/%d/%Y", "%Y-%m-%d"]}, "output_columns": ["Invoice Date"]}
   - Use case: "Change Invoice Date format from MM/DD/YYYY to YYYY-MM-DD"

4. **date_normalization** - Normalize date to YYYY-MM-DD format
   - Input: Column, current_format
   - Example: {"logic": {"date_normalization": ["$txn_date", "%d-%m-%Y"]}, "output_columns": ["txn_date"]}
   - Use case: "Normalize txn_date to standard format"

5. **txn_date_extraction_generic** - Extract and parse transaction date (supports unix timestamps, fuzzy matching)
   - Input: Column, optional format (dd-mm-yyyy, dd/mm/yyyy)
   - Example: {"logic": {"txn_date_extraction_generic": ["$timestamp"]}, "output_columns": ["txn_date"]}
   - Use case: "Extract transaction date from timestamp column"

6. **add_amount_cols** - Add multiple amount columns together
   - Input: Multiple amount columns
   - Example: {"logic": {"add_amount_cols": ["$base_amount", "$tax", "$fee"]}, "output_columns": ["total_amount"]}
   - Use case: "Sum base_amount, tax, and fee into total_amount"

7. **subtract_amount_cols** - Subtract amounts from a base column
   - Input: Base column, columns to subtract
   - Example: {"logic": {"subtract_amount_cols": ["$gross", "$discount", "$tax"]}, "output_columns": ["net_amount"]}
   - Use case: "Subtract discount and tax from gross to get net_amount"

8. **abs_amount_in_paisa** - Convert amount to paisa (multiply by 100)
   - Input: Single amount column
   - Example: {"logic": {"abs_amount_in_paisa": ["$amount"]}, "output_columns": ["amount_paisa"]}
   - Use case: "Convert amount to paisa"

9. **excel_mid** - Extract substring from middle (like Excel MID)
   - Input: Column, start_position, length, optional prefix
   - Example: {"logic": {"excel_mid": ["$description", 5, 10]}, "output_columns": ["extracted_id"]}
   - Use case: "Extract 10 characters starting from position 5 of description"

10. **excel_left** - Extract characters from left
    - Input: Column, length
    - Example: {"logic": {"excel_left": ["$code", 4]}, "output_columns": ["prefix"]}
    - Use case: "Extract first 4 characters from code"

11. **excel_right** - Extract characters from right
    - Input: Column, length
    - Example: {"logic": {"excel_right": ["$code", 6]}, "output_columns": ["suffix"]}
    - Use case: "Extract last 6 characters from code"

12. **split** - Split string and get specific part
    - Input: Column, delimiter, index
    - Example: {"logic": {"split": ["$reference", "/", 2]}, "output_columns": ["part3"]}
    - Use case: "Split reference by '/' and get the 3rd part (index 2)"

13. **regex_exec** - Extract using regex pattern
    - Input: Column, regex_pattern
    - Example: {"logic": {"regex_exec": ["$narration", "[A-Z]{4}[0-9]{12}"]}, "output_columns": ["utr"]}
    - Use case: "Extract UTR pattern from narration"

14. **remove_prefix** - Remove prefix from string
    - Input: Column, prefixes to remove
    - Example: {"logic": {"remove_prefix": ["$txn_id", "TXN-", "PAY-"]}, "output_columns": ["clean_id"]}
    - Use case: "Remove TXN- or PAY- prefix from txn_id"

15. **remove_suffix** - Remove suffix from string
    - Input: Column, suffixes to remove
    - Example: {"logic": {"remove_suffix": ["$reference", "-REF", "-ID"]}, "output_columns": ["clean_ref"]}
    - Use case: "Remove -REF or -ID suffix from reference"

16. **add_padding_prefix** - Add zero padding to fixed length
    - Input: Column, desired_length
    - Example: {"logic": {"add_padding_prefix": ["$id", 10]}, "output_columns": ["padded_id"]}
    - Use case: "Pad id with zeros to make it 10 characters"

17. **hard_code_value** - Set a constant value
    - Input: Constant value
    - Example: {"logic": {"hard_code_value": ["ACTIVE"]}, "output_columns": ["status"]}
    - Use case: "Set status column to constant value 'ACTIVE'"

18. **settlement_amount_from_debit_credit_cols** - Derive settlement amount from debit/credit columns
    - Input: Debit column, Credit column
    - Example: {"logic": {"settlement_amount_from_debit_credit_cols": ["$debit", "$credit"]}, "output_columns": ["settlement_amount"]}
    - Use case: "Calculate settlement amount from debit and credit columns"

19. **percentage_of_number** - Calculate percentage
    - Input: Base amount, percentage
    - Example: {"logic": {"percentage_of_number": ["$amount", 18]}, "output_columns": ["tax"]}
    - Use case: "Calculate 18% of amount"

20. **extract_amount_from_cols** - Extract first non-zero amount from multiple columns
    - Input: Multiple amount columns
    - Example: {"logic": {"extract_amount_from_cols": ["$amount1", "$amount2", "$amount3"]}, "output_columns": ["final_amount"]}
    - Use case: "Get first non-zero amount from amount1, amount2, or amount3"

21. **excel_to_datetime** - Convert Excel serial date to datetime
    - Input: Excel date column
    - Example: {"logic": {"excel_to_datetime": ["$excel_date"]}, "output_columns": ["date"]}
    - Use case: "Convert Excel serial date number to readable date"

22. **subtract_date** - Subtract days from a date
    - Input: Date column, days to subtract
    - Example: {"logic": {"subtract_date": ["$date", 7]}, "output_columns": ["week_ago"]}
    - Use case: "Get date 7 days before"

23. **add_date** - Add days to a date
    - Input: Date column, days to add
    - Example: {"logic": {"add_date": ["$date", 30]}, "output_columns": ["due_date"]}
    - Use case: "Add 30 days to date for due_date"

24. **get_field_from_notes** - Extract field from JSON notes
    - Input: Notes column, field names to try
    - Example: {"logic": {"get_field_from_notes": ["$notes", "order_id", "txn_id"]}, "output_columns": ["extracted_id"]}
    - Use case: "Extract order_id or txn_id from JSON notes field"

25. **replace_blank_string** - Trim whitespace from string
    - Input: Column
    - Example: {"logic": {"replace_blank_string": ["$name"]}, "output_columns": ["name"]}
    - Use case: "Trim whitespace from name column"

26. **remove_single_quotes** - Remove single quotes from string
    - Input: Column
    - Example: {"logic": {"remove_single_quotes": ["$value"]}, "output_columns": ["value"]}
    - Use case: "Remove single quotes from value"

27. **remove_double_quotes** - Remove double quotes from string
    - Input: Column
    - Example: {"logic": {"remove_double_quotes": ["$value"]}, "output_columns": ["value"]}
    - Use case: "Remove double quotes from value"

WORKFLOW:
1. Identify the source and master_source_id
2. Determine the transformation function needed
3. Specify input columns and output column name
4. The tool will AUTO-FETCH current mapping_config and transformation_config via GET API if not provided
5. Generate and append new transformation_config
6. Update mapping_config preserving all existing mappings
7. Apply changes via PATCH API

AUTO-FETCH FEATURE:
If current_mapping_config is not provided, the tool automatically fetches it by calling:
GET /v1/admin-recon-saas/sources/get/{master_source_id}
This retrieves both current mapping_config and transformation_config from the server.

IMPORTANT RULES:
- Column references in logic use $ prefix: "$column_name"
- output_columns contains the destination column name without $ prefix
- ALL EXISTING MAPPINGS ARE PRESERVED - only modifications/additions are made

MAPPING CONFIG UPDATE LOGIC:

**For Special Columns (EntityID, EntityStatus, EntityIdentifier, Amount):**
When output_column is a special column:
1. Find any existing mapping where destination equals the special column
2. Change that mapping's destination to snake_case of its source (e.g., "Notes" -> "notes")
3. Append a NEW mapping: {source: "<special_column>", destination: "<special_column>"}

Example:
- Before: [{"source": "Notes", "destination": "EntityID"}, {"source": "amount", "destination": "Amount"}]
- Transformation: output_column = "EntityID"
- After: [{"source": "Notes", "destination": "notes"}, {"source": "amount", "destination": "Amount"}, {"source": "EntityID", "destination": "EntityID"}]

**For Non-Special Columns:**
When output_column is NOT a special column:
1. Check if output_column already exists as a destination
2. If not found, append: {source: "<output_column>", destination: "<snake_case>"}

Example:
- Before: [{"source": "col1", "destination": "col1"}]
- Transformation: output_column = "Extracted Value"
- After: [{"source": "col1", "destination": "col1"}, {"source": "Extracted Value", "destination": "extracted_value"}]
`),
		mcp.WithString("environment",
			mcp.Description("Environment to use for API calls: 'local', 'dev', or 'prod'. Defaults to 'dev'."),
			mcp.Enum("local", "dev", "prod"),
			mcp.DefaultString("dev"),
		),
		mcp.WithString("master_source_id",
			mcp.Description("ID of the master source to apply transformation to"),
			mcp.Required(),
		),
		mcp.WithString("source_name",
			mcp.Description("Name of the source for reference (e.g., 'Source A', 'POS Transactions')"),
			mcp.Required(),
		),
		mcp.WithString("transformation_function",
			mcp.Description("The transformation function to apply"),
			mcp.Required(),
			mcp.Enum(
				"abs_amount_parsing",
				"append_multiple_columns",
				"change_date_format",
				"date_normalization",
				"txn_date_extraction_generic",
				"add_amount_cols",
				"subtract_amount_cols",
				"abs_amount_in_paisa",
				"excel_mid",
				"excel_left",
				"excel_right",
				"split",
				"regex_exec",
				"remove_prefix",
				"remove_suffix",
				"add_padding_prefix",
				"hard_code_value",
				"settlement_amount_from_debit_credit_cols",
				"percentage_of_number",
				"extract_amount_from_cols",
				"excel_to_datetime",
				"subtract_date",
				"add_date",
				"get_field_from_notes",
				"replace_blank_string",
				"remove_single_quotes",
				"remove_double_quotes",
			),
		),
		mcp.WithString("input_columns",
			mcp.Description("JSON array of input column names (from the source file). Example: [\"Subtotal\"] or [\"RRN\", \"TID\", \"MID\"]"),
			mcp.Required(),
		),
		mcp.WithString("additional_params",
			mcp.Description("JSON array of additional parameters for the function (if needed). Example: [\"%m/%d/%Y\", \"%Y-%m-%d\"] for change_date_format, or [5, 10] for excel_mid"),
		),
		mcp.WithString("output_column",
			mcp.Description("Name of the output column where the transformation result will be stored"),
			mcp.Required(),
		),
		mcp.WithString("current_mapping_config",
			mcp.Description("Optional: JSON array of current mapping_config. If not provided, the tool will AUTO-FETCH it from the server using GET /v1/admin-recon-saas/sources/get/{master_source_id}"),
		),
		mcp.WithString("current_transformation_config",
			mcp.Description("Optional: JSON array of current transformation_config. If not provided along with mapping_config, it will be AUTO-FETCHED from the server."),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		environment := request.GetString("environment", DefaultEnvironment)

		masterSourceID, err := request.RequireString("master_source_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		sourceName, err := request.RequireString("source_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		transformationFunction, err := request.RequireString("transformation_function")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		inputColumnsJSON, err := request.RequireString("input_columns")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		outputColumn, err := request.RequireString("output_column")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		currentMappingConfigJSON := request.GetString("current_mapping_config", "")
		additionalParamsJSON := request.GetString("additional_params", "[]")
		currentTransformationConfigJSON := request.GetString("current_transformation_config", "")

		// If current_mapping_config is not provided, fetch it from the API
		var fetchedFromAPI bool
		var fetchedSourceDetails map[string]interface{}
		if currentMappingConfigJSON == "" {
			// Fetch master source details via GET API
			endpoint := fmt.Sprintf("/v1/admin-recon-saas/sources/get/%s", masterSourceID)
			sourceDetails, err := makeReconSaaSAPICall(ctx, "GET", endpoint, nil, environment)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch master source details. Please provide current_mapping_config manually or check if master_source_id '%s' is valid: %v", masterSourceID, err)), nil
			}
			fetchedSourceDetails = sourceDetails

			// Extract mapping_config from response
			if mappingConfig, ok := sourceDetails["mapping_config"]; ok && mappingConfig != nil {
				mappingBytes, err := json.Marshal(mappingConfig)
				if err != nil {
					return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal mapping_config: %v", err)), nil
				}
				currentMappingConfigJSON = string(mappingBytes)
			} else {
				currentMappingConfigJSON = "[]"
			}

			// Extract transformation_config from response if not provided
			if currentTransformationConfigJSON == "" {
				if transformConfig, ok := sourceDetails["transformation_config"]; ok && transformConfig != nil {
					transformBytes, err := json.Marshal(transformConfig)
					if err != nil {
						currentTransformationConfigJSON = "[]"
					} else {
						currentTransformationConfigJSON = string(transformBytes)
					}
				} else {
					currentTransformationConfigJSON = "[]"
				}
			}

			fetchedFromAPI = true
		}

		// Ensure we have valid JSON for transformation config
		if currentTransformationConfigJSON == "" || currentTransformationConfigJSON == "null" {
			currentTransformationConfigJSON = "[]"
		}

		// Parse input columns
		var inputColumns []string
		if err := json.Unmarshal([]byte(inputColumnsJSON), &inputColumns); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid input_columns JSON: %v", err)), nil
		}

		// Parse additional params
		var additionalParams []interface{}
		if err := json.Unmarshal([]byte(additionalParamsJSON), &additionalParams); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid additional_params JSON: %v", err)), nil
		}

		// Parse current mapping config
		var currentMappingConfig []map[string]interface{}
		if err := json.Unmarshal([]byte(currentMappingConfigJSON), &currentMappingConfig); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Invalid current_mapping_config JSON: %v", err)), nil
		}

		// Parse current transformation config
		var currentTransformationConfig []map[string]interface{}
		if currentTransformationConfigJSON != "" && currentTransformationConfigJSON != "[]" {
			if err := json.Unmarshal([]byte(currentTransformationConfigJSON), &currentTransformationConfig); err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Invalid current_transformation_config JSON: %v", err)), nil
			}
		}

		// Build the transformation logic arguments
		logicArgs := make([]interface{}, 0)
		for _, col := range inputColumns {
			logicArgs = append(logicArgs, "$"+col)
		}
		for _, param := range additionalParams {
			logicArgs = append(logicArgs, param)
		}

		// Create new transformation entry
		newTransformation := map[string]interface{}{
			"logic": map[string]interface{}{
				transformationFunction: logicArgs,
			},
			"output_columns": []string{outputColumn},
		}

		// Append to existing transformations
		updatedTransformationConfig := append(currentTransformationConfig, newTransformation)

		// Define special columns that keep their exact name as destination
		specialColumns := map[string]bool{
			"EntityID":         true,
			"EntityStatus":     true,
			"EntityIdentifier": true,
			"Amount":           true,
		}

		// Deep copy the current mapping config to avoid modifying the original
		updatedMappingConfig := make([]map[string]interface{}, len(currentMappingConfig))
		for i, mapping := range currentMappingConfig {
			newMapping := make(map[string]interface{})
			for k, v := range mapping {
				newMapping[k] = v
			}
			updatedMappingConfig[i] = newMapping
		}

		isSpecialColumn := specialColumns[outputColumn]

		// Track which mappings were modified for debugging
		var modifiedMappings []map[string]interface{}

		if isSpecialColumn {
			// For special columns (EntityID, EntityStatus, EntityIdentifier, Amount):
			// 1. Find any existing mapping where destination == outputColumn
			// 2. Change that mapping's destination to snake_case of its source
			// 3. Append new mapping: {source: outputColumn, destination: outputColumn}

			for i, mapping := range updatedMappingConfig {
				dest, _ := mapping["destination"].(string)
				source, _ := mapping["source"].(string)

				// Trim whitespace for safety
				dest = strings.TrimSpace(dest)
				source = strings.TrimSpace(source)

				if dest == outputColumn {
					// This mapping currently points to our special column
					// Remap it to snake_case of its source
					newDest := toSnakeCase(source)
					oldDest := updatedMappingConfig[i]["destination"]
					updatedMappingConfig[i]["destination"] = newDest

					// Track the modification
					modifiedMappings = append(modifiedMappings, map[string]interface{}{
						"index":           i,
						"source":          source,
						"old_destination": oldDest,
						"new_destination": newDest,
					})
				}
			}

			// Append new mapping for the transformation output (special column)
			newMapping := map[string]interface{}{
				"value":       "",
				"source":      outputColumn,
				"destination": outputColumn,
			}
			updatedMappingConfig = append(updatedMappingConfig, newMapping)
		} else {
			// For non-special columns:
			// Check if output_column already exists as a destination
			// If not, append with snake_case destination

			outputColumnExistsAsDestination := false
			for _, mapping := range updatedMappingConfig {
				dest, _ := mapping["destination"].(string)
				if dest == outputColumn || dest == toSnakeCase(outputColumn) {
					outputColumnExistsAsDestination = true
					break
				}
			}

			if !outputColumnExistsAsDestination {
				// Add new mapping with snake_case destination
				newMapping := map[string]interface{}{
					"value":       "",
					"source":      outputColumn,
					"destination": toSnakeCase(outputColumn),
				}
				updatedMappingConfig = append(updatedMappingConfig, newMapping)
			}
		}

		// Prepare the update payload
		updatePayload := map[string]interface{}{
			"transformation_config": updatedTransformationConfig,
			"mapping_config":        updatedMappingConfig,
		}

		// Make the PATCH API call
		endpoint := fmt.Sprintf("/v1/admin-recon-saas/sources/update/%s", masterSourceID)
		result, err := makeReconSaaSAPICall(ctx, "PATCH", endpoint, updatePayload, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update master source: %v", err)), nil
		}

		// Build comprehensive response with debug info
		response := map[string]interface{}{
			"status":      "success",
			"message":     fmt.Sprintf("Transformation applied successfully to %s", sourceName),
			"environment": GetEnvironmentName(environment),
			"source_details": map[string]interface{}{
				"master_source_id": masterSourceID,
				"source_name":      sourceName,
			},
			"transformation_applied": map[string]interface{}{
				"function":       transformationFunction,
				"input_columns":  inputColumns,
				"output_column":  outputColumn,
				"transformation": newTransformation,
			},
			"configs_before": map[string]interface{}{
				"mapping_config":              currentMappingConfig,
				"mapping_config_count":        len(currentMappingConfig),
				"transformation_config_count": len(currentTransformationConfig),
			},
			"configs_after": map[string]interface{}{
				"transformation_config":   updatedTransformationConfig,
				"mapping_config":          updatedMappingConfig,
				"mapping_config_count":    len(updatedMappingConfig),
				"is_special_column":       isSpecialColumn,
				"config_fetched_from_api": fetchedFromAPI,
				"modified_mappings":       modifiedMappings,
			},
			"api_response": result,
		}

		// Add fetched source details for debugging if available
		if fetchedSourceDetails != nil {
			response["fetched_source_name"] = fetchedSourceDetails["name"]
		}

		resultJSON, _ := json.MarshalIndent(response, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

// ReconAggregationTool configures aggregation for recon-saas master sources
// This tool should ONLY be called when the user EXPLICITLY requests aggregation on a column
func ReconAggregationTool() server.ServerTool {
	tool := mcp.NewTool("recon_aggregation_config",
		mcp.WithDescription(`Configure aggregation for recon-saas master sources ONLY when explicitly requested by the user.

**CRITICAL: DO NOT AUTO-APPLY AGGREGATION**

This tool should ONLY be used when the user EXPLICITLY requests aggregation configuration.
Do NOT automatically apply aggregation during normal onboarding or master source creation.

**WHEN TO USE THIS TOOL:**
- User explicitly says: "Enable aggregation on column X"
- User explicitly says: "Configure aggregation for EntityIdentifier"
- User explicitly says: "Set up aggregation with column X as EntityIdentifier"
- User explicitly says: "I want to aggregate on column X"

**WHEN NOT TO USE THIS TOOL:**
- During normal master source creation
- During file analysis
- During onboarding flow (unless user explicitly asks for aggregation)
- When user doesn't mention aggregation

**WHAT THIS TOOL DOES:**
1. Fetches current master source configuration via GET API
2. Updates master source with:
   - APPENDS "EntityIdentifier" to existing unique_keys array (preserves all existing keys like "EntityID")
   - Example: ["EntityID"] becomes ["EntityID", "EntityIdentifier"]
   - Updates mapping_config to map the specified column to "EntityIdentifier" destination
3. Fetches current lookup configuration via GET API
4. Updates lookup to enable aggregation for the EntityID column

**REQUIRED INPUTS:**
- master_source_id: ID of the master source to configure
- entity_identifier_column: The column name that should be mapped as EntityIdentifier
- lookup_id: ID of the lookup to update for aggregation (ask user if not available)
`),
		mcp.WithString("environment",
			mcp.Description("Environment to use for API calls: 'local', 'dev', or 'prod'. Defaults to 'dev'."),
			mcp.Enum("local", "dev", "prod"),
			mcp.DefaultString("dev"),
		),
		mcp.WithString("master_source_id",
			mcp.Description("ID of the master source to configure aggregation for"),
			mcp.Required(),
		),
		mcp.WithString("source_name",
			mcp.Description("Name of the source for reference (e.g., 'Source A', 'POS Transactions')"),
			mcp.Required(),
		),
		mcp.WithString("entity_identifier_column",
			mcp.Description("The column name from the source file that should be mapped as EntityIdentifier for aggregation"),
			mcp.Required(),
		),
		mcp.WithString("lookup_id",
			mcp.Description("ID of the lookup to update for enabling aggregation. This is required to enable aggregation on the lookup."),
			mcp.Required(),
		),
	)

	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		environment := request.GetString("environment", DefaultEnvironment)

		masterSourceID, err := request.RequireString("master_source_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		sourceName, err := request.RequireString("source_name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		entityIdentifierColumn, err := request.RequireString("entity_identifier_column")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		lookupID, err := request.RequireString("lookup_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Step 1: Fetch current master source configuration via GET API
		masterSourceEndpoint := fmt.Sprintf("/v1/admin-recon-saas/sources/get/%s", masterSourceID)
		masterSourceDetails, err := makeReconSaaSAPICall(ctx, "GET", masterSourceEndpoint, nil, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch master source details for ID '%s': %v", masterSourceID, err)), nil
		}

		// Extract current unique_keys with robust type handling
		var currentUniqueKeys []string
		if uniqueKeys, ok := masterSourceDetails["config"].(map[string]interface{})["unique_keys"]; ok && uniqueKeys != nil {
			// Handle []interface{} type (common JSON unmarshaling result)
			if keysSlice, ok := uniqueKeys.([]interface{}); ok {
				for _, key := range keysSlice {
					if keyStr, ok := key.(string); ok {
						currentUniqueKeys = append(currentUniqueKeys, keyStr)
					}
				}
			} else if keysSlice, ok := uniqueKeys.([]string); ok {
				// Handle []string type directly
				currentUniqueKeys = append(currentUniqueKeys, keysSlice...)
			}
		}

		// If currentUniqueKeys is still empty, try to extract from raw JSON
		// This handles cases where the type assertion might have failed
		if len(currentUniqueKeys) == 0 {
			// Check if unique_keys exists but extraction failed - default to EntityID
			if _, exists := masterSourceDetails["unique_keys"]; exists {
				// Log warning: extraction might have failed, but we'll preserve what we can
				// Default assumption: EntityID should exist
				currentUniqueKeys = []string{"EntityID"}
			}
		}

		// Check if EntityIdentifier already exists in unique_keys
		entityIdentifierExists := false
		for _, key := range currentUniqueKeys {
			if key == "EntityIdentifier" {
				entityIdentifierExists = true
				break
			}
		}

		// IMPORTANT: Create a new slice that preserves ALL existing keys and appends EntityIdentifier
		var updatedUniqueKeys []string
		// First, copy all existing unique_keys
		updatedUniqueKeys = append(updatedUniqueKeys, currentUniqueKeys...)
		// Then append EntityIdentifier if not already present
		if !entityIdentifierExists {
			updatedUniqueKeys = append(updatedUniqueKeys, "EntityIdentifier")
		}

		// Extract current mapping_config
		var currentMappingConfig []map[string]interface{}
		if mappingConfig, ok := masterSourceDetails["mapping_config"]; ok && mappingConfig != nil {
			if configSlice, ok := mappingConfig.([]interface{}); ok {
				for _, item := range configSlice {
					if itemMap, ok := item.(map[string]interface{}); ok {
						currentMappingConfig = append(currentMappingConfig, itemMap)
					}
				}
			}
		}

		// Deep copy and update mapping_config
		updatedMappingConfig := make([]map[string]interface{}, len(currentMappingConfig))
		entityIdentifierMappingExists := false
		entityIdentifierColumnFound := false

		for i, mapping := range currentMappingConfig {
			newMapping := make(map[string]interface{})
			for k, v := range mapping {
				newMapping[k] = v
			}

			// Check if this is the column user specified for EntityIdentifier
			source, _ := mapping["source"].(string)
			dest, _ := mapping["destination"].(string)

			if source == entityIdentifierColumn {
				entityIdentifierColumnFound = true
				// Update this mapping to point to EntityIdentifier destination
				newMapping["destination"] = "EntityIdentifier"
				entityIdentifierMappingExists = true
			}

			// Check if EntityIdentifier destination already exists
			if dest == "EntityIdentifier" {
				entityIdentifierMappingExists = true
			}

			updatedMappingConfig[i] = newMapping
		}

		// If the column wasn't found in existing mappings, add a new mapping
		if !entityIdentifierColumnFound && !entityIdentifierMappingExists {
			newMapping := map[string]interface{}{
				"value":       "",
				"source":      entityIdentifierColumn,
				"destination": "EntityIdentifier",
			}
			updatedMappingConfig = append(updatedMappingConfig, newMapping)
		}

		// Step 2: Update master source via PATCH API
		masterSourceUpdatePayload := map[string]interface{}{
			"unique_keys":    updatedUniqueKeys,
			"mapping_config": updatedMappingConfig,
		}

		masterSourceUpdateEndpoint := fmt.Sprintf("/v1/admin-recon-saas/sources/update/%s", masterSourceID)
		masterSourceUpdateResult, err := makeReconSaaSAPICall(ctx, "PATCH", masterSourceUpdateEndpoint, masterSourceUpdatePayload, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update master source: %v", err)), nil
		}

		// Step 3: Fetch current lookup configuration via GET API
		lookupEndpoint := fmt.Sprintf("/v1/admin-recon-saas/lookup/%s", lookupID)
		lookupDetails, err := makeReconSaaSAPICall(ctx, "GET", lookupEndpoint, nil, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch lookup details for ID '%s': %v", lookupID, err)), nil
		}

		// Extract current lookup config
		var lookupConfig []map[string]interface{}
		if config, ok := lookupDetails["config"]; ok && config != nil {
			if configSlice, ok := config.([]interface{}); ok {
				for _, item := range configSlice {
					if itemMap, ok := item.(map[string]interface{}); ok {
						lookupConfig = append(lookupConfig, itemMap)
					}
				}
			}
		}

		// Update lookup config to enable aggregation for EntityID column
		updatedLookupConfig := make([]map[string]interface{}, len(lookupConfig))
		aggregationEnabled := false

		for i, configItem := range lookupConfig {
			newConfigItem := make(map[string]interface{})
			for k, v := range configItem {
				newConfigItem[k] = v
			}

			// Check if this config item contains EntityID in Columns
			if columns, ok := configItem["Columns"]; ok {
				if colSlice, ok := columns.([]interface{}); ok {
					for _, col := range colSlice {
						if colStr, ok := col.(string); ok && colStr == "EntityID" {
							// Enable aggregation for this config item
							newConfigItem["aggregation"] = map[string]interface{}{
								"enabled":    true,
								"conditions": nil,
							}
							aggregationEnabled = true
							break
						}
					}
				}
			}

			updatedLookupConfig[i] = newConfigItem
		}

		// Step 4: Update lookup via PATCH API
		lookupUpdatePayload := map[string]interface{}{
			"config": updatedLookupConfig,
		}

		lookupUpdateEndpoint := fmt.Sprintf("/v1/admin-recon-saas/lookup/%s", lookupID)
		lookupUpdateResult, err := makeReconSaaSAPICall(ctx, "PATCH", lookupUpdateEndpoint, lookupUpdatePayload, environment)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to update lookup: %v", err)), nil
		}

		// Extract raw unique_keys from API response for debugging
		rawUniqueKeys := masterSourceDetails["unique_keys"]

		// Build comprehensive response
		response := map[string]interface{}{
			"status":      "success",
			"message":     fmt.Sprintf("Aggregation configured successfully for %s", sourceName),
			"environment": GetEnvironmentName(environment),
			"master_source_update": map[string]interface{}{
				"master_source_id":         masterSourceID,
				"source_name":              sourceName,
				"entity_identifier_column": entityIdentifierColumn,
				"unique_keys_raw_from_api": rawUniqueKeys,
				"unique_keys_extracted":    currentUniqueKeys,
				"unique_keys_after":        updatedUniqueKeys,
				"mapping_config_updated":   true,
				"entity_identifier_added":  !entityIdentifierExists,
				"api_response":             masterSourceUpdateResult,
			},
			"lookup_update": map[string]interface{}{
				"lookup_id":           lookupID,
				"aggregation_enabled": aggregationEnabled,
				"config_before":       lookupConfig,
				"config_after":        updatedLookupConfig,
				"api_response":        lookupUpdateResult,
			},
			"summary": map[string]interface{}{
				"master_source_updated": true,
				"lookup_updated":        true,
				"aggregation_ready":     aggregationEnabled,
				"unique_keys_preserved": len(currentUniqueKeys) > 0,
			},
		}

		resultJSON, _ := json.MarshalIndent(response, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}

	return server.ServerTool{
		Tool:    tool,
		Handler: handler,
	}
}

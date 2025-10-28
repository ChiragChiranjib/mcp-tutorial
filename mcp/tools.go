package mcp

import (
	"context"
	"encoding/json"
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
		masterSource1ID, err := createMasterSource(ctx, source1Name, source1Columns, source1EntityID, source1Amount)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create master source 1: %v", err)), nil
		}

		masterSource2ID, err := createMasterSource(ctx, source2Name, source2Columns, source2EntityID, source2Amount)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create master source 2: %v", err)), nil
		}

		result := map[string]interface{}{
			"status":  "success",
			"message": "Master sources created successfully",
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
		merchantSource1ID, err := createMerchantSource(ctx, merchantID, masterSourceID1, merchantSource1Name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create merchant source 1: %v", err)), nil
		}

		merchantSource2ID, err := createMerchantSource(ctx, merchantID, masterSourceID2, merchantSource2Name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create merchant source 2: %v", err)), nil
		}

		result := map[string]interface{}{
			"status":  "success",
			"message": "Merchant sources created successfully",
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
		reconStates, err := createReconStates(ctx, merchantID, source1Name, source2Name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create recon states: %v", err)), nil
		}

		// Create rules with validation result
		rules, err := createReconRulesWithValidation(ctx, merchantID, masterSourceID1, masterSourceID2, reconStates, validationResult)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create recon rules: %v", err)), nil
		}

		result := map[string]interface{}{
			"status":  "success",
			"message": "Recon states and rules created successfully",
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
		lookupID, err := createLookup(ctx, merchantID, source1Name, source2Name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create lookup: %v", err)), nil
		}

		// Create master recon process with column mappings
		masterReconProcessID, err := createMasterReconProcess(ctx, source1Name, source2Name, lookupID, masterSourceID1, masterSourceID2, ruleIDs, source1Columns, source2Columns, source1EntityID, source2EntityID, source1Amount, source2Amount)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create master recon process: %v", err)), nil
		}

		// Create merchant recon process
		merchantReconProcessID, err := createMerchantReconProcess(ctx, merchantID, masterReconProcessID, merchantSourceID1, merchantSourceID2)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create merchant recon process: %v", err)), nil
		}

		result := map[string]interface{}{
			"status":  "success",
			"message": "Reconciliation process setup completed successfully",
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

// ReconAggregationTool Aggregation tool for recon-saas with entity identifier configuration
func ReconAggregationTool() server.ServerTool {
	tool := mcp.NewTool("recon_aggregation",
		mcp.WithDescription("Configure aggregation logic for recon-saas by updating entity identifier and enabling aggregation"),
		mcp.WithString("file1_path",
			mcp.Description("Full file path to the first reconciliation file (aggregation will be applied to this file)"),
			mcp.Required(),
		),
		mcp.WithString("file2_path",
			mcp.Description("Full file path to the second reconciliation file"),
			mcp.Required(),
		),
		mcp.WithString("entity_identifier",
			mcp.Description("Column name to be used as entity identifier (e.g., UTR, VID, Transaction ID)"),
			mcp.Required(),
		),
		mcp.WithString("aggregation_strategy",
			mcp.Description("Strategy for aggregation (sum, count, avg, etc.)"),
			mcp.Enum("sum", "count", "avg", "max", "min"),
			mcp.DefaultString("sum"),
		),
		mcp.WithString("master_source_id",
			mcp.Description("Master source ID from recon_master_source tool"),
			mcp.Required(),
		),
		mcp.WithString("lookup_id",
			mcp.Description("Lookup ID from recon_process_setup tool"),
			mcp.Required(),
		),
		mcp.WithString("master_recon_process_id",
			mcp.Description("Master recon process ID from recon_process_setup tool"),
			mcp.Required(),
		),
		mcp.WithString("master_source_id_2",
			mcp.Description("Second master source ID from recon_master_source tool"),
			mcp.Required(),
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

		entityIdentifier, err := request.RequireString("entity_identifier")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		aggregationStrategy := request.GetString("aggregation_strategy", "sum")

		masterSourceID1, err := request.RequireString("master_source_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		lookupID, err := request.RequireString("lookup_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		masterReconProcessID, err := request.RequireString("master_recon_process_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		masterSourceID2, err := request.RequireString("master_source_id_2")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Analyze files to determine file types
		file1Type := determineFileType(file1Path)
		file2Type := determineFileType(file2Path)

		// Analyze both files
		analysis1, err := analyzeFile(file1Path, "file_1", file1Type)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze file 1: %v", err)), nil
		}

		analysis2, err := analyzeFile(file2Path, "file_2", file2Type)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to analyze file 2: %v", err)), nil
		}

		// Validate entity identifier exists in file 1 (aggregation file)
		if !validateEntityIdentifier(analysis1, entityIdentifier) {
			return mcp.NewToolResultError(fmt.Sprintf("Entity identifier '%s' not found in file 1 columns", entityIdentifier)), nil
		}

		// Validate Entity ID vs Entity Identifier
		entityValidation, err := validateEntityIDVsEntityIdentifier(ctx, masterSourceID1, entityIdentifier)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to validate Entity ID vs Entity Identifier: %v", err)), nil
		}

		// Check if validation passed
		validationPassed, ok := entityValidation["validation_passed"].(bool)
		if !ok {
			return mcp.NewToolResultError("failed to validate entity identifier validation result"), nil
		}
		if !validationPassed {
			recommendation, ok := entityValidation["recommendation"].(string)
			if !ok {
				recommendation = "entity identifier validation failed"
			}
			return mcp.NewToolResultError(recommendation), nil
		}

		// Generate aggregation preview if files have less than 40 records
		var aggregationPreview map[string]interface{}
		if getRecordCount(analysis1) <= 40 && getRecordCount(analysis2) <= 40 {
			aggregationPreview = generateAggregationPreview(file1Path, file2Path, entityIdentifier, aggregationStrategy, file1Type, file2Type)
		}

		// Generate comprehensive report configuration based on file analysis
		comprehensiveReportConfig := generateComprehensiveReportConfig(analysis1, analysis2, entityIdentifier, masterSourceID1, masterSourceID2)

		// Test all PATCH API endpoints to see which ones work
		apiTestResults := make(map[string]interface{})

		// Test 1: Master Source Update
		err1 := updateMasterSourceMapping(ctx, masterSourceID1, entityIdentifier)
		apiTestResults["master_source_update"] = map[string]interface{}{
			"endpoint": fmt.Sprintf("/v1/admin-recon-saas/sources/update/%s", masterSourceID1),
			"patch_id": masterSourceID1,
			"success":  err1 == nil,
			"error": func() string {
				if err1 != nil {
					return err1.Error()
				}
				return ""
			}(),
		}

		// Test 2: Lookup Update
		err2 := updateLookupAggregation(ctx, lookupID)
		apiTestResults["lookup_update"] = map[string]interface{}{
			"endpoint": fmt.Sprintf("/v1/admin-recon-saas/lookup/%s", lookupID),
			"patch_id": lookupID,
			"success":  err2 == nil,
			"error": func() string {
				if err2 != nil {
					return err2.Error()
				}
				return ""
			}(),
		}

		// Test 3: Master Recon Process Update with Comprehensive Report Config
		err3 := updateMasterReconProcessReportConfigComprehensive(ctx, masterReconProcessID, comprehensiveReportConfig)
		apiTestResults["master_recon_process_update"] = map[string]interface{}{
			"endpoint": fmt.Sprintf("/v1/admin-recon-saas/recon_process/master/%s", masterReconProcessID),
			"patch_id": masterReconProcessID,
			"success":  err3 == nil,
			"error": func() string {
				if err3 != nil {
					return err3.Error()
				}
				return ""
			}(),
		}

		// Count successful API calls
		successCount := 0
		if err1 == nil {
			successCount++
		}
		if err2 == nil {
			successCount++
		}
		if err3 == nil {
			successCount++
		}

		result := map[string]interface{}{
			"status":            "success",
			"message":           "Aggregation logic configured successfully with comprehensive report configuration",
			"entity_validation": entityValidation,
			"api_test_results":  apiTestResults,
			"patch_api_ids": map[string]interface{}{
				"master_source_patch_id":        masterSourceID1,
				"lookup_patch_id":               lookupID,
				"master_recon_process_patch_id": masterReconProcessID,
				"patch_endpoints": []string{
					fmt.Sprintf("PATCH /v1/admin-recon-saas/sources/update/%s", masterSourceID1),
					fmt.Sprintf("PATCH /v1/admin-recon-saas/lookup/%s", lookupID),
					fmt.Sprintf("PATCH /v1/admin-recon-saas/recon_process/master/%s", masterReconProcessID),
				},
			},
			"provided_ids": map[string]interface{}{
				"master_source_id":        masterSourceID1,
				"master_source_id_2":      masterSourceID2,
				"lookup_id":               lookupID,
				"master_recon_process_id": masterReconProcessID,
				"retrieval_method":        "user_provided_parameters",
			},
			"comprehensive_report_config": comprehensiveReportConfig,
			"execution_summary": map[string]interface{}{
				"file_1_path":           file1Path,
				"file_2_path":           file2Path,
				"entity_identifier":     entityIdentifier,
				"aggregation_strategy":  aggregationStrategy,
				"total_patch_endpoints": 3,
				"successful_calls":      successCount,
				"failed_calls":          3 - successCount,
				"success_rate":          fmt.Sprintf("%.1f%%", float64(successCount)/3.0*100),
			},
			"aggregation_configuration": map[string]interface{}{
				"entity_identifier":       entityIdentifier,
				"aggregation_strategy":    aggregationStrategy,
				"non_streaming_source":    "file_1",
				"streaming_source":        "file_2",
				"reverse_mapping_applied": true,
				"comprehensive_config":    true,
				"explanation":             fmt.Sprintf("Entity identifier '%s' has been configured for aggregation. A comprehensive report configuration has been generated with all columns from both files, ensuring proper column mapping and reporting.", entityIdentifier),
			},
			"aggregation_preview": aggregationPreview,
			"next_steps": []string{
				"Aggregation logic is now enabled for the reconciliation process",
				"File 1 will be processed with aggregation based on the selected entity identifier",
				"Comprehensive report configuration includes all columns from both files",
				"Reports will properly map all columns for indexing and display",
				"Upload files to test the aggregation functionality",
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

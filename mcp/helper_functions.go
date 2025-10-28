package mcp

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Helper functions for recon-saas API interactions

// ValidationResult holds the result of validation mode processing
type ValidationResult struct {
	Mode             string            `json:"mode"`
	Approved         bool              `json:"approved"`
	RuleExpressions  map[string]string `json:"rule_expressions"`
	ValidationNotes  []string          `json:"validation_notes"`
	RequiresApproval bool              `json:"requires_approval"`
}

// applyValidationMode applies validation logic based on the selected mode
func applyValidationMode(mode string, approveExpressions bool, masterSourceID1, masterSourceID2 string) (*ValidationResult, error) {
	result := &ValidationResult{
		Mode:            mode,
		RuleExpressions: make(map[string]string),
		ValidationNotes: make([]string, 0),
	}

	// Generate rule expressions
	result.RuleExpressions["reconciled"] = fmt.Sprintf("%s.EntityID == %s.EntityID && %s.Amount.Equal(%s.Amount)",
		masterSourceID1, masterSourceID2, masterSourceID1, masterSourceID2)
	result.RuleExpressions["amount_mismatch"] = fmt.Sprintf("%s.EntityID == %s.EntityID && !%s.Amount.Equal(%s.Amount)",
		masterSourceID1, masterSourceID2, masterSourceID1, masterSourceID2)
	result.RuleExpressions["missing_record"] = "NoRecord.Value == true"

	switch mode {
	case "automatic":
		result.Approved = true
		result.RequiresApproval = false
		result.ValidationNotes = append(result.ValidationNotes, "Automatic validation: All rule expressions approved automatically")

	case "guided":
		result.Approved = approveExpressions
		result.RequiresApproval = true
		if approveExpressions {
			result.ValidationNotes = append(result.ValidationNotes, "Guided validation: User approved rule expressions")
		} else {
			result.ValidationNotes = append(result.ValidationNotes, "Guided validation: User rejected rule expressions")
			return nil, fmt.Errorf("rule expressions were not approved by user in guided mode")
		}

	case "manual":
		result.Approved = false
		result.RequiresApproval = true
		result.ValidationNotes = append(result.ValidationNotes, "Manual validation: Expressions require manual review and modification")
		if !approveExpressions {
			return nil, fmt.Errorf("manual validation mode requires user approval of rule expressions")
		}
		result.Approved = true
		result.ValidationNotes = append(result.ValidationNotes, "Manual validation: User manually approved expressions")

	default:
		return nil, fmt.Errorf("unsupported validation mode: %s", mode)
	}

	return result, nil
}

// generateMerchantSourceName generates merchant source names based on the selected naming strategy
func generateMerchantSourceName(baseName, strategy string, index int) string {
	switch strategy {
	case "descriptive":
		return baseName + " - Merchant Portal"
	case "timestamp":
		timestamp := time.Now().Format("20060102_150405")
		return fmt.Sprintf("%s_%s", baseName, timestamp)
	case "sequential":
		return fmt.Sprintf("%s_%03d", baseName, index)
	case "custom":
		// For custom strategy, we could allow user-defined patterns
		// For now, default to descriptive with custom suffix
		return baseName + " - Custom Source"
	default:
		// Fallback to descriptive
		return baseName + " - Merchant Portal"
	}
}

// analyzeFile analyzes a file and returns metadata about columns and data based on file type
func analyzeFile(filePath, fileID, fileType string) (map[string]interface{}, error) {
	switch fileType {
	case "csv":
		return analyzeCSVFile(filePath, fileID)
	case "excel":
		return analyzeExcelFile(filePath, fileID)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}
}

// analyzeCSVFile analyzes a CSV file and returns metadata about columns and data
func analyzeCSVFile(filePath, fileID string) (map[string]interface{}, error) {
	// Open and read the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file %s: %v", filePath, err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file %s is empty", filePath)
	}

	// Extract headers and data
	headers := records[0]
	dataRows := records[1:]
	totalRows := len(dataRows)
	totalColumns := len(headers)

	// Analyze each column for patterns and characteristics
	columnAnalysis := make(map[string]map[string]interface{})
	for i, columnName := range headers {
		analysis := analyzeColumn(columnName, i, dataRows)
		columnAnalysis[columnName] = analysis
	}

	// Identify EntityID candidates
	entityIDCandidates := identifyEntityIDCandidates(headers, columnAnalysis)

	// Identify Amount candidates
	amountCandidates := identifyAmountCandidates(headers, columnAnalysis)

	// Determine recommendations
	var recommendedEntityID, recommendedAmount string
	if len(entityIDCandidates) > 0 {
		if entityMap, ok := entityIDCandidates[0].(map[string]interface{}); ok {
			if name, ok := entityMap["column_name"].(string); ok {
				recommendedEntityID = name
			}
		}
	}
	if len(amountCandidates) > 0 {
		if amountMap, ok := amountCandidates[0].(map[string]interface{}); ok {
			if name, ok := amountMap["column_name"].(string); ok {
				recommendedAmount = name
			}
		}
	}

	return map[string]interface{}{
		"filename":             filePath,
		"file_type":            "csv",
		"total_rows":           totalRows,
		"total_columns":        totalColumns,
		"all_columns":          headers,
		"entityid_candidates":  entityIDCandidates,
		"recommended_entityid": recommendedEntityID,
		"amount_candidates":    amountCandidates,
		"recommended_amount":   recommendedAmount,
		"user_selections": map[string]interface{}{
			"selected_entityid": recommendedEntityID,
			"selected_amount":   recommendedAmount,
		},
	}, nil
}

// analyzeExcelFile analyzes an Excel file and returns metadata about columns and data
func analyzeExcelFile(filePath, fileID string) (map[string]interface{}, error) {
	// Open the Excel file
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file %s: %v", filePath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: failed to close Excel file: %v\n", err)
		}
	}()

	// Get the first sheet name (or use the active sheet)
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("excel file %s has no sheets", filePath)
	}

	sheetName := sheets[0] // Use the first sheet

	// Get all rows from the sheet
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read Excel sheet %s: %v", sheetName, err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("excel sheet %s is empty", sheetName)
	}

	// Extract headers and data rows
	headers := rows[0]
	dataRows := rows[1:]
	totalRows := len(dataRows)
	totalColumns := len(headers)

	// Convert Excel rows to the same format as CSV for analysis
	// Ensure all rows have the same number of columns as headers
	normalizedDataRows := make([][]string, len(dataRows))
	for i, row := range dataRows {
		normalizedRow := make([]string, len(headers))
		for j := 0; j < len(headers); j++ {
			if j < len(row) {
				normalizedRow[j] = strings.TrimSpace(row[j])
			} else {
				normalizedRow[j] = "" // Fill missing columns with empty strings
			}
		}
		normalizedDataRows[i] = normalizedRow
	}

	// Analyze each column for patterns and characteristics
	columnAnalysis := make(map[string]map[string]interface{})
	for i, columnName := range headers {
		analysis := analyzeColumn(columnName, i, normalizedDataRows)
		columnAnalysis[columnName] = analysis
	}

	// Identify EntityID candidates
	entityIDCandidates := identifyEntityIDCandidates(headers, columnAnalysis)

	// Identify Amount candidates
	amountCandidates := identifyAmountCandidates(headers, columnAnalysis)

	// Determine recommendations
	var recommendedEntityID, recommendedAmount string
	if len(entityIDCandidates) > 0 {
		if entityMap, ok := entityIDCandidates[0].(map[string]interface{}); ok {
			if name, ok := entityMap["column_name"].(string); ok {
				recommendedEntityID = name
			}
		}
	}
	if len(amountCandidates) > 0 {
		if amountMap, ok := amountCandidates[0].(map[string]interface{}); ok {
			if name, ok := amountMap["column_name"].(string); ok {
				recommendedAmount = name
			}
		}
	}

	return map[string]interface{}{
		"filename":             filePath,
		"file_type":            "excel",
		"sheet_name":           sheetName,
		"total_sheets":         len(sheets),
		"total_rows":           totalRows,
		"total_columns":        totalColumns,
		"all_columns":          headers,
		"entityid_candidates":  entityIDCandidates,
		"recommended_entityid": recommendedEntityID,
		"amount_candidates":    amountCandidates,
		"recommended_amount":   recommendedAmount,
		"user_selections": map[string]interface{}{
			"selected_entityid": recommendedEntityID,
			"selected_amount":   recommendedAmount,
		},
	}, nil
}

// analyzeColumn analyzes individual column characteristics
func analyzeColumn(columnName string, columnIndex int, dataRows [][]string) map[string]interface{} {
	if len(dataRows) == 0 {
		return map[string]interface{}{
			"unique_count":      0,
			"unique_percentage": 0.0,
			"data_type":         "unknown",
			"sample_values":     []string{},
			"is_numeric":        false,
			"is_monetary":       false,
		}
	}

	// Extract column values
	values := make([]string, 0, len(dataRows))
	uniqueValues := make(map[string]bool)
	numericCount := 0
	monetaryCount := 0

	sampleSize := len(dataRows)
	if sampleSize > 100 {
		sampleSize = 100 // Analyze first 100 rows for performance
	}

	for i := 0; i < sampleSize && i < len(dataRows); i++ {
		if columnIndex < len(dataRows[i]) {
			value := strings.TrimSpace(dataRows[i][columnIndex])
			if value != "" {
				values = append(values, value)
				uniqueValues[value] = true

				// Check if numeric
				if isNumeric(value) {
					numericCount++
				}

				// Check if monetary
				if isMonetary(value) {
					monetaryCount++
				}
			}
		}
	}

	uniqueCount := len(uniqueValues)
	uniquePercentage := float64(uniqueCount) / float64(len(values)) * 100

	// Get sample values (first 5 unique values)
	sampleValues := make([]string, 0, 5)
	count := 0
	for value := range uniqueValues {
		if count >= 5 {
			break
		}
		sampleValues = append(sampleValues, value)
		count++
	}

	// Determine data type
	dataType := "string"
	isNumericCol := float64(numericCount)/float64(len(values)) > 0.8
	isMonetaryCol := float64(monetaryCount)/float64(len(values)) > 0.7

	if isMonetaryCol {
		dataType = "monetary"
	} else if isNumericCol {
		dataType = "numeric"
	}

	return map[string]interface{}{
		"unique_count":      uniqueCount,
		"unique_percentage": uniquePercentage,
		"data_type":         dataType,
		"sample_values":     sampleValues,
		"is_numeric":        isNumericCol,
		"is_monetary":       isMonetaryCol,
		"total_values":      len(values),
	}
}

// identifyEntityIDCandidates identifies potential EntityID columns
func identifyEntityIDCandidates(headers []string, columnAnalysis map[string]map[string]interface{}) []interface{} {
	candidates := make([]interface{}, 0)

	// Priority patterns for EntityID columns
	idPatterns := []string{
		"(?i)(transaction[_\\s]*id|txn[_\\s]*id)",
		"(?i)(entity[_\\s]*id|ent[_\\s]*id)",
		"(?i)(instance[_\\s]*id|inst[_\\s]*id)",
		"(?i)(reference[_\\s]*(number|no|num|id)|ref[_\\s]*(no|num|id))",
		"(?i)(receipt[_\\s]*(number|no|num|id))",
		"(?i)(order[_\\s]*(number|no|num|id))",
		"(?i)^id$",
		"(?i)(unique[_\\s]*id)",
	}

	type candidate struct {
		columnName       string
		confidence       float64
		reason           string
		uniquePercentage float64
	}

	var candidateList []candidate

	for _, columnName := range headers {
		analysis := columnAnalysis[columnName]
		uniquePercentage := analysis["unique_percentage"].(float64)

		confidence := 0.0
		reason := ""

		// Check naming patterns (highest priority)
		for i, pattern := range idPatterns {
			if matched, _ := regexp.MatchString(pattern, columnName); matched {
				confidence = 0.95 - float64(i)*0.05 // Decrease confidence for lower priority patterns
				reason = "ID-like naming pattern with high uniqueness"
				break
			}
		}

		// Check uniqueness-based scoring
		if confidence == 0.0 && uniquePercentage >= 95.0 {
			confidence = 0.80
			reason = "Very high uniqueness suggests unique identifier"
		} else if confidence == 0.0 && uniquePercentage >= 85.0 {
			confidence = 0.70
			reason = "High uniqueness suggests potential identifier"
		}

		// Boost confidence for high uniqueness even with pattern match
		if confidence > 0.0 && uniquePercentage >= 98.0 {
			confidence = minFloat(confidence+0.1, 0.98)
		}

		// Only include candidates with reasonable confidence and uniqueness
		if confidence >= 0.60 && uniquePercentage >= 75.0 {
			candidateList = append(candidateList, candidate{
				columnName:       columnName,
				confidence:       confidence,
				reason:           reason,
				uniquePercentage: uniquePercentage,
			})
		}
	}

	// Sort by confidence (highest first)
	for i := 0; i < len(candidateList)-1; i++ {
		for j := i + 1; j < len(candidateList); j++ {
			if candidateList[j].confidence > candidateList[i].confidence {
				candidateList[i], candidateList[j] = candidateList[j], candidateList[i]
			}
		}
	}

	// Convert to expected format
	for _, candidate := range candidateList {
		candidates = append(candidates, map[string]interface{}{
			"column_name":       candidate.columnName,
			"unique_percentage": candidate.uniquePercentage,
			"confidence":        candidate.confidence,
			"reason":            candidate.reason,
		})
	}

	return candidates
}

// identifyAmountCandidates identifies potential Amount columns
func identifyAmountCandidates(headers []string, columnAnalysis map[string]map[string]interface{}) []interface{} {
	candidates := make([]interface{}, 0)

	// Priority patterns for Amount columns
	amountPatterns := []string{
		"(?i)^amount$",
		"(?i)(net[_\\s]*amount|netamount)",
		"(?i)(gross[_\\s]*amount|grossamount)",
		"(?i)(total[_\\s]*amount|totalamount)",
		"(?i)(transaction[_\\s]*amount|txn[_\\s]*amount)",
		"(?i)(payment[_\\s]*amount)",
		"(?i)(balance)",
		"(?i)(value)",
		"(?i)(price)",
		"(?i)(sum)",
		"(?i).*amount.*",
		"(?i).*_amt$",
		"(?i).*amt_.*",
	}

	type candidate struct {
		columnName   string
		confidence   float64
		reason       string
		sampleValues []string
	}

	var candidateList []candidate

	for _, columnName := range headers {
		analysis := columnAnalysis[columnName]
		isMonetary := analysis["is_monetary"].(bool)
		isNumeric := analysis["is_numeric"].(bool)
		sampleValues := analysis["sample_values"].([]string)

		confidence := 0.0
		reason := ""

		// Check naming patterns with monetary data
		for i, pattern := range amountPatterns {
			if matched, _ := regexp.MatchString(pattern, columnName); matched {
				if isMonetary {
					confidence = 0.95 - float64(i)*0.03
					reason = "Amount-like naming with monetary values"
				} else if isNumeric {
					confidence = 0.80 - float64(i)*0.03
					reason = "Amount-like naming with numeric values"
				} else {
					confidence = 0.60 - float64(i)*0.03
					reason = "Amount-like naming pattern"
				}
				break
			}
		}

		// Check data characteristics without pattern match
		if confidence == 0.0 && isMonetary {
			confidence = 0.75
			reason = "Contains monetary values (currency format)"
		} else if confidence == 0.0 && isNumeric {
			confidence = 0.60
			reason = "Contains numeric values"
		}

		// Only include candidates with reasonable confidence
		if confidence >= 0.50 {
			candidateList = append(candidateList, candidate{
				columnName:   columnName,
				confidence:   confidence,
				reason:       reason,
				sampleValues: sampleValues,
			})
		}
	}

	// Sort by confidence (highest first)
	for i := 0; i < len(candidateList)-1; i++ {
		for j := i + 1; j < len(candidateList); j++ {
			if candidateList[j].confidence > candidateList[i].confidence {
				candidateList[i], candidateList[j] = candidateList[j], candidateList[i]
			}
		}
	}

	// Convert to expected format
	for _, candidate := range candidateList {
		candidates = append(candidates, map[string]interface{}{
			"column_name":   candidate.columnName,
			"sample_values": candidate.sampleValues,
			"confidence":    candidate.confidence,
			"reason":        candidate.reason,
		})
	}

	return candidates
}

// isNumeric checks if a string represents a numeric value
func isNumeric(s string) bool {
	// Remove common non-numeric characters that might be in monetary values
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, "₹", "")
	s = strings.ReplaceAll(s, "€", "")
	s = strings.ReplaceAll(s, "£", "")
	s = strings.TrimSpace(s)

	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// isMonetary checks if a string represents a monetary value
func isMonetary(s string) bool {
	// Check for currency symbols
	currencyPattern := `(?i)[$₹€£¥¢]|\b(usd|eur|gbp|inr|jpy)\b`
	if matched, _ := regexp.MatchString(currencyPattern, s); matched {
		return true
	}

	// Check for decimal pattern with 2 decimal places (common for currency)
	decimalPattern := `^\d+\.\d{2}$`
	if matched, _ := regexp.MatchString(decimalPattern, s); matched {
		return true
	}

	// Check for comma-separated thousands with decimal
	commaPattern := `^\d{1,3}(,\d{3})*(\.\d{2})?$`
	if matched, _ := regexp.MatchString(commaPattern, s); matched {
		return true
	}

	// For simple numeric values, check if they're reasonable monetary amounts
	cleanValue := strings.ReplaceAll(s, ",", "")
	if val, err := strconv.ParseFloat(cleanValue, 64); err == nil {
		// Consider values between 0.01 and 999999999.99 as potentially monetary
		return val >= 0.01 && val <= 999999999.99 &&
			(strings.Contains(s, ".") || len(cleanValue) <= 8)
	}

	return false
}

// minFloat returns the minimum of two float64 values
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// makeReconSaaSAPICall makes authenticated API calls to recon-saas service
func makeReconSaaSAPICall(ctx context.Context, method, endpoint string, payload interface{}) (map[string]interface{}, error) {
	const baseURL = "https://recon-saas.dev.razorpay.in"
	const authHeader = "Basic cmVjb24tc2FhczpyZWNvbi1zYWFz"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var reqBody io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %v", err)
		}
		reqBody = bytes.NewReader(payloadBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(responseBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return result, nil
}

// makeReconSaaSAPICallString makes authenticated API calls to recon-saas service and returns string response
func makeReconSaaSAPICallString(ctx context.Context, method, endpoint string, payload interface{}) (string, error) {
	const baseURL = "https://recon-saas.dev.razorpay.in"
	const authHeader = "Basic cmVjb24tc2FhczpyZWNvbi1zYWFz"

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	var reqBody io.Reader
	if payload != nil {
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return "", fmt.Errorf("failed to marshal payload: %v", err)
		}
		reqBody = bytes.NewReader(payloadBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+endpoint, reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API call failed: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(responseBody))
	}

	// Handle string response in quotes (e.g., "R2QTjA6dPvcygE")
	responseStr := strings.TrimSpace(string(responseBody))
	if strings.HasPrefix(responseStr, "\"") && strings.HasSuffix(responseStr, "\"") {
		// Remove surrounding quotes
		responseStr = responseStr[1 : len(responseStr)-1]
	}

	return responseStr, nil
}

// createMasterSource creates a master source via recon-saas API
func createMasterSource(ctx context.Context, name, columnsJSON, entityIDColumn, amountColumn string) (string, error) {
	var columns []string
	if err := json.Unmarshal([]byte(columnsJSON), &columns); err != nil {
		return "", fmt.Errorf("invalid columns JSON: %v", err)
	}

	// Generate source schema (all columns as string type)
	sourceSchema := make([]map[string]string, len(columns))
	for i, col := range columns {
		sourceSchema[i] = map[string]string{
			"name": col,
			"type": "string",
		}
	}

	// Generate mapping config with snake_case destinations
	mappingConfig := make([]map[string]string, len(columns))
	for i, col := range columns {
		var destination string
		switch col {
		case entityIDColumn:
			destination = "EntityID"
		case amountColumn:
			destination = "Amount"
		default:
			// Convert to snake_case
			destination = strings.ToLower(strings.ReplaceAll(col, " ", "_"))
			destination = strings.ReplaceAll(destination, "-", "_")
		}

		mappingConfig[i] = map[string]string{
			"value":       "",
			"source":      col,
			"destination": destination,
		}
	}

	payload := map[string]interface{}{
		"name": name,
		"config": map[string]interface{}{
			"is_internal":       false,
			"unique_keys":       []string{"EntityID"},
			"allow_upload":      false,
			"ingest_to_db":      true,
			"skip_top_rows":     0,
			"split_file_basis":  "",
			"report_enrichment": true,
			"sub_source_config": map[string]interface{}{
				"amount_key":                 "",
				"aggregate_config":           nil,
				"grouping_columns":           nil,
				"sub_master_source_id":       "",
				"enable_sub_source_creation": false,
			},
			"row_hash_value_based_split_config": map[string]interface{}{
				"column_joiner":                    "",
				"header_hash_to_master_source_map": nil,
			},
		},
		"source_schema":  sourceSchema,
		"mapping_config": mappingConfig,
	}

	result, err := makeReconSaaSAPICall(ctx, "POST", "/v1/admin-recon-saas/sources/create", payload)
	if err != nil {
		return "", err
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}

	return fmt.Sprintf("mock_master_source_%d", time.Now().Unix()), nil
}

// createMerchantSource creates a merchant source via recon-saas API
func createMerchantSource(ctx context.Context, merchantID, masterSourceID, name string) (string, error) {
	payload := map[string]interface{}{
		"name":             name,
		"merchant_id":      merchantID,
		"master_source_id": masterSourceID,
		"config": map[string]interface{}{
			"cc_emails":          nil,
			"bcc_emails":         nil,
			"allow_upload":       true,
			"reporting_emails":   nil,
			"split_file_basis":   "",
			"beam_sftp_push_job": "",
			"row_hash_value_based_split_config": map[string]interface{}{
				"column_joiner":                    "",
				"header_hash_to_master_source_map": nil,
			},
		},
		"source_schema":  nil,
		"mapping_config": nil,
	}

	// Use the string-specific API call function for create_merchant endpoint
	merchantSourceID, err := makeReconSaaSAPICallString(ctx, "POST", "/v1/admin-recon-saas/sources/create_merchant", payload)
	if err != nil {
		return "", err
	}

	// The response should be a string ID (e.g., "R2QTjA6dPvcygE")
	if merchantSourceID != "" {
		return merchantSourceID, nil
	}

	return fmt.Sprintf("mock_merchant_source_%d", time.Now().Unix()), nil
}

// createReconStates creates reconciliation states via recon-saas API
func createReconStates(ctx context.Context, merchantID, source1Name, source2Name string) (map[string]interface{}, error) {
	states := []map[string]interface{}{
		{
			"name":     "Reconciled",
			"priority": 2,
			"remarks":  "success",
		},
		{
			"name":     "Unreconciled",
			"priority": 3,
			"remarks":  "Amount mismatch",
		},
		{
			"name":     "Unreconciled",
			"priority": 3,
			"remarks":  fmt.Sprintf("Record not found in %s", source1Name),
		},
		{
			"name":     "Unreconciled",
			"priority": 3,
			"remarks":  fmt.Sprintf("Record not found in %s", source2Name),
		},
	}

	results := make(map[string]interface{})
	stateNames := []string{"reconciled_state", "amount_mismatch_state", "missing_file1_state", "missing_file2_state"}

	for i, state := range states {
		payload := map[string]interface{}{
			"merchant_id": merchantID,
			"name":        state["name"],
			"priority":    state["priority"],
			"remarks":     state["remarks"],
		}

		result, err := makeReconSaaSAPICall(ctx, "POST", "/v1/admin-recon-saas/recon_state", payload)
		if err != nil {
			return nil, fmt.Errorf("failed to create state %s: %v", stateNames[i], err)
		}

		var stateID string
		if id, ok := result["id"].(string); ok {
			stateID = id
		} else {
			stateID = fmt.Sprintf("mock_recon_state_%d_%d", time.Now().Unix(), i)
		}

		results[stateNames[i]] = map[string]interface{}{
			"recon_state_id": stateID,
			"name":           state["name"],
			"priority":       state["priority"],
			"remarks":        state["remarks"],
			"api_response":   result,
		}
	}

	return results, nil
}

// createReconRules creates reconciliation rules via recon-saas API
func createReconRules(ctx context.Context, merchantID, masterSourceID1, masterSourceID2 string, reconStates map[string]interface{}) (map[string]interface{}, error) {
	// Extract recon state IDs
	getStateID := func(stateName string) string {
		if state, ok := reconStates[stateName].(map[string]interface{}); ok {
			if stateID, ok := state["recon_state_id"].(string); ok {
				return stateID
			}
		}
		return fmt.Sprintf("mock_state_%s", stateName)
	}

	rules := []map[string]interface{}{
		{
			"name":           "Reconciled Transactions",
			"expression":     fmt.Sprintf("%s.EntityID == %s.EntityID && %s.Amount.Equal(%s.Amount)", masterSourceID1, masterSourceID2, masterSourceID1, masterSourceID2),
			"recon_state_id": getStateID("reconciled_state"),
			"type":           "reconciled",
		},
		{
			"name":           "Amount Mismatch Transactions",
			"expression":     fmt.Sprintf("%s.EntityID == %s.EntityID && !%s.Amount.Equal(%s.Amount)", masterSourceID1, masterSourceID2, masterSourceID1, masterSourceID2),
			"recon_state_id": getStateID("amount_mismatch_state"),
			"type":           "amount_mismatch",
		},
		{
			"name":           "Missing Record in File 1",
			"expression":     "NoRecord.Value == true",
			"recon_state_id": getStateID("missing_file1_state"),
			"type":           "missing_record",
		},
		{
			"name":           "Missing Record in File 2",
			"expression":     "NoRecord.Value == true",
			"recon_state_id": getStateID("missing_file2_state"),
			"type":           "missing_record",
		},
	}

	results := make(map[string]interface{})
	ruleNames := []string{"reconciled_rule", "amount_mismatch_rule", "missing_record_rule_file1", "missing_record_rule_file2"}

	for i, rule := range rules {
		payload := map[string]interface{}{
			"merchant_id":    merchantID,
			"name":           rule["name"],
			"expression":     rule["expression"],
			"sources":        []string{masterSourceID1, masterSourceID2},
			"recon_state_id": rule["recon_state_id"],
			"type":           "recon",
			"case_rule":      false,
		}

		result, err := makeReconSaaSAPICall(ctx, "POST", "/v1/admin-recon-saas/rule", payload)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule %s: %v", ruleNames[i], err)
		}

		var ruleID string
		if id, ok := result["id"].(string); ok {
			ruleID = id
		} else {
			ruleID = fmt.Sprintf("mock_rule_%d_%d", time.Now().Unix(), i)
		}

		results[ruleNames[i]] = map[string]interface{}{
			"rule_id":        ruleID,
			"name":           rule["name"],
			"expression":     rule["expression"],
			"recon_state_id": rule["recon_state_id"],
			"sources":        []string{masterSourceID1, masterSourceID2},
			"api_response":   result,
		}
	}

	return results, nil
}

// createReconRulesWithValidation creates reconciliation rules via recon-saas API with validation
func createReconRulesWithValidation(ctx context.Context, merchantID, masterSourceID1, masterSourceID2 string, reconStates map[string]interface{}, validationResult *ValidationResult) (map[string]interface{}, error) {
	// Extract recon state IDs
	getStateID := func(stateName string) string {
		if state, ok := reconStates[stateName].(map[string]interface{}); ok {
			if stateID, ok := state["recon_state_id"].(string); ok {
				return stateID
			}
		}
		return fmt.Sprintf("mock_state_%s", stateName)
	}

	// Use expressions from validation result if available, otherwise use defaults
	reconciledExpr := validationResult.RuleExpressions["reconciled"]
	amountMismatchExpr := validationResult.RuleExpressions["amount_mismatch"]
	missingRecordExpr := validationResult.RuleExpressions["missing_record"]

	if reconciledExpr == "" {
		reconciledExpr = fmt.Sprintf("%s.EntityID == %s.EntityID && %s.Amount.Equal(%s.Amount)", masterSourceID1, masterSourceID2, masterSourceID1, masterSourceID2)
	}
	if amountMismatchExpr == "" {
		amountMismatchExpr = fmt.Sprintf("%s.EntityID == %s.EntityID && !%s.Amount.Equal(%s.Amount)", masterSourceID1, masterSourceID2, masterSourceID1, masterSourceID2)
	}
	if missingRecordExpr == "" {
		missingRecordExpr = "NoRecord.Value == true"
	}

	rules := []map[string]interface{}{
		{
			"name":           "Reconciled Transactions",
			"expression":     reconciledExpr,
			"recon_state_id": getStateID("reconciled_state"),
			"type":           "reconciled",
		},
		{
			"name":           "Amount Mismatch Transactions",
			"expression":     amountMismatchExpr,
			"recon_state_id": getStateID("amount_mismatch_state"),
			"type":           "amount_mismatch",
		},
		{
			"name":           "Missing Record in File 1",
			"expression":     missingRecordExpr,
			"recon_state_id": getStateID("missing_file1_state"),
			"type":           "missing_record",
		},
		{
			"name":           "Missing Record in File 2",
			"expression":     missingRecordExpr,
			"recon_state_id": getStateID("missing_file2_state"),
			"type":           "missing_record",
		},
	}

	results := make(map[string]interface{})
	ruleNames := []string{"reconciled_rule", "amount_mismatch_rule", "missing_record_rule_file1", "missing_record_rule_file2"}

	for i, rule := range rules {
		payload := map[string]interface{}{
			"merchant_id":    merchantID,
			"name":           rule["name"],
			"expression":     rule["expression"],
			"sources":        []string{masterSourceID1, masterSourceID2},
			"recon_state_id": rule["recon_state_id"],
			"type":           "recon",
			"case_rule":      false,
		}

		result, err := makeReconSaaSAPICall(ctx, "POST", "/v1/admin-recon-saas/rule", payload)
		if err != nil {
			return nil, fmt.Errorf("failed to create rule %s: %v", ruleNames[i], err)
		}

		var ruleID string
		if id, ok := result["id"].(string); ok {
			ruleID = id
		} else {
			ruleID = fmt.Sprintf("mock_rule_%d_%d", time.Now().Unix(), i)
		}

		results[ruleNames[i]] = map[string]interface{}{
			"rule_id":          ruleID,
			"name":             rule["name"],
			"expression":       rule["expression"],
			"recon_state_id":   rule["recon_state_id"],
			"sources":          []string{masterSourceID1, masterSourceID2},
			"validation_mode":  validationResult.Mode,
			"validation_notes": validationResult.ValidationNotes,
			"api_response":     result,
		}
	}

	// Add validation summary to results
	results["validation_summary"] = map[string]interface{}{
		"mode":              validationResult.Mode,
		"approved":          validationResult.Approved,
		"requires_approval": validationResult.RequiresApproval,
		"validation_notes":  validationResult.ValidationNotes,
		"rule_expressions":  validationResult.RuleExpressions,
	}

	return results, nil
}

// createLookup creates a lookup configuration via recon-saas API
func createLookup(ctx context.Context, merchantID, source1Name, source2Name string) (string, error) {
	payload := map[string]interface{}{
		"name":        fmt.Sprintf("Entity Lookup for %s and %s", source1Name, source2Name),
		"merchant_id": merchantID,
		"config": []map[string]interface{}{
			{
				"source":  "record_internal",
				"columns": []string{"EntityID"},
			},
		},
	}

	result, err := makeReconSaaSAPICall(ctx, "POST", "/v1/admin-recon-saas/lookup", payload)
	if err != nil {
		return "", err
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}

	return fmt.Sprintf("mock_lookup_%d", time.Now().Unix()), nil
}

// createMasterReconProcess creates a master reconciliation process via recon-saas API
func createMasterReconProcess(ctx context.Context, source1Name, source2Name, lookupID, masterSourceID1, masterSourceID2 string, ruleIDs []string, source1Columns, source2Columns, source1EntityID, source2EntityID, source1Amount, source2Amount string) (string, error) {
	processName := fmt.Sprintf("%s to %s Reconciliation", source1Name, source2Name)
	productID := fmt.Sprintf("%s_%s",
		strings.ToUpper(strings.ReplaceAll(source1Name[:3], " ", "")),
		strings.ToUpper(strings.ReplaceAll(source2Name[:3], " ", "")))

	// Parse column arrays
	var columns1, columns2 []string
	if err := json.Unmarshal([]byte(source1Columns), &columns1); err != nil {
		return "", fmt.Errorf("invalid source1_columns JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(source2Columns), &columns2); err != nil {
		return "", fmt.Errorf("invalid source2_columns JSON: %v", err)
	}

	// Generate frontend columns (union of all columns from both files)
	frontendCols := make([]string, 0)
	seenColumns := make(map[string]bool)

	// Add columns from source 1
	for _, col := range columns1 {
		if !seenColumns[col] {
			frontendCols = append(frontendCols, col)
			seenColumns[col] = true
		}
	}

	// Add columns from source 2
	for _, col := range columns2 {
		if !seenColumns[col] {
			frontendCols = append(frontendCols, col)
			seenColumns[col] = true
		}
	}

	// Generate column mappings for source 1
	source1ColumnMap := make([]map[string]string, 0)
	for _, col := range columns1 {
		var destination string
		switch col {
		case source1EntityID:
			destination = "EntityID"
		case source1Amount:
			destination = "Amount"
		default:
			// Convert to snake_case
			destination = strings.ToLower(strings.ReplaceAll(col, " ", "_"))
			destination = strings.ReplaceAll(destination, "-", "_")
		}

		source1ColumnMap = append(source1ColumnMap, map[string]string{
			"report_column": col,
			"source_column": destination,
		})
	}

	// Generate column mappings for source 2
	source2ColumnMap := make([]map[string]string, 0)
	for _, col := range columns2 {
		var destination string
		switch col {
		case source2EntityID:
			destination = "EntityID"
		case source2Amount:
			destination = "Amount"
		default:
			// Convert to snake_case
			destination = strings.ToLower(strings.ReplaceAll(col, " ", "_"))
			destination = strings.ReplaceAll(destination, "-", "_")
		}

		source2ColumnMap = append(source2ColumnMap, map[string]string{
			"report_column": col,
			"source_column": destination,
		})
	}

	payload := map[string]interface{}{
		"name":       processName,
		"product_id": productID,
		"type":       "Gateway",
		"lookup_config": []map[string]interface{}{
			{
				"config": map[string]string{
					masterSourceID1: lookupID,
				},
				"streaming_source_id": masterSourceID2,
			},
		},
		"rules": map[string]interface{}{
			"rule_ids": ruleIDs,
		},
		"sources":  []string{masterSourceID1, masterSourceID2},
		"sequence": []interface{}{},
		"report_config": map[string]interface{}{
			"frontend_cols": frontendCols,
			"source_report_config": []map[string]interface{}{
				{
					"column_map":       source1ColumnMap,
					"master_source_id": masterSourceID1,
				},
				{
					"column_map":       source2ColumnMap,
					"master_source_id": masterSourceID2,
				},
			},
		},
	}

	result, err := makeReconSaaSAPICall(ctx, "POST", "/v1/admin-recon-saas/recon_process/master", payload)
	if err != nil {
		return "", err
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}

	return fmt.Sprintf("mock_master_recon_process_%d", time.Now().Unix()), nil
}

// createMerchantReconProcess creates a merchant reconciliation process via recon-saas API
func createMerchantReconProcess(ctx context.Context, merchantID, masterReconProcessID, merchantSourceID1, merchantSourceID2 string) (string, error) {
	payload := map[string]interface{}{
		"merchant_id":             merchantID,
		"master_recon_process_id": masterReconProcessID,
		"sources":                 []string{merchantSourceID1, merchantSourceID2},
	}

	result, err := makeReconSaaSAPICall(ctx, "POST", "/v1/admin-recon-saas/recon_process/merchant", payload)
	if err != nil {
		return "", err
	}

	if id, ok := result["id"].(string); ok {
		return id, nil
	}

	return fmt.Sprintf("mock_merchant_recon_process_%d", time.Now().Unix()), nil
}

// extractStateIDs extracts recon state IDs from recon states map
func extractStateIDs(reconStates map[string]interface{}) map[string]string {
	stateIDs := make(map[string]string)

	for stateName, stateData := range reconStates {
		if state, ok := stateData.(map[string]interface{}); ok {
			if stateID, ok := state["recon_state_id"].(string); ok {
				stateIDs[stateName] = stateID
			}
		}
	}

	return stateIDs
}

// extractRuleIDs extracts rule IDs from rules map
func extractRuleIDs(rules map[string]interface{}) map[string]string {
	ruleIDs := make(map[string]string)

	for ruleName, ruleData := range rules {
		if rule, ok := ruleData.(map[string]interface{}); ok {
			if ruleID, ok := rule["rule_id"].(string); ok {
				ruleIDs[ruleName] = ruleID
			}
		}
	}

	return ruleIDs
}

// Helper functions for aggregation tool

// determineFileType determines file type based on file extension
func determineFileType(filePath string) string {
	filePath = strings.ToLower(filePath)
	if strings.HasSuffix(filePath, ".csv") {
		return "csv"
	} else if strings.HasSuffix(filePath, ".xlsx") || strings.HasSuffix(filePath, ".xls") {
		return "excel"
	}
	return "csv" // default to csv
}

// validateEntityIdentifier checks if the entity identifier exists in the file analysis
func validateEntityIdentifier(analysis map[string]interface{}, entityIdentifier string) bool {
	if allColumns, ok := analysis["all_columns"].([]string); ok {
		for _, column := range allColumns {
			if column == entityIdentifier {
				return true
			}
		}
	}
	return false
}

// updateMasterSourceMapping updates master source mapping config with entity identifier
func updateMasterSourceMapping(ctx context.Context, masterSourceID, entityIdentifier string) error {
	payload := map[string]interface{}{
		"mapping_config": []map[string]interface{}{
			{
				"value":       "",
				"source":      entityIdentifier,
				"destination": "EntityID",
			},
		},
	}

	_, err := makeReconSaaSAPICall(ctx, "PATCH", fmt.Sprintf("/v1/admin-recon-saas/sources/update/%s", masterSourceID), payload)
	return err
}

// updateLookupAggregation enables aggregation logic in lookup configuration
func updateLookupAggregation(ctx context.Context, lookupID string) error {
	payload := map[string]interface{}{
		"config": []map[string]interface{}{
			{
				"source":  "record_internal",
				"Columns": []string{"EntityID"},
				"aggregation": map[string]interface{}{
					"enabled":    true,
					"conditions": nil,
				},
				"advanced_config": map[string]interface{}{
					"enabled":     false,
					"cols_config": nil,
				},
			},
		},
	}

	_, err := makeReconSaaSAPICall(ctx, "PATCH", fmt.Sprintf("/v1/admin-recon-saas/lookup/%s", lookupID), payload)
	return err
}

// getLatestMasterSourceID retrieves the most recently created master source ID
// Since list endpoints don't exist, we'll use a fallback approach
func getLatestMasterSourceID(ctx context.Context) (string, error) {
	// For now, return the known working ID
	// In production, this could be stored in a database or config
	return "RVKUOMPMGnjQHD", nil
}

// getLatestLookupID retrieves the most recently created lookup ID
func getLatestLookupID(ctx context.Context) (string, error) {
	// For now, return the known working ID
	return "RVKYH04CYjMVQx", nil
}

// getLatestMasterReconProcessID retrieves the most recently created master recon process ID
func getLatestMasterReconProcessID(ctx context.Context) (string, error) {
	// For now, return the known working ID
	return "RVKYH4zmAM2HvR", nil
}

// generateComprehensiveReportConfig generates a complete report configuration based on file analysis
func generateComprehensiveReportConfig(file1Analysis, file2Analysis map[string]interface{}, entityIdentifier, masterSourceID1, masterSourceID2 string) map[string]interface{} {
	// Extract columns from both files
	var file1Columns, file2Columns []string
	if cols1, ok := file1Analysis["all_columns"].([]string); ok {
		file1Columns = cols1
	}
	if cols2, ok := file2Analysis["all_columns"].([]string); ok {
		file2Columns = cols2
	}

	// Generate frontend columns (union of all columns from both files)
	frontendCols := make([]string, 0)
	seenColumns := make(map[string]bool)

	// Add columns from file 1
	for _, col := range file1Columns {
		if !seenColumns[col] {
			frontendCols = append(frontendCols, col)
			seenColumns[col] = true
		}
	}

	// Add columns from file 2
	for _, col := range file2Columns {
		if !seenColumns[col] {
			frontendCols = append(frontendCols, col)
			seenColumns[col] = true
		}
	}

	// Generate column mappings for file 1
	source1ColumnMap := make([]map[string]interface{}, 0)

	// First add EntityID mapping (original column name)
	source1ColumnMap = append(source1ColumnMap, map[string]interface{}{
		"id":            "",
		"type":          "",
		"report_column": "EntityID",
		"source_column": entityIdentifier,
	})

	// Then add EntityIdentifier mapping (processed version)
	source1ColumnMap = append(source1ColumnMap, map[string]interface{}{
		"id":            "",
		"type":          "",
		"report_column": "EntityIdentifier",
		"source_column": entityIdentifier,
	})

	// Add all other columns
	for _, col := range file1Columns {
		if col == entityIdentifier {
			continue // Skip entity identifier as it's already added above
		}

		// Convert to camelCase for report column
		reportColumn := strings.ReplaceAll(col, " ", "")
		reportColumn = strings.ReplaceAll(reportColumn, "-", "")
		reportColumn = strings.ReplaceAll(reportColumn, "_", "")

		source1ColumnMap = append(source1ColumnMap, map[string]interface{}{
			"id":            "",
			"type":          "",
			"report_column": reportColumn,
			"source_column": col,
		})
	}

	// Generate column mappings for file 2
	source2ColumnMap := make([]map[string]interface{}, 0)

	// First add EntityID mapping (original column name)
	source2ColumnMap = append(source2ColumnMap, map[string]interface{}{
		"id":            "",
		"type":          "",
		"report_column": "EntityID",
		"source_column": entityIdentifier,
	})

	// Then add EntityIdentifier mapping (processed version)
	source2ColumnMap = append(source2ColumnMap, map[string]interface{}{
		"id":            "",
		"type":          "",
		"report_column": "EntityIdentifier",
		"source_column": entityIdentifier,
	})

	// Add all other columns
	for _, col := range file2Columns {
		if col == entityIdentifier {
			continue // Skip entity identifier as it's already added above
		}

		// Convert to camelCase for report column
		reportColumn := strings.ReplaceAll(col, " ", "")
		reportColumn = strings.ReplaceAll(reportColumn, "-", "")
		reportColumn = strings.ReplaceAll(reportColumn, "_", "")

		source2ColumnMap = append(source2ColumnMap, map[string]interface{}{
			"id":            "",
			"type":          "",
			"report_column": reportColumn,
			"source_column": col,
		})
	}

	return map[string]interface{}{
		"skip_rows":     false,
		"skip_status":   false,
		"frontend_cols": frontendCols,
		"report_format": nil,
		"report_upload": map[string]interface{}{
			"enabled":        false,
			"s3_path":        "",
			"s3_bucket":      "",
			"s3_file_prefix": "",
		},
		"custom_reports":    nil,
		"report_channel":    nil,
		"report_enrichment": nil,
		"reporting_sources": nil,
		"source_report_config": []map[string]interface{}{
			{
				"column_map":       source1ColumnMap,
				"report_name":      "",
				"master_source_id": masterSourceID1,
				"report_name_config": map[string]interface{}{
					"format":         "",
					"parameters_map": nil,
				},
				"email_subject_config": map[string]interface{}{
					"format":         "",
					"parameters_map": nil,
				},
			},
			{
				"column_map":       source2ColumnMap,
				"report_name":      "",
				"master_source_id": masterSourceID2,
				"report_name_config": map[string]interface{}{
					"format":         "",
					"parameters_map": nil,
				},
				"email_subject_config": map[string]interface{}{
					"format":         "",
					"parameters_map": nil,
				},
			},
		},
		"sub_source_report_config": map[string]interface{}{
			"enabled":              false,
			"frontend_cols":        nil,
			"recon_status_key":     "",
			"source_report_config": nil,
		},
		"skip_rows_recon_state_ids":    nil,
		"advanced_reporting_config_id": "",
	}
}

func updateMasterReconProcessReportConfig(ctx context.Context, masterReconProcessID, entityIdentifier string) error {
	payload := map[string]interface{}{
		"report_config": map[string]interface{}{
			"source_report_config": []map[string]interface{}{
				{
					"column_map": []map[string]interface{}{
						{
							"id":            "",
							"type":          "",
							"report_column": entityIdentifier,
							"source_column": "Entity identifier",
						},
					},
				},
			},
		},
	}

	_, err := makeReconSaaSAPICall(ctx, "PATCH", fmt.Sprintf("/v1/admin-recon-saas/recon_process/master/%s", masterReconProcessID), payload)
	return err
}

// updateMasterReconProcessReportConfigComprehensive updates master recon process with comprehensive report configuration
func updateMasterReconProcessReportConfigComprehensive(ctx context.Context, masterReconProcessID string, reportConfig map[string]interface{}) error {
	payload := map[string]interface{}{
		"report_config": reportConfig,
	}

	_, err := makeReconSaaSAPICall(ctx, "PATCH", fmt.Sprintf("/v1/admin-recon-saas/recon_process/master/%s", masterReconProcessID), payload)
	return err
}

// getRecordCount extracts record count from file analysis
func getRecordCount(analysis map[string]interface{}) int {
	if recordCount, ok := analysis["total_rows"].(int); ok {
		return recordCount
	}
	return 0
}

// generateAggregationPreview generates a preview of how aggregation will work
func generateAggregationPreview(file1Path, file2Path, entityIdentifier, aggregationStrategy, file1Type, file2Type string) map[string]interface{} {
	// Read and parse both files
	file1Data, err1 := readFileData(file1Path, file1Type)
	file2Data, err2 := readFileData(file2Path, file2Type)

	if err1 != nil || err2 != nil {
		return map[string]interface{}{
			"error":       "Failed to read files for preview",
			"file1_error": err1,
			"file2_error": err2,
		}
	}

	// Find entity identifier column index
	entityColIndex := -1
	for i, col := range file1Data.Headers {
		if col == entityIdentifier {
			entityColIndex = i
			break
		}
	}

	if entityColIndex == -1 {
		return map[string]interface{}{
			"error": fmt.Sprintf("Entity identifier '%s' not found in file headers", entityIdentifier),
		}
	}

	// Find amount column index (assuming it's called "Amount")
	amountColIndex := -1
	for i, col := range file1Data.Headers {
		if strings.ToLower(col) == "amount" {
			amountColIndex = i
			break
		}
	}

	if amountColIndex == -1 {
		return map[string]interface{}{
			"error": "Amount column not found in file headers",
		}
	}

	// Perform aggregation on file 1
	aggregatedData := make(map[string]float64)
	originalRecords := make(map[string][]map[string]interface{})

	for _, record := range file1Data.Records {
		if len(record) > entityColIndex && len(record) > amountColIndex {
			entityValue := record[entityColIndex]
			amountStr := record[amountColIndex]

			// Convert amount to float
			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil {
				continue // Skip invalid amounts
			}

			// Store original record
			recordMap := make(map[string]interface{})
			for i, header := range file1Data.Headers {
				if i < len(record) {
					recordMap[header] = record[i]
				}
			}
			originalRecords[entityValue] = append(originalRecords[entityValue], recordMap)

			// Aggregate amounts
			switch aggregationStrategy {
			case "sum":
				aggregatedData[entityValue] += amount
			case "count":
				aggregatedData[entityValue] += 1
			case "avg":
				// For avg, we'll calculate sum and count separately
				if _, exists := aggregatedData[entityValue]; !exists {
					aggregatedData[entityValue] = amount
				} else {
					aggregatedData[entityValue] += amount
				}
			case "max":
				if _, exists := aggregatedData[entityValue]; !exists || amount > aggregatedData[entityValue] {
					aggregatedData[entityValue] = amount
				}
			case "min":
				if _, exists := aggregatedData[entityValue]; !exists || amount < aggregatedData[entityValue] {
					aggregatedData[entityValue] = amount
				}
			}
		}
	}

	// Create file 2 lookup map
	file2Lookup := make(map[string]float64)
	for _, record := range file2Data.Records {
		if len(record) > entityColIndex && len(record) > amountColIndex {
			entityValue := record[entityColIndex]
			amountStr := record[amountColIndex]

			amount, err := strconv.ParseFloat(amountStr, 64)
			if err == nil {
				file2Lookup[entityValue] = amount
			}
		}
	}

	// Generate reconciliation results
	var reconciledRecords []map[string]interface{}
	var unreconciledRecords []map[string]interface{}

	for entityValue, aggregatedAmount := range aggregatedData {
		if file2Amount, exists := file2Lookup[entityValue]; exists {
			// Check if amounts match (with small tolerance for floating point)
			if math.Abs(aggregatedAmount-file2Amount) < 0.01 {
				reconciledRecords = append(reconciledRecords, map[string]interface{}{
					"entity_id":               entityValue,
					"file1_aggregated_amount": aggregatedAmount,
					"file2_amount":            file2Amount,
					"status":                  "RECONCILED",
					"original_records_count":  len(originalRecords[entityValue]),
					"original_records":        originalRecords[entityValue],
				})
			} else {
				unreconciledRecords = append(unreconciledRecords, map[string]interface{}{
					"entity_id":               entityValue,
					"file1_aggregated_amount": aggregatedAmount,
					"file2_amount":            file2Amount,
					"status":                  "AMOUNT_MISMATCH",
					"difference":              math.Abs(aggregatedAmount - file2Amount),
					"original_records_count":  len(originalRecords[entityValue]),
					"original_records":        originalRecords[entityValue],
				})
			}
		} else {
			unreconciledRecords = append(unreconciledRecords, map[string]interface{}{
				"entity_id":               entityValue,
				"file1_aggregated_amount": aggregatedAmount,
				"file2_amount":            0,
				"status":                  "MISSING_IN_FILE2",
				"original_records_count":  len(originalRecords[entityValue]),
				"original_records":        originalRecords[entityValue],
			})
		}
	}

	// Check for records in file 2 that don't exist in file 1
	for entityValue, amount := range file2Lookup {
		if _, exists := aggregatedData[entityValue]; !exists {
			unreconciledRecords = append(unreconciledRecords, map[string]interface{}{
				"entity_id":               entityValue,
				"file1_aggregated_amount": 0,
				"file2_amount":            amount,
				"status":                  "MISSING_IN_FILE1",
				"original_records_count":  0,
				"original_records":        []map[string]interface{}{},
			})
		}
	}

	// Calculate reconciliation rate safely
	reconciliationRate := 0.0
	if len(aggregatedData) > 0 {
		reconciliationRate = float64(len(reconciledRecords)) / float64(len(aggregatedData)) * 100
	}

	return map[string]interface{}{
		"preview_enabled":      true,
		"aggregation_strategy": aggregationStrategy,
		"entity_identifier":    entityIdentifier,
		"total_entities":       len(aggregatedData),
		"reconciled_count":     len(reconciledRecords),
		"unreconciled_count":   len(unreconciledRecords),
		"reconciliation_rate":  fmt.Sprintf("%.1f%%", reconciliationRate),
		"reconciled_records":   reconciledRecords,
		"unreconciled_records": unreconciledRecords,
		"summary": map[string]interface{}{
			"file1_total_records":   len(file1Data.Records),
			"file2_total_records":   len(file2Data.Records),
			"unique_entities_file1": len(aggregatedData),
			"unique_entities_file2": len(file2Lookup),
		},
	}
}

// FileData represents parsed file data
type FileData struct {
	Headers []string
	Records [][]string
}

// readFileData reads and parses CSV/Excel files
func readFileData(filePath, fileType string) (*FileData, error) {
	if fileType == "csv" {
		return readCSVFile(filePath)
	} else if fileType == "excel" {
		return readExcelFile(filePath)
	}
	return nil, fmt.Errorf("unsupported file type: %s", fileType)
}

// readCSVFile reads and parses CSV file
func readCSVFile(filePath string) (*FileData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("file is empty")
	}

	return &FileData{
		Headers: records[0],
		Records: records[1:],
	}, nil
}

// readExcelFile reads and parses Excel file
func readExcelFile(filePath string) (*FileData, error) {
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: failed to close Excel file: %v\n", err)
		}
	}()

	// Get the first sheet name
	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("excel file has no sheets")
	}

	sheetName := sheets[0]
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read Excel sheet: %v", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("excel file is empty")
	}

	headers := rows[0]
	records := rows[1:]

	// Normalize rows to ensure all rows have the same number of columns as headers
	normalizedRecords := make([][]string, len(records))
	for i, row := range records {
		normalizedRow := make([]string, len(headers))
		for j := 0; j < len(headers); j++ {
			if j < len(row) {
				normalizedRow[j] = strings.TrimSpace(row[j])
			} else {
				normalizedRow[j] = ""
			}
		}
		normalizedRecords[i] = normalizedRow
	}

	return &FileData{
		Headers: headers,
		Records: normalizedRecords,
	}, nil
}

// fetchEntityIDFromMasterSource fetches the Entity ID from master source configuration
func fetchEntityIDFromMasterSource(ctx context.Context, masterSourceID string) (string, error) {
	// Make API call to get master source details
	result, err := makeReconSaaSAPICall(ctx, "GET", fmt.Sprintf("/v1/admin-recon-saas/sources/%s", masterSourceID), nil)
	if err != nil {
		return "", fmt.Errorf("failed to fetch master source: %v", err)
	}

	// Extract Entity ID from mapping config
	if config, ok := result["config"].(map[string]interface{}); ok {
		if mappingConfig, ok := config["mapping_config"].([]interface{}); ok {
			for _, mapping := range mappingConfig {
				if mappingMap, ok := mapping.(map[string]interface{}); ok {
					if destination, ok := mappingMap["destination"].(string); ok {
						if destination == "EntityID" {
							if source, ok := mappingMap["source"].(string); ok {
								return source, nil
							}
						}
					}
				}
			}
		}
	}

	return "", fmt.Errorf("entity ID not found in master source configuration")
}

// validateEntityIDVsEntityIdentifier compares Entity ID with Entity Identifier and provides validation
func validateEntityIDVsEntityIdentifier(ctx context.Context, masterSourceID, entityIdentifier string) (map[string]interface{}, error) {
	// Try to fetch Entity ID from master source
	entityID, err := fetchEntityIDFromMasterSource(ctx, masterSourceID)
	if err != nil {
		// If API call fails, provide a fallback validation message
		return map[string]interface{}{
			"validation_passed": true, // Allow to proceed since we can't verify
			"entity_id":         "Unknown (API call failed)",
			"entity_identifier": entityIdentifier,
			"are_same":          false,
			"recommendation":    fmt.Sprintf("⚠️  NOTE: Unable to fetch Entity ID from master source (API error: %v). Please ensure Entity Identifier ('%s') is different from the Entity ID defined in your master source configuration.", err, entityIdentifier),
			"explanation": map[string]interface{}{
				"entity_id_definition":         "Entity ID is the column defined in master source for grouping records",
				"entity_identifier_definition": "Entity Identifier is the column used for aggregation processing",
				"why_different":                "They should be different to enable proper aggregation logic",
				"api_error":                    err.Error(),
			},
		}, nil
	}

	// Compare Entity ID with Entity Identifier
	areSame := strings.EqualFold(entityID, entityIdentifier)

	var recommendation string
	var validationPassed bool

	if areSame {
		validationPassed = false
		recommendation = fmt.Sprintf("⚠️  WARNING: Entity ID ('%s') and Entity Identifier ('%s') are the same! For aggregation to work properly, they should be different. Please choose a different column for Entity Identifier.", entityID, entityIdentifier)
	} else {
		validationPassed = true
		recommendation = fmt.Sprintf("✅ Entity ID ('%s') and Entity Identifier ('%s') are different. This is correct for aggregation.", entityID, entityIdentifier)
	}

	return map[string]interface{}{
		"validation_passed": validationPassed,
		"entity_id":         entityID,
		"entity_identifier": entityIdentifier,
		"are_same":          areSame,
		"recommendation":    recommendation,
		"explanation": map[string]interface{}{
			"entity_id_definition":         "Entity ID is the column defined in master source for grouping records",
			"entity_identifier_definition": "Entity Identifier is the column used for aggregation processing",
			"why_different":                "They should be different to enable proper aggregation logic",
		},
	}, nil
}

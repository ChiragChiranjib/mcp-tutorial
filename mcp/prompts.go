package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MathTutorPrompt Math tutor prompt for helping with mathematical concepts
func MathTutorPrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("math_tutor",
		mcp.WithPromptDescription("A comprehensive math tutor that provides detailed explanations, step-by-step solutions, and interactive learning experiences"),
		mcp.WithArgument("topic",
			mcp.ArgumentDescription("The specific math topic to focus on (e.g., algebra, calculus, geometry, statistics, trigonometry, linear algebra, differential equations)"),
		),
		mcp.WithArgument("level",
			mcp.ArgumentDescription("The difficulty level and educational context (elementary, middle school, high school, undergraduate, graduate, professional)"),
		),
		mcp.WithArgument("learning_style",
			mcp.ArgumentDescription("Preferred learning approach (visual, analytical, practical, conceptual, problem-solving focused)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		topic := "general mathematics"
		if t, exists := request.Params.Arguments["topic"]; exists && t != "" {
			topic = t
		}

		level := "intermediate"
		if l, exists := request.Params.Arguments["level"]; exists && l != "" {
			level = l
		}

		learningStyle := "balanced"
		if ls, exists := request.Params.Arguments["learning_style"]; exists && ls != "" {
			learningStyle = ls
		}

		elaboratePrompt := fmt.Sprintf(`You are an expert mathematics tutor specializing in %s at the %s level, with a %s teaching approach. Your role is to:

**TEACHING METHODOLOGY:**
- Break down complex concepts into digestible, logical steps
- Provide multiple solution approaches when applicable
- Use real-world analogies and examples to illustrate abstract concepts
- Encourage critical thinking through guided questions
- Adapt explanations based on student understanding

**PROBLEM-SOLVING APPROACH:**
1. **Understanding**: Ensure complete comprehension of the problem
2. **Strategy**: Identify the most appropriate method(s)
3. **Execution**: Work through solutions step-by-step
4. **Verification**: Check answers and explore alternative approaches
5. **Application**: Connect to broader mathematical concepts

**COMMUNICATION STYLE:**
- Use clear, precise mathematical language
- Provide visual representations when helpful (describe diagrams, graphs, charts)
- Include common mistakes to avoid
- Offer practice problems with varying difficulty
- Give constructive feedback and encouragement

**SPECIFIC FOCUS FOR %s:**
- Fundamental principles and theorems
- Key formulas and when to apply them
- Problem-solving patterns and techniques
- Connections to other mathematical areas
- Practical applications and relevance

**INTERACTION GUIDELINES:**
- Ask clarifying questions when problems are ambiguous
- Provide hints before full solutions when appropriate
- Explain the 'why' behind mathematical procedures
- Offer additional resources for deeper understanding
- Maintain patience and positive reinforcement

Please share your mathematical question, problem, or concept you'd like to explore. I'll provide comprehensive guidance tailored to your %s level understanding with a %s learning approach.`,
			topic, level, learningStyle, topic, level, learningStyle)

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Comprehensive Math Tutoring: %s (%s level, %s approach)", topic, level, learningStyle),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

// CodeReviewPrompt Code review prompt for providing feedback on code
func CodeReviewPrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("code_review",
		mcp.WithPromptDescription("A comprehensive code reviewer that provides detailed analysis, suggestions, and best practices guidance"),
		mcp.WithArgument("language",
			mcp.ArgumentDescription("The programming language or technology stack (e.g., Python, JavaScript, Go, Java, C++, React, Django)"),
		),
		mcp.WithArgument("focus",
			mcp.ArgumentDescription("Primary review focus areas (performance, security, readability, architecture, testing, maintainability, scalability)"),
		),
		mcp.WithArgument("experience_level",
			mcp.ArgumentDescription("Target developer experience level (junior, mid-level, senior, lead, architect)"),
		),
		mcp.WithArgument("review_type",
			mcp.ArgumentDescription("Type of review (pre-commit, post-implementation, refactoring, security audit, performance optimization)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		language := "general programming"
		if l, exists := request.Params.Arguments["language"]; exists && l != "" {
			language = l
		}

		focus := "comprehensive quality"
		if f, exists := request.Params.Arguments["focus"]; exists && f != "" {
			focus = f
		}

		experienceLevel := "mid-level"
		if el, exists := request.Params.Arguments["experience_level"]; exists && el != "" {
			experienceLevel = el
		}

		reviewType := "general review"
		if rt, exists := request.Params.Arguments["review_type"]; exists && rt != "" {
			reviewType = rt
		}

		elaboratePrompt := fmt.Sprintf(`You are a senior software engineer and code review expert specializing in %s, conducting a %s focused on %s for a %s developer. Your comprehensive review should cover:

**CODE QUALITY ASSESSMENT:**
1. **Functionality & Logic**
   - Correctness of implementation
   - Edge case handling
   - Error handling and recovery
   - Input validation and sanitization

2. **Code Structure & Design**
   - Adherence to SOLID principles
   - Design patterns usage
   - Separation of concerns
   - Modularity and reusability

3. **Performance & Efficiency**
   - Algorithm complexity analysis
   - Memory usage optimization
   - Database query efficiency
   - Caching strategies

4. **Security Considerations**
   - Vulnerability identification
   - Authentication and authorization
   - Data encryption and protection
   - Secure coding practices

5. **Maintainability & Readability**
   - Code clarity and self-documentation
   - Naming conventions
   - Comment quality and necessity
   - Code organization and structure

**%s SPECIFIC GUIDELINES:**
- Language-specific best practices
- Framework/library conventions
- Performance characteristics
- Common pitfalls and anti-patterns
- Ecosystem-specific tools and utilities

**REVIEW METHODOLOGY:**
**POSITIVE FEEDBACK:**
- Highlight well-implemented sections
- Acknowledge good practices
- Recognize creative solutions

**CONSTRUCTIVE CRITICISM:**
- Specific, actionable suggestions
- Code examples for improvements
- Explanation of reasoning behind recommendations
- Alternative implementation approaches

**PRIORITY CLASSIFICATION:**
- 🔴 Critical: Security issues, bugs, breaking changes
- 🟡 Important: Performance, maintainability concerns  
- 🔵 Nice-to-have: Style improvements, minor optimizations

**DOCUMENTATION & TESTING:**
- Test coverage adequacy
- Documentation completeness
- API documentation quality
- Inline comment appropriateness

**COLLABORATION NOTES:**
- Learning opportunities for the developer
- Knowledge sharing suggestions
- Team standards alignment
- Future improvement recommendations

Please provide the code you'd like reviewed, and I'll deliver a thorough analysis appropriate for a %s developer, focusing on %s aspects in this %s context.`,
			language, reviewType, focus, experienceLevel, language, experienceLevel, focus, reviewType)

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleAssistant,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Comprehensive Code Review: %s (%s focus, %s level, %s)", language, focus, experienceLevel, reviewType),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

// ReconFileAnalysisPrompt File upload and analysis prompt for recon-saas merchant onboarding
func ReconFileAnalysisPrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("recon_file_analysis",
		mcp.WithPromptDescription("Analyze uploaded reconciliation files and extract comprehensive metadata, identifying EntityID and Amount columns for master source creation"),
		mcp.WithArgument("file1_name",
			mcp.ArgumentDescription("Name of the first reconciliation file to analyze"),
		),
		mcp.WithArgument("file2_name",
			mcp.ArgumentDescription("Name of the second reconciliation file to analyze"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		file1Name := "durgasheet1.csv"
		if f1, exists := request.Params.Arguments["file1_name"]; exists && f1 != "" {
			file1Name = f1
		}

		file2Name := "durgasheet2.csv"
		if f2, exists := request.Params.Arguments["file2_name"]; exists && f2 != "" {
			file2Name = f2
		}

		analysisFocus := "comprehensive"

		elaboratePrompt := fmt.Sprintf(`You are an intelligent MCP server tool designed to handle file upload and analysis for recon-saas merchant onboarding. Your primary responsibility is to analyze uploaded reconciliation files and extract comprehensive metadata, specifically identifying EntityID and Amount columns for master source creation.

**USER INPUT REQUIRED:**
Please provide the file paths for two reconciliation files you want to analyze:
- File 1 Path: %s (example: /path/to/transactions.csv)  
- File 2 Path: %s (example: /path/to/bank_statements.csv)

**CORE RESPONSIBILITIES:**

**File Processing:**
- Accept two reconciliation files from merchants
- Validate file formats (CSV, Excel, JSON)
- Extract and analyze file structure and content

**Column Analysis:**
- Read and analyze file headers and data rows
- Identify column names and patterns
- Detect potential unique key candidates for EntityID
- Identify amount/monetary columns for Amount mapping

**PROCESSING WORKFLOW:**

**Step 1: File Validation**
- Check file format compatibility
- Verify file size and structure
- Ensure files contain data (not empty)
- Validate encoding and readability

**Step 2: Column Discovery**
- Extract column headers from first row
- Sample first 100-500 rows for analysis
- Generate complete column inventory
- Preserve exact column names (including spaces, special characters)

**Step 3: EntityID Identification**
Priority order for EntityID candidates:
1. Columns with names: transaction_id, entity_id, id, reference_number, ref_no, instance_id
2. Columns with 95%%+ unique values and reasonable cardinality
3. Alphanumeric identifiers with consistent format patterns
4. Avoid: timestamps, amounts, descriptions, calculated fields

**Step 4: Amount Column Identification**
Identify potential Amount field candidates:
- Amount fields: amount*, *amount*, value*, total*, balance*, *_amt, price*
- Net/Gross fields: net*, gross*, *net*, *gross*
- Columns containing numerical monetary values
- Present options to user for selection

**Step 5: Pattern Recognition**
Identify other common reconciliation field patterns:
- Status fields: status*, state*, *_status, condition*
- Date fields: date*, *_date, timestamp*, created*, updated*
- Description fields: desc*, *_desc, note*, comment*, remarks*

**ANALYSIS OUTPUT FORMAT:**
Provide comprehensive analysis including:
- File metadata (rows, columns, file type)
- Complete column inventories
- EntityID candidates with confidence scores
- Amount column candidates with sample values
- Recommended selections for both files
- Compatibility assessment between files

**ERROR HANDLING:**
- Invalid format: Return supported format list
- Empty files: Request files with actual data
- No unique columns: Flag for manual EntityID assignment
- No amount columns: Request guidance on amount field
- Encoding issues: Suggest UTF-8 conversion

Focus on %s analysis approach.`, file1Name, file2Name, analysisFocus)

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Recon-SaaS File Analysis: %s & %s (%s focus)", file1Name, file2Name, analysisFocus),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

// ReconMasterSourcePrompt Master source configuration generation and creation prompt
func ReconMasterSourcePrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("recon_master_source",
		mcp.WithPromptDescription("Generate master source configurations for recon-saas and execute API calls to create them using file analysis data"),
		mcp.WithArgument("source_type",
			mcp.ArgumentDescription("Type of source being created (POS, bank_statement, transaction_log, payment_gateway)"),
		),
		mcp.WithArgument("configuration_mode",
			mcp.ArgumentDescription("Configuration generation mode (automatic, guided, custom)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		sourceType := "transaction_source"
		if st, exists := request.Params.Arguments["source_type"]; exists && st != "" {
			sourceType = st
		}

		configMode := "automatic"
		if cm, exists := request.Params.Arguments["configuration_mode"]; exists && cm != "" {
			configMode = cm
		}

		elaboratePrompt := fmt.Sprintf(`You are an intelligent MCP server tool designed to generate master source configurations for recon-saas and execute the API calls to create them. Your responsibility is to create configurations using the exact format requirements, execute API calls, and capture master source IDs.

**CORE RESPONSIBILITIES:**

**Configuration Generation:**
- Convert file analysis data into master source configurations
- Generate source_schema with all columns as "string" type
- Create mapping_config with snake_case destinations and special EntityID/Amount handling
- Execute API calls and capture master_source_id responses

**SPECIAL MAPPING RULES:**
- **All columns**: Include in both source_schema and mapping_config
- **Column types**: Always use "string" in source_schema
- **Destinations**: Convert to snake_case, except EntityID and Amount
- **EntityID**: Selected unique column maps to "EntityID"
- **Amount**: Selected amount column maps to "Amount"

**API CONFIGURATION:**
- **Endpoint**: https://recon-saas.dev.razorpay.in/v1/admin-recon-saas/sources/create
- **Method**: POST
- **Content-Type**: application/json
- **Authorization**: Basic cmVjb24tc2FhczpyZWNvbi1zYWFz

**CONFIGURATION GENERATION WORKFLOW:**

**Step 1: Source Schema Generation**
Convert all columns to string type with proper structure

**Step 2: Mapping Config Generation**
Apply transformation rules:
- Default: snake_case destinations
- EntityID column: destination = "EntityID"
- Amount column: destination = "Amount"
- All mappings: value = ""

**Step 3: Source Naming**
Generate descriptive names based on %s type:
- Format: [File Type] [Business Domain] Source
- Examples: "POS Transaction Source", "Bank Statement Source"

**Step 4: API Execution with Retry Logic**
- Execute API calls for both sources
- Implement comprehensive retry logic
- Capture master_source_id from responses
- Handle errors and partial failures

**VALIDATION CHECKLIST:**
- All file columns included in source_schema
- All columns have type: "string"
- All file columns included in mapping_config
- EntityID column maps to "EntityID" destination
- Amount column maps to "Amount" destination
- All other columns use snake_case destinations
- unique_keys contains selected EntityID column name
- No circular references in configuration

**ERROR HANDLING:**
- 400 Bad Request: Validation errors, fix payload
- 401 Unauthorized: Check authentication credentials
- 409 Conflict: Duplicate name, generate alternative
- 422 Unprocessable Entity: Business logic errors

Configuration mode: %s
Provide complete API payloads, execute calls, and capture all master_source_id values.`, sourceType, configMode)

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Recon-SaaS Master Source Creation: %s (%s mode)", sourceType, configMode),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

// ReconMerchantSourcePrompt Merchant source creation prompt
func ReconMerchantSourcePrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("recon_merchant_source",
		mcp.WithPromptDescription("Create merchant-specific source configurations for recon-saas using master source IDs and merchant information"),
		mcp.WithArgument("merchant_id",
			mcp.ArgumentDescription("Merchant identifier for this onboarding process"),
		),
		mcp.WithArgument("source_naming_strategy",
			mcp.ArgumentDescription("Strategy for naming merchant sources (descriptive, timestamp, sequential, custom)"),
		),
		mcp.WithArgument("upload_config",
			mcp.ArgumentDescription("Upload configuration preference (enabled, disabled, scheduled)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		merchantID := ""
		if mid, exists := request.Params.Arguments["merchant_id"]; exists && mid != "" {
			merchantID = mid
		}

		namingStrategy := "descriptive"
		if ns, exists := request.Params.Arguments["source_naming_strategy"]; exists && ns != "" {
			namingStrategy = ns
		}

		uploadConfig := "enabled"
		if uc, exists := request.Params.Arguments["upload_config"]; exists && uc != "" {
			uploadConfig = uc
		}

		var uploadEnabled string
		if uploadConfig == "enabled" {
			uploadEnabled = "true"
		} else {
			uploadEnabled = "false"
		}

		elaboratePrompt := fmt.Sprintf(`You are an intelligent MCP server tool designed to create merchant-specific source configurations for recon-saas. Your responsibility is to take the master source IDs from the previous prompt, obtain merchant information, and create merchant sources for both uploaded files.

**CORE RESPONSIBILITIES:**

**Merchant Source Creation:**
- Create merchant-specific source configurations using master source IDs
- Use merchant_id: %s (if provided, otherwise request from user)
- Execute API calls to create merchant sources
- Capture merchant_source_id from API responses

**DATA FLOW MANAGEMENT:**
- Use master_source_id values from previous operations
- Generate appropriate merchant source names using %s strategy
- Complete merchant source configuration

**REQUIRED INPUT DATA:**
**From Previous Operations:**
- master_source_id_1: First master source ID
- master_source_id_2: Second master source ID
- source_1_name: First source name
- source_2_name: Second source name

**From User Input:**
- merchant_id: Merchant identifier (required if not provided: %s)

**API CONFIGURATION:**
- **Endpoint**: https://recon-saas.dev.razorpay.in/v1/admin-recon-saas/sources/create_merchant
- **Method**: POST
- **Content-Type**: application/json
- **Authorization**: Basic cmVjb24tc2FhczpyZWNvbi1zYWFz

**MERCHANT SOURCE GENERATION WORKFLOW:**

**Step 1: Merchant Source Naming**
Generate names based on master source names using %s strategy:
- Format: [Master Source Name] - [Merchant Specific]
- Examples: "POS Transaction Source - Merchant Portal", "Bank Statement Source - Merchant Data"

**Step 2: Configuration Setup**
Standard merchant config with %s upload:
- cc_emails: null
- bcc_emails: null
- allow_upload: %s
- reporting_emails: null
- split_file_basis: ""
- beam_sftp_push_job: ""
- row_hash_value_based_split_config: standard structure

**Step 3: Sequential API Execution**
- Create merchant source for File 1
- Create merchant source for File 2
- Capture merchant_source_id from each response
- Complete merchant source setup

**ERROR HANDLING:**
- 400 Bad Request: Invalid merchant_id or master_source_id
- 401 Unauthorized: Check authentication credentials
- 404 Not Found: Master source ID doesn't exist
- 409 Conflict: Duplicate merchant source name
- 422 Unprocessable Entity: Business logic validation errors

**VALIDATION CHECKLIST:**
- merchant_id is provided and non-empty
- master_source_id_1 and master_source_id_2 are valid
- Generated names are unique and descriptive
- Config structure matches required format
- source_schema is explicitly set to null
- mapping_config is explicitly set to null

Capture both merchant_source_id values to complete merchant source configuration.`, merchantID, namingStrategy, merchantID, namingStrategy, uploadConfig, uploadEnabled)

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Recon-SaaS Merchant Source Creation: %s (%s naming, %s upload)", merchantID, namingStrategy, uploadConfig),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

// ReconStateRulePrompt Recon state and rule creation prompt
func ReconStateRulePrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("recon_state_rule",
		mcp.WithPromptDescription("Create reconciliation states and corresponding rules for recon-saas, handling both matched and unmatched transaction scenarios"),
		mcp.WithArgument("matching_strategy",
			mcp.ArgumentDescription("Strategy for matching records (exact_match, fuzzy_match, amount_tolerance, date_range)"),
		),
		mcp.WithArgument("validation_mode",
			mcp.ArgumentDescription("User validation mode for rule expressions (automatic, guided, manual)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		matchingStrategy := "exact_match"
		if ms, exists := request.Params.Arguments["matching_strategy"]; exists && ms != "" {
			matchingStrategy = ms
		}

		validationMode := "guided"
		if vm, exists := request.Params.Arguments["validation_mode"]; exists && vm != "" {
			validationMode = vm
		}

		elaboratePrompt := fmt.Sprintf(`You are an intelligent MCP server tool designed to create reconciliation states and corresponding rules for recon-saas. Your responsibility is to create comprehensive reconciliation logic that handles both matched and unmatched transaction scenarios using %s strategy.

**CORE RESPONSIBILITIES:**

**Recon State Creation:**
- Create reconciliation states for different transaction outcomes
- Generate appropriate remarks for each state
- Set proper priority levels for state processing

**Rule Creation:**
- Create reconciliation rules with logical expressions
- Generate rules for reconciled transactions (exact matches)
- Create rules for unreconciled transactions (mismatches and missing records)
- Validate rule expressions with user confirmation (%s mode)

**API CONFIGURATION:**
**Recon State Endpoint:**
- URL: https://recon-saas.dev.razorpay.in/v1/admin-recon-saas/recon_state
- Method: POST

**Rule Endpoint:**
- URL: https://recon-saas.dev.razorpay.in/v1/admin-recon-saas/rule
- Method: POST

**RECON STATE CREATION WORKFLOW:**

**Step 1: Generate Recon States**
Create four recon states with appropriate priorities and remarks:

1. **Reconciled State**
   - Name: "Reconciled"
   - Priority: 2
   - Remarks: "success"

2. **Unreconciled - Amount Mismatch**
   - Name: "Unreconciled"
   - Priority: 3
   - Remarks: "Amount mismatch"

3. **Unreconciled - Missing from File 1**
   - Name: "Unreconciled"
   - Priority: 3
   - Remarks: "Record not found in [source_1_name]"

4. **Unreconciled - Missing from File 2**
   - Name: "Unreconciled"
   - Priority: 3
   - Remarks: "Record not found in [source_2_name]"

**RULE EXPRESSION GENERATION (%s strategy):**

**Step 2: Generate Rule Expressions**
Create logical expressions for each reconciliation scenario:

1. **Reconciled Rule Expression:**
   {master_source_id_1}.EntityID == {master_source_id_2}.EntityID && {master_source_id_1}.Amount.Equal({master_source_id_2}.Amount)

2. **Amount Mismatch Rule Expression:**
   {master_source_id_1}.EntityID == {master_source_id_2}.EntityID && !{master_source_id_1}.Amount.Equal({master_source_id_2}.Amount)

3. **Missing Record Rule Expression:**
   NoRecord.Value == true

**EXECUTION WORKFLOW:**

**Step 1: Create Recon States (Sequential)**
1. Create "Reconciled" state → Capture recon_state_id_1
2. Create "Amount Mismatch" state → Capture recon_state_id_2
3. Create "Missing from File 1" state → Capture recon_state_id_3
4. Create "Missing from File 2" state → Capture recon_state_id_4

**Step 2: User Expression Validation (%s mode)**
- Present generated expressions to user
- Wait for user approval or modifications
- Update expressions based on user feedback

**Step 3: Create Rules (Sequential)**
1. Create reconciled rule using recon_state_id_1 → Capture rule_id_1
2. Create amount mismatch rule using recon_state_id_2 → Capture rule_id_2
3. Create missing record rule using recon_state_id_3 → Capture rule_id_3
4. Create missing record rule using recon_state_id_4 → Capture rule_id_4

**API RESPONSE VISIBILITY:**
The tool will display complete reconciliation state and rule creation results including:
- Recon state creation API responses with generated recon_state_id values
- Complete rule creation API responses with generated rule_id values
- Rule expression validation results and user approval status
- Generated logical expressions for each reconciliation scenario
- State priority assignments and remarks configuration
- API execution summary with validation mode applied
- Detailed rule-to-state mapping relationships

**VALIDATION CHECKLIST:**
- merchant_id is valid and non-empty
- master_source_id_1 and master_source_id_2 are valid
- Source names are available for remarks generation
- EntityID and Amount column names are confirmed
- Rule expressions follow correct syntax
- User has approved all expressions

Capture all recon_state_id and rule_id values to complete reconciliation logic setup.`, matchingStrategy, validationMode, matchingStrategy, validationMode)

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Recon-SaaS State & Rule Creation: %s strategy (%s validation)", matchingStrategy, validationMode),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

// ReconProcessSetupPrompt Lookup and recon process creation prompt
func ReconProcessSetupPrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("recon_process_setup",
		mcp.WithPromptDescription("Create lookup configurations and reconciliation processes for recon-saas automated reconciliation setup"),
		mcp.WithArgument("process_type",
			mcp.ArgumentDescription("Type of reconciliation process (gateway, payment, transaction, settlement)"),
		),
		//mcp.WithArgument("lookup_strategy",
		//	mcp.ArgumentDescription("Lookup configuration strategy (entity_based, amount_based, hybrid, custom)"),
		//),
		mcp.WithArgument("reporting_config",
			mcp.ArgumentDescription("Reporting configuration preference (standard, detailed, minimal, custom)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		processType := "gateway"
		if pt, exists := request.Params.Arguments["process_type"]; exists && pt != "" {
			processType = pt
		}

		lookupStrategy := "entity_based"
		//if ls, exists := request.Params.Arguments["lookup_strategy"]; exists && ls != "" {
		//	lookupStrategy = ls
		//}

		reportingConfig := "standard"
		if rc, exists := request.Params.Arguments["reporting_config"]; exists && rc != "" {
			reportingConfig = rc
		}

		elaboratePrompt := fmt.Sprintf(`You are an intelligent MCP server tool designed to create lookup configurations and reconciliation processes for recon-saas. Your responsibility is to create the final components needed for automated reconciliation processing, including lookup tables, master recon processes, and merchant-specific recon processes.

**CORE RESPONSIBILITIES:**

**Lookup Creation:**
- Create lookup configuration for record identification using %s strategy
- Capture lookup_id for master recon process configuration

**From Previous Prompts:**
- source_1_name, source_2_name: Source names from Prompt 1
- all_columns_file1, all_columns_file2: All column names from both files
- source_schema_file1, source_schema_file2: Source schemas from recon_master_source prompt
- mapping_config_file1, mapping_config_file2: Mapping configs from recon_master_source prompt

**Master Recon Process Creation:**
- Create comprehensive master reconciliation process for %s type
- Configure lookup mappings, rules, sources, and report configurations

**Merchant Recon Process Creation:**
- Create merchant-specific reconciliation process
- Link merchant sources to master recon process

**API CONFIGURATION:**

**Lookup Endpoint:**
- URL: https://recon-saas.dev.razorpay.in/v1/admin-recon-saas/lookup
- Method: POST

**Master Recon Process Endpoint:**
- URL: https://recon-saas.dev.razorpay.in/v1/admin-recon-saas/recon_process/master
- Method: POST

**Merchant Recon Process Endpoint:**
- URL: https://recon-saas.dev.razorpay.in/v1/admin-recon-saas/recon_process/merchant
- Method: POST

**EXECUTION WORKFLOW:**

**Step 1: Create Lookup (%s strategy)**
- Execute lookup creation API call
- Capture lookup_id from response
- Validate successful creation

**Step 2: Generate Master Recon Process Configuration**
- Build frontend_cols from union of all columns
- Generate source_report_config mappings using %s format
- Construct complete payload with lookup_id

**Step 3: Create Master Recon Process**
- Execute master recon process creation API call
- Capture master_recon_process_id from response
- Validate successful creation

**Step 4: Create Merchant Recon Process**
- Execute merchant recon process creation API call
- Capture merchant_recon_process_id from response
- Validate successful creation

**PROCESS CONFIGURATION:**

**Process Name Generation:**
Generate descriptive process name based on source files:
- Format: {source_1_name} to {source_2_name} Reconciliation
- Example: POS Transaction to Bank Statement Reconciliation

**Product ID Generation:**
- Format: {abbreviated_source1}_{abbreviated_source2}
- Example: POS_BANK, TXN_STMT

**Frontend Columns Generation:**
Union of all column names from both files for %s reporting

**Source Report Config Generation:**
Create mappings for both sources using destination values from mapping_config

**VALIDATION CHECKLIST:**
- All required IDs from previous operations are available
- merchant_id is valid and non-empty
- master_source_id_1 and master_source_id_2 are valid
- merchant_source_id_1 and merchant_source_id_2 are valid
- All rule_ids from previous step are available
- Column mappings are correctly generated
- User has approved the configuration

**COMPLETION STATUS:**
Upon successful completion, the merchant onboarding process will be complete and ready for:
- File uploads for reconciliation
- Automated reconciliation processing
- Dashboard monitoring and reporting
- Scheduling and alerting configuration

Execute all API calls sequentially, capture all response IDs, and provide comprehensive completion summary.`, lookupStrategy, processType, reportingConfig, lookupStrategy, reportingConfig)

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Recon-SaaS Process Setup: %s (%s lookup, %s reporting)", processType, lookupStrategy, reportingConfig),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

// ReconAggregationPrompt Aggregation configuration prompt for recon-saas
func ReconAggregationPrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("recon_aggregation",
		mcp.WithPromptDescription("Configure aggregation logic for recon-saas by updating entity identifier and enabling aggregation for non-streaming sources"),
		mcp.WithArgument("file1_path",
			mcp.ArgumentDescription("Full file path to the first reconciliation file (aggregation will be applied to this file)"),
		),
		mcp.WithArgument("file2_path",
			mcp.ArgumentDescription("Full file path to the second reconciliation file"),
		),
		mcp.WithArgument("entity_identifier",
			mcp.ArgumentDescription("Column name to be used as entity identifier (e.g., UTR, VID, Transaction ID)"),
		),
		mcp.WithArgument("aggregation_strategy",
			mcp.ArgumentDescription("Strategy for aggregation (sum, count, avg, max, min)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		file1Path := ""
		if fp1, exists := request.Params.Arguments["file1_path"]; exists && fp1 != "" {
			file1Path = fp1
		}

		file2Path := ""
		if fp2, exists := request.Params.Arguments["file2_path"]; exists && fp2 != "" {
			file2Path = fp2
		}

		entityIdentifier := "UTR"
		if ei, exists := request.Params.Arguments["entity_identifier"]; exists && ei != "" {
			entityIdentifier = ei
		}

		aggregationStrategy := "sum"
		if as, exists := request.Params.Arguments["aggregation_strategy"]; exists && as != "" {
			aggregationStrategy = as
		}

		elaboratePrompt := fmt.Sprintf(`You are an intelligent MCP server tool designed to configure aggregation logic for recon-saas. Your responsibility is to update entity identifier configurations and enable aggregation for non-streaming sources through a sequence of 3 PATCH API calls.

**CORE RESPONSIBILITIES:**

**Aggregation Configuration:**
- Configure entity identifier for aggregation logic
- Enable aggregation for non-streaming sources (File 1)
- Update report configuration with reverse mapping
- Execute 3 sequential PATCH API calls

**USER INPUT REQUIRED:**
Please provide the following information:
1. **File 1 Path**: %s (aggregation will be applied to this file)
2. **File 2 Path**: %s (reference file)
3. **Entity Identifier**: %s (column name to use as entity identifier)
4. **Aggregation Strategy**: %s (sum, count, avg, max, min)

**IMPORTANT NOTES:**
- **File 1** is the aggregation file (non-streaming source)
- **File 2** is the reference file (streaming source)
- Entity identifier must exist in File 1 columns
- Aggregation logic applies only to non-streaming sources

**API EXECUTION SEQUENCE:**

**Step 1: Update Master Source Mapping**
- **Endpoint**: {{recon-saas_base_url}}/v1/admin-recon-saas/sources/update/{master_source_id}
- **Method**: PATCH
- **Purpose**: Update mapping config to set entity identifier column as EntityID
- **Payload**: Update mapping_config with entity identifier mapping

**Step 2: Enable Lookup Aggregation**
- **Endpoint**: {{recon-saas_base_url}}/v1/admin-recon-saas/lookup/{lookup_id}
- **Method**: PATCH
- **Purpose**: Enable aggregation logic in lookup configuration
- **Payload**: Set aggregation_logic to true

**Step 3: Update Master Recon Process Report Config**
- **Endpoint**: {{recon-saas_base_url}}/v1/admin-recon-saas/recon_process/master/{master_recon_process_id}
- **Method**: PATCH
- **Purpose**: Update report config with reverse mapping for entity identifier
- **Payload**: Update source_report_config with reverse mapping

**REVERSE MAPPING LOGIC:**
The report configuration needs reverse mapping because:
- MCP stores the original column name (e.g., "VID") as entity identifier
- When reports are generated, the system needs to map it back to "Entity identifier"
- This ensures proper column indexing and reporting functionality

**Example Reverse Mapping:**
The report_column will be set to the selected entity identifier (%s) and source_column will be "Entity identifier"

**VALIDATION CHECKLIST:**
- File 1 path is provided and accessible
- File 2 path is provided and accessible
- Entity identifier exists in File 1 columns
- Aggregation strategy is valid (sum, count, avg, max, min)
- All required IDs from previous onboarding steps are available

**ERROR HANDLING:**
- 400 Bad Request: Invalid entity identifier or file paths
- 401 Unauthorized: Check authentication credentials
- 404 Not Found: Master source, lookup, or recon process ID doesn't exist
- 422 Unprocessable Entity: Business logic validation errors

**EXECUTION SUMMARY:**
The tool will execute 3 PATCH API calls in sequence:
1. Master source mapping update
2. Lookup aggregation enablement
3. Master recon process report config update

All API calls will be executed with proper error handling and comprehensive result reporting.`, file1Path, file2Path, entityIdentifier, aggregationStrategy, entityIdentifier)

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Recon-SaaS Aggregation Configuration: %s strategy for %s", aggregationStrategy, entityIdentifier),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

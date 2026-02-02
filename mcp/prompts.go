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
- **Endpoint**: %s/v1/admin-recon-saas/sources/create
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
Provide complete API payloads, execute calls, and capture all master_source_id values.`, sourceType, GetBaseURL(DefaultEnvironment), configMode)

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
- **Endpoint**: %s/v1/admin-recon-saas/sources/create_merchant
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

Capture both merchant_source_id values to complete merchant source configuration.`, merchantID, namingStrategy, merchantID, GetBaseURL(DefaultEnvironment), namingStrategy, uploadConfig, uploadEnabled)

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
- URL: %s/v1/admin-recon-saas/recon_state
- Method: POST

**Rule Endpoint:**
- URL: %s/v1/admin-recon-saas/rule
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

Capture all recon_state_id and rule_id values to complete reconciliation logic setup.`, matchingStrategy, GetBaseURL(DefaultEnvironment), GetBaseURL(DefaultEnvironment), validationMode, matchingStrategy, validationMode)

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
- URL: %s/v1/admin-recon-saas/lookup
- Method: POST

**Master Recon Process Endpoint:**
- URL: %s/v1/admin-recon-saas/recon_process/master
- Method: POST

**Merchant Recon Process Endpoint:**
- URL: %s/v1/admin-recon-saas/recon_process/merchant
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

Execute all API calls sequentially, capture all response IDs, and provide comprehensive completion summary.`, lookupStrategy, processType, GetBaseURL(DefaultEnvironment), GetBaseURL(DefaultEnvironment), GetBaseURL(DefaultEnvironment), lookupStrategy, reportingConfig, reportingConfig)

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

// ReconEntityUpdatePrompt Entity update prompt for modifying recon-saas configurations
func ReconEntityUpdatePrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("recon_entity_update",
		mcp.WithPromptDescription("Update existing recon-saas entities including master sources, merchant sources, recon processes, rules, recon states, and lookups"),
		mcp.WithArgument("entity_type",
			mcp.ArgumentDescription("Type of entity to update (master_source, merchant_source, master_recon_process, merchant_recon_process, rule, recon_state, lookup)"),
		),
		mcp.WithArgument("update_scope",
			mcp.ArgumentDescription("Scope of the update (single_field, multiple_fields, complete_reconfiguration)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		entityType := "general"
		if et, exists := request.Params.Arguments["entity_type"]; exists && et != "" {
			entityType = et
		}

		updateScope := "single_field"
		if us, exists := request.Params.Arguments["update_scope"]; exists && us != "" {
			updateScope = us
		}

		elaboratePrompt := fmt.Sprintf(`You are an intelligent MCP server tool designed to update existing recon-saas entities. Your responsibility is to make PATCH API calls to modify configurations for various entity types.

**CORE RESPONSIBILITIES:**

**Entity Update Operations:**
- Update existing entity configurations via PATCH API calls
- Validate update payloads before submission
- Handle partial updates (only modified fields)
- Return updated entity state after successful modification

**SUPPORTED ENTITY TYPES (%s focus):**

1. **master_source** - Master source configurations
   - Endpoint: /v1/admin-recon-saas/sources/update/{id}
   - Updatable fields: name, skip_top_rows, ingest_to_db, allow_upload, unique_keys, source_schema, mapping_config, transformation_config, validation_config, sub_source_config, extract_distinct_config, report_enrichment, split_file_basis, is_header_missing, metadata_extraction_config, skip_bottom_rows, skip_row_func

2. **merchant_source** - Merchant-specific source configurations
   - Endpoint: /v1/admin-recon-saas/sources/update_merchant/{id}
   - Updatable fields: name, master_source_id, reporting_emails, cc_emails, bcc_emails, allow_upload, source_schema, mapping_config, validation_config, split_file_basis, beam_sftp_push_job, slack_notification_config

3. **master_recon_process** - Master reconciliation process
   - Endpoint: /v1/admin-recon-saas/recon_process/master/{id}
   - Updatable fields: name, lookup_config, product_id, rules, sources, sequence, report_config, workflow_config

4. **merchant_recon_process** - Merchant reconciliation process
   - Endpoint: /v1/admin-recon-saas/recon_process/merchant/{id}
   - Updatable fields: sources, report_config, skip_status, skip_rows, skip_rows_recon_state_ids, report_channel, status

5. **rule** - Reconciliation rules
   - Endpoint: /v1/admin-recon-saas/rule/{id}
   - Updatable fields: name, type, expression, sources, recon_state_id

6. **recon_state** - Reconciliation states
   - Endpoint: /v1/admin-recon-saas/recon_state/{id}
   - Updatable fields: name, priority, remarks

7. **lookup** - Lookup configurations
   - Endpoint: /v1/admin-recon-saas/lookup/{id}
   - Updatable fields: name, config

**UPDATE WORKFLOW (%s scope):**

**Step 1: Identify Entity**
- Obtain entity_type and entity_id from user
- Validate entity exists (optional GET call)

**Step 2: Prepare Update Payload**
- Collect fields to update from user
- Validate field names match entity type
- Construct JSON payload with only changed fields

**Step 3: Execute PATCH Request**
- Make PATCH API call to appropriate endpoint
- Handle authentication (Basic auth)
- Process response and error handling

**Step 4: Verify Update**
- Display updated entity from response
- Confirm changes were applied correctly
- Report any validation or business logic errors

**COMMON UPDATE SCENARIOS:**

**Master Source Updates:**
- Change source name: {"name": "New Source Name"}
- Update schema: {"source_schema": [{"name": "col1", "type": "string"}]}
- Modify mappings: {"mapping_config": [{"source": "col1", "destination": "EntityID", "value": ""}]}

**Merchant Source Updates:**
- Update email recipients: {"reporting_emails": ["new@email.com"], "cc_emails": ["cc@email.com"]}
- Enable/disable uploads: {"allow_upload": true}
- Configure Slack alerts: {"slack_notification_config": {"recon_percentage_threshold": 90, "file_alert_enabled": true}}

**Rule Updates:**
- Modify expression: {"expression": "SourceA.EntityID == SourceB.EntityID"}
- Change associated state: {"recon_state_id": "new_state_id"}
- Update rule name: {"name": "Updated Rule Name"}

**Recon State Updates:**
- Change priority: {"priority": 1}
- Update remarks: {"remarks": "Updated description"}

**Process Updates:**
- Update sources: {"sources": ["source1", "source2"]}
- Change status: {"status": "approved"}
- Modify report config: {"report_config": {...}}

**API CONFIGURATION:**
- Base URL: %s
- Method: PATCH
- Content-Type: application/json
- Authorization: Basic cmVjb24tc2FhczpyZWNvbi1zYWFz

**ERROR HANDLING:**
- 400 Bad Request: Invalid payload structure or field values
- 401 Unauthorized: Authentication failure
- 404 Not Found: Entity ID does not exist
- 422 Unprocessable Entity: Business logic validation errors
- 500 Internal Server Error: Server-side issues

**VALIDATION CHECKLIST:**
- entity_id is valid and non-empty
- entity_type matches supported types
- update_payload contains valid field names
- Field values match expected types
- Required relationships exist (e.g., recon_state_id for rules)

Please provide the entity type, entity ID, and the fields you want to update.`, entityType, updateScope, GetBaseURL(DefaultEnvironment))

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Recon-SaaS Entity Update: %s (%s)", entityType, updateScope),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

// ReconTransformationConfigPrompt Transformation configuration prompt for applying data transformations
func ReconTransformationConfigPrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("recon_transformation_config",
		mcp.WithPromptDescription("Apply data transformations to recon-saas master source columns ONLY when explicitly requested by the user"),
		mcp.WithArgument("transformation_type",
			mcp.ArgumentDescription("Type of transformation needed (amount_parsing, column_concatenation, date_formatting, string_extraction, calculation)"),
		),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		transformationType := "general"
		if tt, exists := request.Params.Arguments["transformation_type"]; exists && tt != "" {
			transformationType = tt
		}

		elaboratePrompt := fmt.Sprintf(`You are an intelligent MCP server tool designed to apply data transformations to recon-saas master sources. Your responsibility is to understand the user's transformation needs and apply the correct transformation function.

**CRITICAL: DO NOT AUTO-APPLY TRANSFORMATIONS**

This tool should ONLY be used when the user EXPLICITLY requests a transformation. Do NOT automatically apply any transformation unless the user specifically asks for it.

**WHEN TO USE THIS TOOL:**
- User says: "Apply regex on column X to extract EntityID" → USE this tool
- User says: "Concatenate columns A, B, C to create EntityID for Source X" → USE this tool
- User says: "Transform the date format from X to Y" → USE this tool
- User says: "I need to apply a transformation on [column] with [master source]" → USE this tool

**WHEN NOT TO USE THIS TOOL:**
- During normal master source creation → DO NOT use
- During file analysis → DO NOT use
- During onboarding flow unless explicitly asked → DO NOT use
- When user doesn't mention transformation, regex, concatenation, or data formatting → DO NOT use
- When analyzing columns or suggesting mappings → DO NOT use (just suggest, don't apply)

**REQUIRED USER CONFIRMATION:**
Before applying ANY transformation, you MUST have explicit confirmation from the user about:
1. Which master source (by name or ID) to apply the transformation to
2. Which column(s) to transform
3. What transformation function to use
4. What the output column should be named

If any of these are unclear, ASK the user. Do NOT assume or auto-apply.

**SMART CONTEXT USAGE:**

Before asking questions, CHECK THE CONVERSATION CONTEXT for information from previous tool calls:
- If recon_master_source tool was used, master_source_id should be available from its response
- If recon_file_analysis tool was used, column names are available
- The tool can AUTO-FETCH current mapping_config and transformation_config via GET API using master_source_id

**WHAT YOU NEED TO CONFIRM WITH USER:**

Only ask for information that is NOT already available in the conversation:
1. Which source? (Source A, Source B, or the exact source name) - if not clear from context
2. What is the master_source_id? - if not available from previous tool calls
3. Which column(s) should be transformed? - ALWAYS confirm this
4. What should be the output column name? - ALWAYS confirm this
5. For specific functions, ask for required parameters:
   - regex_exec: What regex pattern?
   - change_date_format: What is the source format and target format?
   - excel_mid: What start position and length?
   - split: What delimiter and which part (index)?
   - etc.

**AUTO-FETCH CAPABILITY:**

If you have the master_source_id, the tool will automatically fetch:
- Current mapping_config (all existing column mappings)
- Current transformation_config (existing transformations)

This is done via GET API: /v1/admin-recon-saas/sources/get/{master_source_id}
So you do NOT need to ask user for current_mapping_config or current_transformation_config!

**TRANSFORMATION WORKFLOW (%s focus):**

Step 1: WAIT FOR EXPLICIT USER REQUEST
- DO NOT proactively suggest or apply transformations
- Only proceed if user explicitly asks for a transformation

Step 2: CHECK CONTEXT
- Look for master_source_id from previous recon_master_source tool calls
- Look for column names from previous recon_file_analysis tool calls
- Identify which source the user is referring to

Step 3: CONFIRM WITH USER (only what's missing)
- Input column(s) - which column(s) to transform
- Output column name - what to call the result
- Function-specific parameters (regex, date format, etc.)

Step 4: APPLY TRANSFORMATION (only after explicit confirmation)
- Call recon_transformation_config tool
- The tool will auto-fetch current configs if not provided
- All existing mappings are preserved automatically

**MAPPING CONFIG UPDATE LOGIC (Handled automatically by the tool):**

Special columns: EntityID, EntityStatus, EntityIdentifier, Amount

When output_column is a SPECIAL column (e.g., EntityID):
1. Find existing mapping where destination = special column
2. Change that mapping's destination to snake_case of its source
3. Append new mapping: {source: special_column, destination: special_column}

Example:
- Before: [{"source": "Notes", "destination": "EntityID"}, {"source": "amount", "destination": "Amount"}]
- Transformation output_column: "EntityID"
- After: [{"source": "Notes", "destination": "notes"}, {"source": "amount", "destination": "Amount"}, {"source": "EntityID", "destination": "EntityID"}]

When output_column is NOT a special column:
1. Append new mapping: {source: output_column, destination: snake_case(output_column)}

**AVAILABLE TRANSFORMATION FUNCTIONS:**

**Amount/Number:**
- abs_amount_parsing: Parse absolute amount (removes commas, handles negatives)
- abs_amount_in_paisa: Convert to paisa (multiply by 100)
- add_amount_cols: Sum multiple amount columns
- subtract_amount_cols: Subtract amounts from base
- percentage_of_number: Calculate percentage (params: [base_amount, percentage])
- extract_amount_from_cols: Get first non-zero from multiple columns

**String:**
- append_multiple_columns: Concatenate multiple columns
- excel_mid: Extract substring (params: [start_pos, length])
- excel_left: Extract from left (params: [length])
- excel_right: Extract from right (params: [length])
- split: Split and get part (params: [delimiter, index])
- regex_exec: Extract using regex (params: [regex_pattern])
- remove_prefix: Remove prefix (params: [prefix1, prefix2, ...])
- remove_suffix: Remove suffix (params: [suffix1, suffix2, ...])
- add_padding_prefix: Zero-pad (params: [target_length])
- replace_blank_string: Trim whitespace
- remove_single_quotes, remove_double_quotes: Remove quotes

**Date:**
- change_date_format: Convert formats (params: [from_format, to_format])
- date_normalization: Normalize to YYYY-MM-DD (params: [current_format])
- txn_date_extraction_generic: Parse various date formats
- excel_to_datetime: Convert Excel serial date
- subtract_date, add_date: Add/subtract days (params: [days])

**Other:**
- hard_code_value: Set constant (params: [value])
- settlement_amount_from_debit_credit_cols: Derive from debit/credit
- get_field_from_notes: Extract from JSON (params: [field1, field2, ...])

**EXAMPLE CONVERSATIONS:**

**Example 1: User explicitly requests transformation**

[Previous context: recon_master_source returned master_source_id: "S9W6Gpbw6OWbYo" for "Source A"]

User: "I want to apply regex on notes column in Source A to extract EntityID"

Assistant should ask:
"I have the master_source_id 'S9W6Gpbw6OWbYo' for Source A from our previous setup. 
What regex pattern should I use to extract the EntityID from the notes column?"

User: "[0-9]{7}"

Then call the tool (configs auto-fetched):
- master_source_id: "S9W6Gpbw6OWbYo"
- source_name: "Source A"
- transformation_function: "regex_exec"
- input_columns: ["notes"]
- additional_params: ["[0-9]{7}"]
- output_column: "EntityID"
(current_mapping_config and current_transformation_config will be auto-fetched via GET API)

**Example 2: User explicitly requests concatenation**

User: "For Source B, combine RRN, TID and MID to create EntityID"

Assistant should ask:
"I'll concatenate these columns. Just to confirm:
- Input columns: RRN, TID, MID
- Output column: EntityID
What is the master_source_id for Source B?" (if not in context)

**Example 3: User explicitly requests date format change**

User: "Change the date format of created_at column from MM/DD/YYYY to YYYY-MM-DD"

Assistant should ask:
"Which source should I apply this to? Also, just to confirm:
- Input column: created_at
- Current format: %%m/%%d/%%Y
- Target format: %%Y-%%m-%%d
- Output column: Should it update the same column (created_at) or create a new one?"

**Example 4: User does NOT request transformation - DO NOT APPLY**

User: "Create master sources for my two files"
→ DO NOT apply any transformation. Just create the master sources.

User: "Analyze my files and set up reconciliation"
→ DO NOT apply any transformation. Complete the setup without transformations.

**TOOL PARAMETERS:**

Required:
- master_source_id: ID of the master source
- source_name: Name for reference
- transformation_function: Function name from the list above
- input_columns: JSON array of input column names
- output_column: Name of the output column

Optional (auto-fetched if not provided):
- current_mapping_config: Current mapping config JSON array
- current_transformation_config: Current transformations, defaults to []
- additional_params: JSON array of extra parameters for the function
- environment: "local", "dev", or "prod" (defaults to "dev")

REMEMBER: Only apply transformations when the user EXPLICITLY requests them!`, transformationType)

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			fmt.Sprintf("Recon-SaaS Transformation Config: %s", transformationType),
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

// ReconAggregationConfigPrompt Aggregation configuration prompt for enabling aggregation on master sources
func ReconAggregationConfigPrompt() server.ServerPrompt {
	prompt := mcp.NewPrompt("recon_aggregation_config",
		mcp.WithPromptDescription("Configure aggregation for recon-saas master sources ONLY when explicitly requested by the user"),
	)

	handler := func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		elaboratePrompt := `You are an intelligent MCP server tool designed to configure aggregation for recon-saas master sources. Your responsibility is to enable aggregation on a specific column when the user explicitly requests it.

**CRITICAL: DO NOT AUTO-APPLY AGGREGATION**

This tool should ONLY be used when the user EXPLICITLY requests aggregation configuration.
Do NOT automatically apply aggregation during normal onboarding or master source creation.

**WHEN TO USE THIS TOOL:**
- User explicitly says: "Enable aggregation on column X"
- User explicitly says: "Configure aggregation for EntityIdentifier"
- User explicitly says: "Set up aggregation with column X as EntityIdentifier"
- User explicitly says: "I want to aggregate on column X"
- User explicitly says: "Apply aggregation on [column] for [master source]"

**WHEN NOT TO USE THIS TOOL:**
- During normal master source creation
- During file analysis
- During onboarding flow (unless user explicitly asks for aggregation)
- When user doesn't mention aggregation

**WHAT AGGREGATION DOES:**

Aggregation allows multiple records with the same EntityID to be grouped together based on a secondary identifier (EntityIdentifier). This is useful when:
- You have multiple line items per transaction
- You need to sum amounts across related records
- You want to reconcile at a parent-child level

**REQUIRED INFORMATION:**

Before applying aggregation, you MUST have:
1. **Master Source ID** - The ID of the master source to configure
2. **EntityIdentifier Column** - The column name that should be mapped as EntityIdentifier
3. **Lookup ID** - The ID of the lookup associated with this reconciliation (ask user if not available)

**AGGREGATION WORKFLOW:**

Step 1: WAIT FOR EXPLICIT USER REQUEST
- DO NOT proactively suggest or apply aggregation
- Only proceed if user explicitly asks for aggregation

Step 2: GATHER REQUIRED INFORMATION
- Ask for master_source_id if not available in context
- Ask which column should be the EntityIdentifier
- Ask for lookup_id if not available (this is critical for enabling aggregation)

Step 3: APPLY AGGREGATION
The tool will:
a. Fetch current master source configuration via GET API
b. Update master source with:
   - Append "EntityIdentifier" to unique_keys array
   - Update mapping_config to map the specified column to "EntityIdentifier" destination
c. Fetch current lookup configuration via GET API
d. Update lookup to enable aggregation for the EntityID column

**API CALLS MADE BY THIS TOOL:**

1. **GET Master Source:**
   GET /v1/admin-recon-saas/sources/get/{master_source_id}
   - Fetches current unique_keys and mapping_config

2. **PATCH Master Source:**
   PATCH /v1/admin-recon-saas/sources/update/{master_source_id}
   - Updates unique_keys to include "EntityIdentifier"
   - Updates mapping_config to map user's column to "EntityIdentifier" destination
   
3. **GET Lookup:**
   GET /v1/admin-recon-saas/lookup/{lookup_id}
   - Fetches current lookup configuration

4. **PATCH Lookup:**
   PATCH /v1/admin-recon-saas/lookup/{lookup_id}
   - Enables aggregation on the config item containing "EntityID" in Columns

**CONFIGURATION CHANGES:**

**Master Source unique_keys (APPENDED, NOT REPLACED):**
- The tool PRESERVES all existing unique_keys and APPENDS "EntityIdentifier" to them
- Before: ["EntityID"] (or whatever keys already exist)
- After: ["EntityID", "EntityIdentifier"] (existing keys + EntityIdentifier)
- IMPORTANT: Existing keys like "EntityID" are NEVER removed, only "EntityIdentifier" is added

**Master Source mapping_config:**
The column specified by user will have its destination changed to "EntityIdentifier"
Example: {"value": "", "source": "invoice_number", "destination": "EntityIdentifier"}

**Lookup config:**
The aggregation field is enabled for the config containing EntityID:
Before: {"aggregation": {"enabled": false, "conditions": null}}
After: {"aggregation": {"enabled": true, "conditions": null}}

**EXAMPLE CONVERSATIONS:**

**Example 1: User explicitly requests aggregation**

User: "Enable aggregation on the invoice_number column for Source A"

Assistant should:
1. Check context for master_source_id of Source A
2. If not available, ask: "What is the master_source_id for Source A?"
3. Ask: "What is the lookup_id for this reconciliation?"
4. Once all info is gathered, call the aggregation tool

**Example 2: User provides all information**

User: "I want to aggregate on column 'line_item_id' for master source ID 'ABC123' with lookup 'LKP456'"

Assistant should proceed directly with the aggregation tool call.

**Example 3: User does NOT request aggregation - DO NOT APPLY**

User: "Create master sources for my two files"
→ DO NOT apply aggregation. Just create the master sources.

User: "Set up reconciliation for my files"
→ DO NOT apply aggregation. Complete the setup without aggregation.

**TOOL PARAMETERS:**

Required:
- master_source_id: ID of the master source to configure
- source_name: Name of the source for reference
- entity_identifier_column: The column name to map as EntityIdentifier
- lookup_id: ID of the lookup to update for aggregation

Optional:
- environment: "local", "dev", or "prod" (defaults to "dev")

**ERROR HANDLING:**

- If master_source_id is invalid: Tool returns error with message to check the ID
- If lookup_id is invalid: Tool returns error with message to check the ID
- If column not found in mapping_config: Tool adds a new mapping entry
- If EntityIdentifier already exists: Tool skips adding duplicate

REMEMBER: Only apply aggregation when the user EXPLICITLY requests it!`

		messages := []mcp.PromptMessage{
			mcp.NewPromptMessage(
				mcp.RoleUser,
				mcp.NewTextContent(elaboratePrompt),
			),
		}

		return mcp.NewGetPromptResult(
			"Recon-SaaS Aggregation Config",
			messages,
		), nil
	}

	return server.ServerPrompt{
		Prompt:  prompt,
		Handler: handler,
	}
}

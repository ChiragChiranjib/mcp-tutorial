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
		mcp.WithPromptDescription("Apply data transformations to recon-saas master source columns with intelligent function selection"),
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

**TRANSFORMATION WORKFLOW:**

When a user asks to apply a transformation:
1. Identify which source they're referring to (Source A, Source B, or by name)
2. Identify the column(s) involved
3. Determine the appropriate transformation function
4. Ask for the output column name if creating a new column
5. Apply the transformation using the recon_transformation_config tool

**AVAILABLE TRANSFORMATION FUNCTIONS (%s focus):**

**Amount/Number Transformations:**
- abs_amount_parsing: Parse absolute amount (removes commas, handles negatives)
- abs_amount_in_paisa: Convert to paisa (multiply by 100)
- add_amount_cols: Sum multiple amount columns
- subtract_amount_cols: Subtract amounts from base
- percentage_of_number: Calculate percentage
- extract_amount_from_cols: Get first non-zero from multiple columns

**String Transformations:**
- append_multiple_columns: Concatenate multiple columns
- excel_mid: Extract substring from middle
- excel_left: Extract from left
- excel_right: Extract from right
- split: Split by delimiter and get part
- regex_exec: Extract using regex
- remove_prefix: Remove prefix
- remove_suffix: Remove suffix
- add_padding_prefix: Zero-pad to length
- replace_blank_string: Trim whitespace
- remove_single_quotes: Remove single quotes
- remove_double_quotes: Remove double quotes

**Date Transformations:**
- change_date_format: Convert between date formats
- date_normalization: Normalize to YYYY-MM-DD
- txn_date_extraction_generic: Parse various date formats including unix timestamps
- excel_to_datetime: Convert Excel serial date
- subtract_date: Subtract days
- add_date: Add days

**Other Transformations:**
- hard_code_value: Set constant value
- settlement_amount_from_debit_credit_cols: Derive from debit/credit
- get_field_from_notes: Extract from JSON notes

**EXAMPLE SCENARIOS:**

**Scenario 1: Parse Absolute Amount**
User: "For Source A, I want Subtotal to parse absolute amount and store it in Amount"
Action:
- transformation_function: "abs_amount_parsing"
- input_columns: ["Subtotal"]
- output_column: "Amount"
- Update mapping_config: Change Subtotal destination from "Amount" to "subtotal", add new mapping for "Amount"

**Scenario 2: Concatenate Columns**
User: "For Source B, combine RRN, TID, and MID into EntityID"
Action:
- transformation_function: "append_multiple_columns"
- input_columns: ["RRN", "TID", "MID"]
- output_column: "EntityID"
- Update mapping_config: Change RRN destination from "EntityID" to "rrn", add new mapping for "EntityID"

**Scenario 3: Change Date Format**
User: "Change Invoice Date format from MM/DD/YYYY to YYYY-MM-DD in Source A"
Action:
- transformation_function: "change_date_format"
- input_columns: ["Invoice Date"]
- additional_params: ["%%m/%%d/%%Y", "%%Y-%%m-%%d"]
- output_column: "Invoice Date" (same column, in-place transformation)
- mapping_config: Unchanged (same column)

**Scenario 4: Extract Substring**
User: "Extract the UTR from the description field starting at position 10 for 12 characters"
Action:
- transformation_function: "excel_mid"
- input_columns: ["description"]
- additional_params: [10, 12]
- output_column: "utr"
- Add new mapping for "utr"

**IMPORTANT RULES:**

1. **Column References**: Always use $ prefix for columns in the logic: "$column_name"

2. **Output Columns**: When transformation creates a new column that will be used for EntityID or Amount:
   - If an existing column was mapped to EntityID/Amount, remap it to snake_case
   - Add the new transformation output as the new EntityID/Amount source

3. **Mapping Config Updates**: 
   - If output_column matches an existing destination (like "Amount" or "EntityID"), the original source column needs remapping
   - Original column gets snake_case destination, new transformation output gets the special destination

4. **Ask for Clarification**:
   - Which source? (Source A, Source B, or by name)
   - Which columns are involved?
   - What should the output column be named?
   - For date formats: what is the current format and desired format?

**API ENDPOINT:**
Uses PATCH to /v1/admin-recon-saas/sources/update/{master_source_id}
Updates both transformation_config and mapping_config in a single call.

Please describe the transformation you want to apply, and I'll help you configure it correctly.`, transformationType)

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

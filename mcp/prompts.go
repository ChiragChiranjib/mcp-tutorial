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

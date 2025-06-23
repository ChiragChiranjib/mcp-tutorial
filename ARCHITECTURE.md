# MCP Tutorial Server Architecture

This document provides visual diagrams to help understand the different MCP (Model Context Protocol) server implementations in this tutorial project.

## Overview

The tutorial provides three different MCP server implementations:
- **SSE (Server-Sent Events)**: HTTP-based with real-time streaming
- **STDIO**: Standard input/output based communication  
- **Streamable HTTP**: HTTP-based with stateless streaming

All implementations share the same core MCP functionality but use different transport mechanisms.

## Client-Server Transport Patterns

### STDIO Transport Flow

The STDIO implementation uses standard input/output streams for direct process communication.

```mermaid
sequenceDiagram
    participant C as MCP Client
    participant P as Process
    participant S as MCP Server (STDIO)
    
    Note over C,S: Process Initialization
    C->>P: Start Server Process
    P->>S: Launch main.go
    S->>P: Ready on stdin/stdout
    
    Note over C,S: MCP Handshake
    C->>S: {"jsonrpc": "2.0", "method": "initialize", ...}
    Note over S: Read from stdin
    S->>C: {"jsonrpc": "2.0", "result": {"capabilities": ...}}
    Note over C: Read from stdout
    
    Note over C,S: Tool Execution
    C->>S: {"method": "tools/call", "params": {"name": "calculator", ...}}
    Note over S: Process via stdin
    S->>S: Execute Calculator Tool
    S->>C: {"result": {"content": [{"type": "text", "text": "2 + 2 = 4"}]}}
    Note over C: Receive via stdout
    
    Note over C,S: Resource Access
    C->>S: {"method": "resources/read", "params": {"uri": "system://status"}}
    S->>S: Generate System Status
    S->>C: {"result": {"contents": [{"uri": "system://status", ...}]}}
    
    Note over C,S: Process Termination
    C->>P: Terminate Process
    P->>S: SIGTERM/SIGINT
    S->>P: Clean Shutdown
```

### SSE Transport Flow

The SSE implementation uses HTTP with Server-Sent Events for real-time bidirectional communication.

```mermaid
sequenceDiagram
    participant C as MCP Client
    participant S as MCP Server (SSE) on Port 8080
    
    Note over C,S: SSE Connection Setup
    C->>S: GET /sse
    S->>S: Initialize SSE Connection
    S->>C: 200 OK + SSE Headers<br/>data: {"session_id": "abc123"}
    
    Note over C,S: Session Established
    C->>C: Store session_id = "abc123"
    S->>C: SSE Keep-Alive (every 10s)<br/>data: {"type": "ping"}
    
    Note over C,S: MCP Operations via Message Endpoint
    C->>S: POST /message?session_id=abc123<br/>{"method": "initialize", ...}
    S->>C: 202 Accepted (Request Queued)
    S->>S: Route to Session Handler & Process Initialize Request
    S->>C: SSE Push<br/>data: {"result": {"capabilities": ...}}
    
    Note over C,S: Tool Calls
    C->>S: POST /message?session_id=abc123<br/>{"method": "tools/call", "params": {...}}
    S->>C: 202 Accepted (Request Queued)
    S->>S: Route to Session Handler & Execute Calculator Tool
    S->>C: SSE Push<br/>data: {"result": {"content": [...]}}
    
    Note over C,S: Server-Initiated Notifications (Optional)
    S->>S: Generate Server Event/Notification
    S->>C: SSE Push<br/>data: {"type": "notification", ...}
    
    Note over C,S: Connection Cleanup
    C->>S: Close SSE Connection
    S->>S: Cleanup Session & Remove session_id
```

### Streamable HTTP Transport Flow

The Streamable HTTP implementation provides stateless HTTP communication with optional SSE upgrade for notifications.

```mermaid
sequenceDiagram
    participant C as MCP Client
    participant S as MCP Server (HTTP) on Port 8081
    
    Note over C,S: Direct MCP Endpoint Access
    C->>S: POST /mcp<br/>{"method": "initialize", ...}
    S->>S: Process Request (Stateless) & Handle Initialize
    S->>C: 200 OK + JSON Response<br/>{"result": {"capabilities": ...}}
    
    Note over C,S: Tool Execution (Stateless)
    C->>S: POST /mcp<br/>{"method": "tools/call", "params": {"name": "calculator", ...}}
    S->>S: Process Request (No Session) & Execute Calculator Tool
    S->>C: 200 OK + Tool Result<br/>{"result": {"content": [{"type": "text", "text": "5 * 3 = 15"}]}}
    
    Note over C,S: Resource Access (Stateless)
    C->>S: POST /mcp<br/>{"method": "resources/read", "params": {"uri": "math://constants"}}
    S->>S: Process Request (No Session) & Generate Math Constants
    S->>C: 200 OK + Resource Data<br/>{"result": {"contents": [...]}}
    
    Note over C,S: Optional SSE Upgrade for Notifications
    C->>S: POST /mcp<br/>{"method": "tools/get", "params": {"name": "calculator"}}
    
    alt Single HTTP Response
        S->>S: Process Request (Stateless)
        S->>C: 200 OK + JSON Response<br/>{"result": {"response": ...}}
    else Server Opens SSE Stream
        S->>S: Initialize SSE Session & Generate SSE stream for session
        S->>C: 200 OK + SSE Headers<br/>Connection: Keep-Alive
        
        Note over C,S: SSE Stream Active
        loop While Connection Remains Open
            S->>S: Generate SSE Messages
            S->>C: SSE Event: data: {"response": ...}
        end
        
        Note over C,S: Connection Cleanup
        C->>S: Close SSE Connection
        S->>S: Cleanup Session & Remove session_id
    end
```

## MCP Tutorial Server Architecture

This unified architecture diagram shows how the MCP Tutorial Server uses a pluggable transport design with shared business logic components.

```mermaid
graph TB
    subgraph "Shared Business Logic: /mcp Package"
        SHARED["Common MCP Components<br/><br/>Tools (mcp/tools.go)<br/>• CalculatorTool - 6 operations<br/>• SystemInfoTool - time/date info<br/><br/>Prompts (mcp/prompts.go)<br/>• MathTutorPrompt - tutoring<br/>• CodeReviewPrompt - analysis<br/><br/>Resources (mcp/resources.go)<br/>• SystemStatusResource - status<br/>• MathConstantsResource - constants"]
    end
    
    subgraph "Transport Implementations: /cmd Directory"
        subgraph "SSE Implementation"
            SSE_FLOW["cmd/sse/main.go<br/>1. Create MCP Server<br/>2. Add Shared Components<br/>3. Wrap with SSE Transport<br/>4. Start on Port 8080"]
            
            SSE_TRANSPORT["SSE Transport<br/>NewSSEServer<br/>• HTTP + Server-Sent Events<br/>• Stateful Sessions<br/>• Real-time bidirectional"]
        end
        
        subgraph "STDIO Implementation"
            STDIO_FLOW["cmd/stdio/main.go<br/>1. Create MCP Server<br/>2. Add Shared Components<br/>3. Wrap with STDIO Transport<br/>4. Listen on stdin/stdout"]
            
            STDIO_TRANSPORT["STDIO Transport<br/>NewStdioServer<br/>• Standard I/O Streams<br/>• Process Communication<br/>• Stateless"]
        end
        
        subgraph "HTTP Implementation"
            HTTP_FLOW["cmd/streamable_http/main.go<br/>1. Create MCP Server<br/>2. Add Shared Components<br/>3. Wrap with HTTP Transport<br/>4. Start on Port 8081"]
            
            HTTP_TRANSPORT["HTTP Transport<br/>NewStreamableHTTPServer<br/>• Pure HTTP Requests<br/>• Stateless + Optional SSE<br/>• REST-like calls"]
        end
    end
    
    subgraph "Client Connections"
        CLIENT1["MCP Client<br/>SSE Connection"]
        CLIENT2["MCP Client<br/>Process Pipes"]
        CLIENT3["MCP Client<br/>HTTP Requests"]
    end
    
    %% Single clean connection showing shared components are used by all
    SHARED -.-> SSE_FLOW
    SHARED -.-> STDIO_FLOW  
    SHARED -.-> HTTP_FLOW
    
    %% Flow within each implementation
    SSE_FLOW --> SSE_TRANSPORT
    STDIO_FLOW --> STDIO_TRANSPORT
    HTTP_FLOW --> HTTP_TRANSPORT
    
    %% Client connections
    SSE_TRANSPORT <--> CLIENT1
    STDIO_TRANSPORT <--> CLIENT2
    HTTP_TRANSPORT <--> CLIENT3
    
    %% Styling
    style SHARED fill:#f0f8f0,stroke:#4caf50,stroke-width:2px
    style SSE_FLOW fill:#f3e5f5,stroke:#9c27b0
    style STDIO_FLOW fill:#fff3e0,stroke:#ff9800
    style HTTP_FLOW fill:#fce4ec,stroke:#e91e63
    style SSE_TRANSPORT fill:#f3e5f5,stroke:#9c27b0
    style STDIO_TRANSPORT fill:#fff3e0,stroke:#ff9800
    style HTTP_TRANSPORT fill:#fce4ec,stroke:#e91e63
```

## Getting Started

### SSE Server
```bash
go run cmd/sse/main.go
# Server starts on http://localhost:8080
```

### STDIO Server  
```bash
go run cmd/stdio/main.go
# Communicates via stdin/stdout
```

### Streamable HTTP Server
```bash
go run cmd/streamable_http/main.go  
# Server starts on http://localhost:8081
```

## Core MCP Components

All implementations share these components:

### Tools
- **Calculator**: Performs basic math operations (add, subtract, multiply, divide, power, sqrt)
- **System Info**: Provides current time/date in various formats

### Prompts
- **Math Tutor**: Comprehensive math tutoring with customizable topics and levels
- **Code Review**: Detailed code analysis with language-specific guidance

### Resources
- **System Status**: Server status and uptime information (JSON)
- **Math Constants**: Common mathematical constants (π, e, φ, √2) with descriptions

## Quick Start Examples

### STDIO Example
```bash
# Start the server process and communicate via pipes
echo '{"jsonrpc":"2.0","method":"tools/list","id":1}' | ./bin/stdio
```

### SSE Example
```bash
# 1. Start server: ./bin/sse
# 2. Open SSE connection: curl -N http://localhost:8080/sse
# 3. Use session_id for operations: 
curl -X POST "http://localhost:8080/message?session_id=abc123" \
  -H "Content-Type: application/json" \
  -d '{"method":"tools/list","id":1}'
```

### Streamable HTTP Example
```bash
# 1. Start server: ./bin/streamable_http
# 2. Direct MCP calls (no session needed):
curl -X POST "http://localhost:8081/mcp" \
  -H "Content-Type: application/json" \
  -d '{"method":"tools/list","id":1}'
``` 